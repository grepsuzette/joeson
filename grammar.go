/*
Package joeson is a packrat left recursive parser in Go
ported from https://github.com/jaekwon/JoeScript 's joeson.coffee
*/
package joeson

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// See examples and tests to build a grammar.
// Then just use ParseString().
// Like in the original implementation, Grammar "is" a GNode.
type Grammar struct {
	*GNodeImpl
	rank     Parser         // a *Rank or a Ref to a rank
	NumRules int            // Each Ast can have rules, recursively. This however is the total count in the grammar
	id2Rule  map[int]Parser // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605
	TraceOptions
	wasInitialized bool
}

func (gm *Grammar) GetRule(name string) Parser {
	if x, exists := gm.GNodeImpl.rules[name]; exists {
		return x
	} else {
		return nil
	}
}

// General options to build a grammar
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
	// The name is also set afterwards in the coffeescript version
	newgm.SetName(name)
	newgm.postinit()
	return newgm
}

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
	return gm.Parse(newParseContext(code, gm.NumRules, attrs, gm.TraceOptions))
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
		for pos := ctx.Code.Pos; pos < len(ctx.Frames); pos++ {
			posFrames := ctx.Frames[pos]
			for _, frame := range posFrames {
				if frame != nil {
					maxAttempt = pos
					if frame.Result != nil {
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

/*
func (gm *Grammar) parseOrFail(ctx *ParseContext) (Ast, error) {
	// func (gm *Grammar) Parse(ctx *ParseContext) (Ast, error) {
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
		for pos := ctx.Code.Pos; pos < len(ctx.Frames); pos++ {
			posFrames := ctx.Frames[pos]
			for _, frame := range posFrames {
				if frame != nil {
					maxAttempt = pos
					if frame.Result != nil {
						maxSuccess = pos
						break
					}
				}
			}
		}
		// TODO this is kept as original, but seems it was not finished
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
		return nil, errors.New(sErr)
	}
	// joeson.coffee has a opts.returnContext but won't implement it
	return result, nil
}
*/

// -- after this are the lower level stuffs --

func newEmptyGrammar() *Grammar { return newEmptyGrammarWithOptions(DefaultTraceOptions()) }

func newEmptyGrammarWithOptions(opts TraceOptions) *Grammar {
	name := "__empty__"
	gm := &Grammar{NewGNode(), nil, 0, map[int]Parser{}, opts, false}
	gm.GNodeImpl.name = name
	gm.GNodeImpl.node = gm
	return gm
}

// Destroy the grammar. Only tests should use this.
func (gm *Grammar) Bomb() {
	gm.rank = newEmptyRank("bombd")
	gm.GNodeImpl = NewGNode()
	gm.NumRules = 0
	gm.id2Rule = nil
	gm.wasInitialized = false
}

func (gm *Grammar) GetGNode() *GNodeImpl { return gm.GNodeImpl }

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
		panic("You can only call grammar.Postinit() after some rank has been set")
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
				mono := monochoice.GetGNode()
				// Merge label
				if mono.label == "" {
					mono.label = choice.Label()
				}
				// Merge included rules
				for k, v := range choice.GetGNode().rules {
					mono.rules[k] = v
					mono.rulesK = append(mono.rulesK, k)
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
			gnode := node.GetGNode()
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
						gnode.rule = parent.GetGNode().rule
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
			gnode := node.GetGNode()
			if gnode == nil { // TODO delete this if ASAP, but need test, wip now
				return ""
			}
			if IsRule(node) {
				gm.GetGNode().rulesK = append(gm.GetGNode().rulesK, gnode.name)
				gm.GetGNode().rules[gnode.name] = node
				gnode.id = gm.NumRules
				gm.NumRules++
				gm.id2Rule[gnode.id] = node
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
	fmt.Println("| rules        : " + strconv.Itoa(gm.NumRules))
	fmt.Println("| ")
	if gm.NumRules <= 0 {
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
	for i := 0; i < gm.NumRules; i++ {
		v := gm.id2Rule[i]
		sParentName := "-"
		if v.GetGNode().parent != nil {
			switch father := v.GetGNode().parent.(type) {
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
			helpers.PadLeft(strconv.Itoa(v.GetGNode().id), 3),
			helpers.PadLeft(helpers.TypeOfToString(v), 20),
			helpers.PadLeft(helpers.BoolToString(v.Capture()), 3),
			helpers.PadLeft(v.Label(), 7),
			helpers.PadLeft(strings.Join(v.GetGNode().labels_.Get(), ","), 21),
			helpers.PadLeft(sParentName, 16),
			helpers.PadLeft(v.ContentString(), 30),
		)
	}
	fmt.Println("| ")
}
