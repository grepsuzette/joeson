package joeson

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

type Grammar struct {
	*Attr
	*TraceOptions
	*rule
	rank           Parser         // The toplevel rank with which the grammar was defined
	numrules       int            // Each Ast can have rules, recursively. This however is the total count in the grammar
	id2rule        map[int]Parser // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605
	wasInitialized bool
}

type GrammarOptions struct {
	// Govern what is traced during initialization (unless SkipSetup == true) or parsing
	TraceOptions *TraceOptions

	// Leave this nil unless you know what you're doing.
	// This lazy function must return the grammar to use when some uncompiled
	// string rules (sLine) are encountered. This is internally used to bootstrap
	// the grammar. `nil` will have the grammar use the joeson_handcompiled
	// grammar to parse the given grammar.
	LazyGrammar *helpers.Lazy[*Grammar]
}

// Prepare a new grammar from the rules in `lines`.
func GrammarFromLines(name string, lines []Line) *Grammar {
	return GrammarWithOptionsFromLines(
		name,
		GrammarOptions{
			TraceOptions: DefaultTraceOptions(),
			LazyGrammar:  nil,
		},
		lines,
	)
}

// The way to pass options will need to be reworked at some point.
// Here is an example:
// ```
//
//	gm := j.GrammarWithOptionsFromLines(
//		"leftRecursion",
//		j.GrammarOptions{TraceOptions: j.Verbose()},
//		[]j.Line{
//			o(named("Expression", `Expression _ binary_op _ Expression | UnaryExpr`)),
//			i(named("binary_op", `'+'`)),
//			i(named("UnaryExpr", `[0-9]+`)),
//			i(named("_", `[ \t]*`)),
//		},
//	)
//
// ```
func GrammarWithOptionsFromLines(name string, options GrammarOptions, lines []Line) *Grammar {
	if options.TraceOptions == nil {
		options.TraceOptions = Mute()
	}
	rank := rankFromLines(lines, name, options)
	newgm := newEmptyGrammarWithOptions(options.TraceOptions)
	newgm.rank = rank
	newgm.getRule().name = name
	newgm.postinit()
	return newgm
}

func (gm *Grammar) assertNode()     {}
func (gm *Grammar) CountRules() int { return gm.numrules }

// Parse a string with given parse options.
// E.g. `gm.ParseString("blah", j.Debug{true})`.
// Return nil if failed to parse.
func (gm *Grammar) ParseString(s string, options ...ParseOption) Ast {
	ctx := newParseContext(
		NewRuneStream(s),
		gm.numrules,
		gm.TraceOptions,
	)
	for _, option := range options {
		ctx.applyOption(option)
	}
	return gm.parse(ctx)
}

// Parse provided TokenStream.
// E.g. `gm.ParseString("blah", j.Debug{true})`.
// Return nil if failed to parse.
func (gm *Grammar) ParseTokens(tokens *TokenStream, options ...ParseOption) Ast {
	ctx := newParseContext(
		tokens,
		gm.numrules,
		gm.TraceOptions,
	)
	for _, option := range options {
		ctx.applyOption(option)
	}
	return gm.parse(ctx)
}

