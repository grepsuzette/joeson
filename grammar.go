// joeson is a packrat left recursive parser in Go
package joeson

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

type Grammar struct {
	*gnodeimpl
	rank     Parser         // a *Rank or a Ref to a rank
	numrules int            // Each Ast can have rules, recursively. This however is the total count in the grammar
	id2rule  map[int]Parser // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605
	TraceOptions
	wasInitialized bool
}


type GrammarOptions struct {
	// Options governing what is traced or not during the initialization or the parsing
	TraceOptions TraceOptions

	// Leave this nil unless you know. Nil specifies the joeson_handcompiled grammar.
	// This lazy function must return the (compiled) grammar to use when some uncompiled
	// string rules (sLine) are encountered. Once a grammar has been compiled, it can
	// parse and therefore be used here.
	LazyGrammar *helpers.Lazy[*Grammar]
}

// Make a new grammar from the rules in `lines`.
// A Rank will be internally created.
// Options can be omitted.
func GrammarFromLines(lines []Line, name string, options ...GrammarOptions) *Grammar {
	var opts GrammarOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = GrammarOptions{
			TraceOptions: DefaultTraceOptions(),
			LazyGrammar:  nil,
		}
	}
	ranke := rankFromLines(lines, name, opts)
	newgm := newEmptyGrammarWithOptions(opts.TraceOptions)
	newgm.rank = ranke
	newgm.SetName(name)
	newgm.postinit()
	return newgm
}

func (gm *Grammar) CountRules() int { return gm.numrules }

// Parse functions don't panic.
// A parser returns nil when declining to parse. When there is an error
// but the parser takes the responsability (denying any other parser the
// chance to parse), it returns a ParseError.
func (gm *Grammar) ParseString(sCode string, attrs ...ParseOptions) Ast {
	if len(attrs) > 0 {
		return gm.ParseCode(NewCodeStream(sCode), attrs[0])
	} else {
		return gm.ParseCode(NewCodeStream(sCode), ParseOptions{})
	}
}

func (gm *Grammar) ParseCode(code *CodeStream, attrs ParseOptions) Ast {
	return gm.Parse(newParseContext(code, gm.numrules, attrs, gm.TraceOptions))
}

func (gm *Grammar) Parse(ctx *ParseContext) Ast {
	var oldTrace bool
	if ctx.Debug {
		// temporarily enable stack tracing
		oldTrace = gm.TraceOptions.Stack
		gm.TraceOptions.Stack = true
	}
	if gm.rank == nil {
		panic("Grammar.rank is nil")
	}
	result := gm.rank.Parse(ctx)
	// undo temporary stack tracing
	if ctx.Debug {
		gm.TraceOptions.Stack = oldTrace
	}
	// if parse is incomplete, compute error message
	if ctx.Code.Pos != ctx.Code.Length() {
		// find the maximum parsed entity
		maxAttempt := ctx.Code.Pos
		maxSuccess := ctx.Code.Pos
		for pos := ctx.Code.Pos; pos < len(ctx.frames); pos++ {
			posFrames := ctx.frames[pos]
			for _, frame := range posFrames {
				if frame != nil {
					maxAttempt = pos
					if frame.result != nil {
						maxSuccess = pos
						break
					}
				}
			}
		}
		sErr := fmt.Sprintf("Error parsing at char:%d=(line:%d,col:%d).", maxSuccess, ctx.Code.PosToLine(maxSuccess), ctx.Code.PosToCol(maxSuccess))
		sErr += "\n" + ctx.Code.Print()
		sErr += "\nDetails:\n"
		sErr += green("OK") + "/"
		sErr += yellow("Parsing") + "/"
		sErr += red("Suspect") + "/"
		sErr += white("Unknown") + "\n\n"
		sErr += green(ctx.Code.Peek(NewPeek().BeforeLines(2)))
		sErr += yellow(ctx.Code.Peek(NewPeek().AfterChars(maxSuccess - ctx.Code.Pos)))
		ctx.Code.Pos = maxSuccess
		sErr += red(ctx.Code.Peek(NewPeek().AfterChars(maxAttempt-ctx.Code.Pos))) + "/"
		ctx.Code.Pos = maxAttempt
		sErr += white(ctx.Code.Peek(NewPeek().AfterLines(2))) + "\n"
		return NewParseError(ctx, sErr)
	}
	return result
}

// -- after this are the lower level stuffs --

func newEmptyGrammar() *Grammar { return newEmptyGrammarWithOptions(DefaultTraceOptions()) }

func newEmptyGrammarWithOptions(opts TraceOptions) *Grammar {
	name := "__empty__"
	gm := &Grammar{NewGNode(), nil, 0, map[int]Parser{}, opts, false}
	gm.gnodeimpl.name = name
	gm.gnodeimpl.node = gm
	return gm
}

// Destroy the grammar. Only tests should use this.
func (gm *Grammar) Bomb() {
	gm.rank = newEmptyRank("bombd")
	gm.gnodeimpl = NewGNode()
	gm.numrules = 0
	gm.id2rule = nil
	gm.wasInitialized = false
}

func (gm *Grammar) getgnode() *gnodeimpl { return gm.gnodeimpl }

func (gm *Grammar) getRule(name string) Parser {
	if x, exists := gm.gnodeimpl.rules[name]; exists {
		return x
	} else {
		return nil
	}
}


func (gm *Grammar) Prepare()                {}
func (gm *Grammar) HandlesChildLabel() bool { return false }
func (gm *Grammar) ContentString() string {
	if gm.rank == nil {
		return magenta("GRAMMAR{}")
	} else {
		return magenta("GRAMMAR{") + String(gm.rank) + magenta("}")
	}
}

func (gm *Grammar) IsReady() bool { return gm.rank != nil && gm.wasInitialized }

