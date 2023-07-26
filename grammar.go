// joeson is a packrat left recursive parser in Go
package joeson

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

type Grammar struct {
	Attr
	*gnodeimpl
	rank     Parser         // a *Rank or a Ref to a rank
	numrules int            // Each Ast can have rules, recursively. This however is the total count in the grammar
	id2rule  map[int]Parser // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605
	TraceOptions
	wasInitialized bool
}

type GrammarOptions struct {
	// Govern what is traced during initialization or parsing
	TraceOptions TraceOptions

	// Leave this nil unless you know what you're doing.
	// This lazy function must return the grammar to use when some uncompiled
	// string rules (sLine) are encountered. This is internally used to bootstrap
	// the grammar. `nil` will have the grammar use the joeson_handcompiled
	// grammar to parse the given grammar.
	LazyGrammar *helpers.Lazy[*Grammar]
}

// Prepare a new grammar from the rules in `lines`.
// The way to pass options will need to be reworked at some point.
// Here is an example:
// ```
//
//	gm := joeson.GrammarFromLines([]joeson.Line{
//			o(named("Input", "expr:Expression")),
//			i(named("Expression", "Expression _ binary_op _ Expression | UnaryExpr")),
//			i(named("binary_op", "'+'")),
//			i(named("UnaryExpr", "[0-9]+")),
//			i(named("_", "[ \t]*")),
//		}, "leftRecursion", joeson.GrammarOptions{TraceOptions: joeson.Verbose()})
//
// ```
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
	rank := rankFromLines(lines, name, opts)
	newgm := newEmptyGrammarWithOptions(opts.TraceOptions)
	newgm.rank = rank
	newgm.SetRuleName(name)
	newgm.postinit()
	return newgm
}

func (gm *Grammar) assertNode()     {}
func (gm *Grammar) CountRules() int { return gm.numrules }

// Parse functions don't panic.
// A parser returns nil when it cannot to parse. When there is an error
// but the parser takes the responsability (denying any other parser the
// chance to parse), it returns a ParseError instead.
func (gm *Grammar) ParseString(sCode string) Ast {
	return gm.ParseCode(NewCodeStream(sCode))
}

// CodeStream comes from original Joeson implementation
// Prefer to use ParseString() or ParseTokens()
func (gm *Grammar) ParseCode(code *CodeStream) Ast {
	return gm.Parse(newParseContext(code, gm.numrules, gm.TraceOptions))
}

func (gm *Grammar) Parse(ctx *ParseContext) Ast {
	var oldTrace bool
	ctx.GrammarName = gm.GetRuleName()
	if ctx.parseOptions.Debug {
		// temporarily enable stack tracing
		oldTrace = gm.TraceOptions.Stack
		gm.TraceOptions.Stack = true
	}
	if gm.rank == nil {
		panic("Grammar.rank is nil")
	}
	result := gm.rank.Parse(ctx)
	// undo temporary stack tracing
	if ctx.parseOptions.Debug {
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
		sErr += Green("OK") + "/"
		sErr += Yellow("Parsing") + "/"
		sErr += Red("Suspect") + "/"
		sErr += White("Unknown") + "\n\n"
		sErr += Green(ctx.Code.Peek(NewPeek().BeforeLines(2)))
		sErr += Yellow(ctx.Code.Peek(NewPeek().AfterChars(maxSuccess - ctx.Code.Pos)))
		ctx.Code.Pos = maxSuccess
		sErr += Red(ctx.Code.Peek(NewPeek().AfterChars(maxAttempt-ctx.Code.Pos))) + "/"
		ctx.Code.Pos = maxAttempt
		sErr += White(ctx.Code.Peek(NewPeek().AfterLines(2))) + "\n"
		return NewParseError(ctx, sErr)
	}
	return result
}

// lower level stuffs

func newEmptyGrammar() *Grammar { return newEmptyGrammarWithOptions(DefaultTraceOptions()) }

func newEmptyGrammarWithOptions(opts TraceOptions) *Grammar {
	name := "__empty__"
	gm := &Grammar{newAttr(), newGNode(), nil, 0, map[int]Parser{}, opts, false}
	gm.gnodeimpl.name = name
	gm.gnodeimpl.node = gm
	return gm
}

// Destroy the grammar. Only tests should use this.
func (gm *Grammar) Bomb() {
	gm.rank = newEmptyRank("bombd")
	gm.gnodeimpl = newGNode()
	gm.numrules = 0
	gm.id2rule = nil
	gm.wasInitialized = false
}

func (gm *Grammar) gnode() *gnodeimpl { return gm.gnodeimpl }

func (gm *Grammar) getRule(name string) Parser {
	if x, exists := gm.gnodeimpl.rules[name]; exists {
		return x
	} else {
		return nil
	}
}

func (gm *Grammar) prepare()                {}
func (gm *Grammar) handlesChildLabel() bool { return false }
func (gm *Grammar) String() string {
	if gm.rank == nil {
		return Magenta("GRAMMAR{}")
	} else {
		return Magenta("GRAMMAR{") + String(gm.rank) + Magenta("}")
	}
}

func (gm *Grammar) IsReady() bool { return gm.rank != nil && gm.wasInitialized }

func (gm *Grammar) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren rank: {type:Rank}
	gm.rules = ForEachChildInRules(gm, f)
	if gm.rank != nil {
		gm.rank = f(gm.rank)
	}
	return gm
}

// after Rank has already been set,
// connect and collect rules, simplify the rule tree, set wasInitialized.
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
				mono := monochoice.gnode()
				// Merge label
				if mono.label == "" {
					mono.label = choice.GetRuleLabel()
				}
				// Merge included rules
				for k, v := range choice.rules {
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
			gnode := node.gnode()
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
						gnode.rule = parent.gnode().rule
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
			gnode := node.gnode()
			if IsRule(node) {
				gm.rulesK = append(gm.rulesK, gnode.name)
				gm.rules[gnode.name] = node
				gnode.id = gm.numrules
				gm.numrules++
				gm.id2rule[gnode.id] = node
				if opts.Loop { // print out id->rulename for convenience
					fmt.Println("Loop " + Red(strconv.Itoa(gnode.id)) + ":\t" + String(node))
				}
			}
			return ""
		},
	})

	// Prepare all the nodes, children first.
	Walk(gm, nil, WalkPrepost{
		Post: func(node Parser, parent Parser) string {
			node.prepare()
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
	fmt.Println("| name         : " + Bold(gm.GetRuleName()))
	fmt.Println("| contentString: " + gm.String())
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
		if v.gnode().parent != nil {
			switch father := v.gnode().parent.(type) {
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
					sParentName = father.GetRuleName()
				}
			default:
				// sParentName = fmt.Sprintf("%T", v)
				sParentName = v.GetRuleName()
			}
			// sParentName = v.Parent.Name
		}
		fmt.Println("|  ",
			helpers.PadLeft(v.GetRuleName(), 14),
			helpers.PadLeft(strconv.Itoa(v.gnode().id), 3),
			helpers.PadLeft(helpers.TypeOfToString(v), 20),
			helpers.PadLeft(helpers.BoolToString(v.Capture()), 3),
			helpers.PadLeft(v.GetRuleLabel(), 7),
			helpers.PadLeft(strings.Join(v.gnode().cachedLabels.Get(), ","), 21),
			helpers.PadLeft(sParentName, 16),
			helpers.PadLeft(v.String(), 30),
		)
	}
	fmt.Println("| ")
}