// Use ParseString() or ParseTokens()
// This parse(ctx) is for Grammar to implement Parser
func (gm *Grammar) parse(ctx *ParseContext) Ast {
	oldTrace := gm.TraceOptions.Stack
	ctx.GrammarName = gm.getRule().name
	if ctx.parseOptions.debug {
		// temporarily enable stack tracing
		gm.TraceOptions.Stack = true
	}
	if gm.rank == nil {
		panic("Grammar.rank is nil")
	}
	result := gm.rank.parse(ctx)
	// undo temporary stack tracing
	gm.TraceOptions.Stack = oldTrace
	// if parse is incomplete, compute error message
	if ctx.Code.Pos() != ctx.Code.workLength() {
		// find the maximum parsed entity
		maxAttempt := ctx.Code.Pos()
		maxSuccess := ctx.Code.Pos()
		for pos := ctx.Code.Pos(); pos < len(ctx.frames); pos++ {
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
		sErr += Green(ctx.Code.PeekLines(-2))
		sErr += Yellow(ctx.Code.PeekRunes(maxSuccess - ctx.Code.Pos()))
		ctx.Code.SetPos(maxSuccess)
		sErr += Red(ctx.Code.PeekRunes(maxAttempt-ctx.Code.Pos())) + "/"
		ctx.Code.SetPos(maxAttempt)
		sErr += White(ctx.Code.PeekLines(+2)) + "\n"
		return NewParseError(ctx, sErr)
	}
	return result
}

// lower level stuffs

func newEmptyGrammar() *Grammar { return newEmptyGrammarWithOptions(DefaultTraceOptions()) }

func newEmptyGrammarWithOptions(opts *TraceOptions) *Grammar {
	name := "__empty__"
	gm := &Grammar{newAttr(), opts, newRule(), nil, 0, map[int]Parser{}, false}
	gm.rule.name = name
	gm.rule.node = gm
	return gm
}

// Destroy the grammar. Only tests should use this
// (it is used by bootstrapping test)
func (gm *Grammar) Bomb() {
	gm.rank = newEmptyRank("bombd")
	gm.rule = newRule()
	gm.numrules = 0
	gm.id2rule = nil
	gm.wasInitialized = false
}

func (gm *Grammar) getRule() *rule { return gm.rule }

func (gm *Grammar) getRuleRef(name string) Parser {
	if x, exists := gm.rule.rules[name]; exists {
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

func (gm *Grammar) forEachChild(f func(Parser) Parser) Parser {
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
		node.forEachChild(func(child Parser) Parser {
			if choice, ok := child.(*choice); ok && choice.isMonoChoice() {
				monochoice := choice.choices[0]
				mono := monochoice.getRule()
				// Merge label
				if mono.label == "" {
					mono.label = choice.getRule().label
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
			rule := node.getRule()
			if rule == nil {
				return ""
			}
			// sanity check: it must have no parent yet if it's not a rule
			if !IsRule(node) && rule != nil && rule.parent != nil {
				panic("Grammar tree should be a DAG, nodes should not be referenced more than once.")
			}

			rule.grammar = gm
			rule.parent = parent
			if false {
				// "inline rules are special" in original coffeescript
				// but the bit of code seem unreachable anyway
				panic("assert")
			} else {
				// set node.rule, the root node for this rule
				if rule.parser == nil {
					if parent != nil {
						rule.parser = parent.getRule().parser
					} else {
						// TODO rule.Rule = NewNativeUndefined()
						//   we used nil here if there is any pb...
						rule.parser = nil
					}
				}
			}
			return ""
		},
		Post: func(node Parser, parent Parser) string {
			rule := node.getRule()
			if IsRule(node) {
				gm.rulesK = append(gm.rulesK, rule.name)
				gm.rules[rule.name] = node
				rule.id = gm.numrules
				gm.numrules++
				gm.id2rule[rule.id] = node
				if opts.Loop { // print out id->rulename for convenience
					fmt.Println("Loop " + Red(strconv.Itoa(rule.id)) + ":\t" + String(node))
				}
			}
			return ""
		},
	})

	// Prepare all the nodes, children first.
	Walk(gm, nil, WalkPrepost{
		Post: func(parser Parser, parent Parser) string {
			parser.prepare()
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
	fmt.Println("| name         : " + Bold(gm.getRule().name))
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
		if v.getRule().parent != nil {
			switch father := v.getRule().parent.(type) {
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
					sParentName = father.getRule().name
				}
			default:
				sParentName = v.getRule().name
			}
		}
		fmt.Println("|  ",
			helpers.PadLeft(v.getRule().name, 14),
			helpers.PadLeft(strconv.Itoa(v.getRule().id), 3),
			helpers.PadLeft(helpers.TypeOfToString(v), 20),
			helpers.PadLeft(helpers.BoolToString(v.getRule().capture), 3),
			helpers.PadLeft(v.getRule().label, 7),
			helpers.PadLeft(strings.Join(v.getRule().labels_.Get(), ","), 21),
			helpers.PadLeft(sParentName, 16),
			helpers.PadLeft(v.String(), 30),
		)
	}
	fmt.Println("| ")
}