func (gm *Grammar) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren rank: {type:Rank}
	gm.GetGNode().rules = ForEachChild_InRules(gm, f)
	if gm.rank != nil {
		gm.rank = f(gm.rank)
	}
	return gm
}

// after Rank has already been set,
// collect and collect rules, simplify the rule tree etc.
func (gm *Grammar) postinit() {
	if gm.rank == nil {
		panic("grammar.rank is nil")
	}
	opts := gm.TraceOptions

	// from joeson.coffee
	// "TODO refactor into translation passes." <empty>

	// Merge Choices with just a single choice.
	walkOptimizeAwayMonochoice := WalkPrepost{Pre: nil} // predeclare, because cb will need to call itself just below
	walkOptimizeAwayMonochoice.Pre = func(node Parser, parent Parser) string {
		/*
			BEFORE:	  somenode            ** foreach child, if one child is a mono choice, optimize it away
					   \
					   Choices(1)
						 \
						  node below.unaffected

			AFTER:	  somenode+choiceattrs
						  \
						   node below.unaffected
		*/
		node.ForEachChild(func(child Parser) Parser {
			if choice, ok := child.(*choice); ok && choice.isMonoChoice() {
				monochoice := choice.choices[0]
				mono := monochoice.getgnode()
				// Merge label
				if mono.label == "" {
					mono.label = choice.Label()
				}
				// Merge included rules
				for k, v := range choice.GetGNode().rules {
					mono.rules[k] = v
				}
				// grandchild becomes the child
				return monochoice
			} else {
				return Walk(child, node, walkOptimizeAwayMonochoice)
			}
		})
		return ""
	}
	Walk(gm, nil, walkOptimizeAwayMonochoice)

	// Connect all the nodes and collect dereferences into @rules
	Walk(gm, nil, WalkPrepost{
		Pre: func(node Parser, parent Parser) string {
			gnode := node.getgnode()
			if gnode == nil {
				return ""
			}
			// sanity check: it must have no parent yet if it's not a rule
			if !IsRule(node) && gnode != nil && gnode.parent != nil {
				panic("Grammar tree should be a DAG, nodes should not be referenced more than once.")
			}

			gnode.grammar = gm
			gnode.parent = parent
			if false {
				// "inline rules are special" in original coffeescript
				// but the bit of code seem unreachable anyway
				panic("assert")
			} else {
				// set node.rule, the root node for this rule
				if gnode.rule == nil {
					if parent != nil {
						gnode.rule = parent.getgnode().rule
					} else {
						// TODO gnode.Rule = NewNativeUndefined()
						//   we used nil here if there is any pb...
						gnode.rule = nil
					}
				}
			}
			return ""
		},
		Post: func(node Parser, parent Parser) string {
			gnode := node.getgnode()
			if IsRule(node) {
				gm.GetGNode().rules[gnode.name] = node
				gnode.id = gm.numrules
				gm.numrules++
				gm.id2rule[gnode.id] = node
				if opts.Loop { // print out id->rulename for convenience
					fmt.Println("Loop " + red(strconv.Itoa(gnode.id)) + ":\t" + String(node))
				}
			}
			return ""
		},
	})

	// Prepare all the nodes, children first.
	Walk(gm, nil, WalkPrepost{
		Post: func(node Parser, parent Parser) string {
			node.Prepare()
			return ""
		},
	})
	gm.wasInitialized = true

	if opts.Grammar {
		gm.PrintRules()
	}
}

func (gm *Grammar) PrintRules() {
	fmt.Println("+--------------- Grammar.Debug() ----------------------------------")
	fmt.Println("| name         : " + bold(gm.Name()))
	fmt.Println("| contentString: " + gm.ContentString())
	fmt.Println("| rules        : " + strconv.Itoa(gm.numrules))
	fmt.Println("| ")
	if gm.numrules <= 0 {
		return
	}
	fmt.Println("| ",
		helpers.PadLeft("key/name", 14),
		helpers.PadLeft("id", 3),
		helpers.PadLeft("type", 20),
		helpers.PadLeft("cap", 3),
		helpers.PadLeft("label", 7),
		helpers.PadLeft("labels()", 21),
		helpers.PadLeft("parent.name", 16),
		helpers.PadLeft("contentString", 30),
	)
	fmt.Println("|   -------------------------------------------------------------------------------------")
	for i := 0; i < gm.numrules; i++ {
		v := gm.id2rule[i]
		sParentName := "-"
		if v.getgnode().parent != nil {
			switch father := v.getgnode().parent.(type) {
			// case *Grammar:
			// 	sParentName = "__grammar__" // instead show name, use same as js for diffing
			case *rank:
				// 2 kind of ranks:
				//   1. parent is a grammar
				//   2. regular rank (subrank)
				switch father.parent.(type) {
				case *Grammar:
					sParentName = "__grammar__" // instead show name, use same as js for diffing
				default:
					sParentName = father.Name()
				}
			default:
				// sParentName = fmt.Sprintf("%T", v)
				sParentName = v.Name()
			}
			// sParentName = v.GetGNode().Parent.GetGNode().Name
		}
		fmt.Println("|  ",
			helpers.PadLeft(v.Name(), 14),
			helpers.PadLeft(strconv.Itoa(v.getgnode().id), 3),
			helpers.PadLeft(helpers.TypeOfToString(v), 20),
			helpers.PadLeft(helpers.BoolToString(v.Capture()), 3),
			helpers.PadLeft(v.Label(), 7),
			helpers.PadLeft(strings.Join(v.getgnode().labels_.Get(), ","), 21),
			helpers.PadLeft(sParentName, 16),
			helpers.PadLeft(v.ContentString(), 30),
		)
	}
	fmt.Println("| ")
}
