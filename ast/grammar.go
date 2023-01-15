package ast

import (
	"errors"
	"fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"strconv"
	"strings"
)

// Grammar is the one to use in this package.
// See examples and tests to build a grammar.
// Then just use ParseString().

type Grammar struct {
	*GNode
	rank           Ast         // a *Rank or a Ref to a rank
	NumRules       int         // Each Ast can have rules, recursively. This however i the total count in the grammar
	Id2Rule        map[int]Ast // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605
	wasInitialized bool
}

func NewEmptyGrammar() *Grammar { return NewEmptyGrammarNamed("empty grammar") }
func NewEmptyGrammarNamed(name string) *Grammar {
	gm := &Grammar{NewGNode(), nil, 0, map[int]Ast{}, false}
	gm.GNode.Name = name
	gm.GNode.Node = gm
	return gm
}

// ♥ call this one (MAIN GRAMMAR PARSE FUNCTION)
func (gm *Grammar) ParseString(sCode string, attrs ...ParseOptions) (Ast, error) {
	if len(attrs) > 0 {
		return gm.ParseCode(NewCodeStream(sCode), attrs[0])
	} else {
		return gm.ParseCode(NewCodeStream(sCode), ParseOptions{})
	}
}

func (gm *Grammar) ParseCode(code *CodeStream, attrs ParseOptions) (Ast, error) {
	return gm.parseOrFail(NewParseContext(code, gm, attrs))
}

// -- after this are the lower level stuffs --

func (gm *Grammar) GetGNode() *GNode { return gm.GNode }

// this one conforms the interface, but you would normally call
// grammar.ParseString() or grammar.ParseCode().
func (gm *Grammar) Parse(ctx *ParseContext) Ast {
	if ast, error := gm.parseOrFail(ctx); error == nil {
		return ast
	} else {
		panic(error)
	}
}

func (gm *Grammar) parseOrFail(ctx *ParseContext) (Ast, error) {
	var oldTrace TraceSettings
	if ctx.Debug {
		// temporarily enable stack tracing
		oldTrace = Trace
		Trace.Stack = true
	}
	if gm.rank == nil {
		panic("Grammar.rank is nil")
	}
	result := gm.rank.Parse(ctx)
	// undo temprary stack tracing
	if ctx.Debug {
		Trace = oldTrace
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
		return nil, errors.New(sErr)
	}
	// joeson.coffee has a opts.returnContext but won't implement it
	return result, nil
}

func (gm *Grammar) Prepare()                {}
func (gm *Grammar) HandlesChildLabel() bool { return false }
func (gm *Grammar) ContentString() string {
	return Magenta("GRAMMAR{") + helpers.TypeOfToString(gm.rank) + Prefix(gm.rank) + gm.rank.ContentString() + Magenta("}")
}

func (gm *Grammar) CountRules() int { return gm.NumRules }
func (gm *Grammar) IsReady() bool   { return gm.rank != nil && gm.wasInitialized }
func (gm *Grammar) SetRankIfEmpty(rank Ast) {
	if gm.rank != nil {
		return
	}
	if gm.IsReady() {
		panic("Grammar is already defined and can not be changed on the fly at the moment")
	}
	gm.rank = rank
}
func (gm *Grammar) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren rank: {type:Rank}
	gm.GetGNode().Rules = ForEachChild_InRules(gm, f)
	if gm.rank != nil {
		gm.rank = f(gm.rank)
	}
	return gm
}

// after Rank has already been set,
// collect and collect rules, simplify the rule tree etc.
func (gm *Grammar) Postinit() {
	if gm.rank == nil {
		panic("You can only call grammar.Postinit() after some rank has been set")
	}

	// from joeson.coffee
	// "TODO refactor into translation passes." <empty>

	// Merge Choices with just a single choice.
	walkOptimizeAwayMonochoice := WalkPrepost{Pre: nil}
	walkOptimizeAwayMonochoice.Pre = func(node Ast, parent Ast) string {
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
		node.ForEachChild(func(child Ast) Ast {
			if choice, ok := child.(*Choice); ok && choice.IsMonoChoice() {
				monochoice := choice.choices[0]
				mono := monochoice.GetGNode()
				// Merge label
				if mono.Label == "" {
					mono.Label = choice.GetGNode().Label
				}
				// Merge included rules
				for k, v := range choice.GetGNode().Rules {
					mono.Rules[k] = v
					mono.RulesK = append(mono.RulesK, k)
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
		Pre: func(node Ast, parent Ast) string {
			gnode := node.GetGNode()
			if gnode == nil {
				return ""
			}
			// sanity check: it must have no parent yet if it's not a rule
			if !IsRule(node) && gnode != nil && gnode.Parent != nil {
				panic("Grammar tree should be a DAG, nodes should not be referenced more than once.")
			}

			gnode.Grammar = gm
			gnode.Parent = parent
			if false {
				// "inline rules are special" in original coffeescript
				// but the bit of code seem unreachable anyway
				panic("assert")
			} else {
				// set node.rule, the root node for this rule
				if gnode.Rule == nil {
					var r Ast
					if parent != nil {
						r = parent.GetGNode().Rule
					} else {
						r = NewNativeUndefined()
					}
					gnode.Rule = r
				}
			}
			return ""
		},
		Post: func(node Ast, parent Ast) string {
			gnode := node.GetGNode()
			if gnode == nil {
				return ""
			}
			if IsRule(node) {
				gm.GetGNode().RulesK = append(gm.GetGNode().RulesK, gnode.Name)
				gm.GetGNode().Rules[gnode.Name] = node
				gnode.Id = gm.NumRules
				gm.NumRules++
				gm.Id2Rule[gnode.Id] = node
				if Trace.Loop { // print out id->rulename for convenience
					fmt.Println("Loop " + Red(strconv.Itoa(gnode.Id)) + ":\t" + Prefix(node) + node.ContentString())
				}
			}
			return ""
		},
	})

	// just show the tree (very low level debug)
	// nbNodes := 0
	// Walk(gm, nil, WalkPrepost{
	// 	Pre: func(node Astnode, parent Astnode) string {
	// 		s := "undefined"
	// 		if parent != nil {
	// 			s = parent.ContentString()
	// 		}
	// 		fmt.Println("grammar PRE node:" + node.ContentString() + "/" + helpers.TypeOfToString(node) + " parent:" + s)
	// 		nbNodes++
	// 		depth := func(ast Astnode) int {
	// 			deep := 0
	// 			var x Astnode = ast
	// 			var parent = x.GetGNode().Parent
	// 			for parent != nil && parent != x {
	// 				deep++
	// 				x = parent
	// 				parent = x.GetGNode().Parent
	// 			}
	// 			return deep
	// 		}
	// 		sParentName := "parent: -"
	// 		if parent != nil && parent.GetGNode() != nil {
	// 			if parent.GetGNode().Name == "" {
	// 				sParentName = "parent: undefined"
	// 			} else {
	// 				sParentName = "parent: " + parent.GetGNode().Name
	// 			}
	// 		}
	// 		fmt.Println("DEEP " + helpers.PadLeft(sParentName, 34) + strconv.Itoa(depth(node)) + " Node " + Prefix(node) + node.ContentString())
	// 		return ""
	// 	},
	// })

	// Prepare all the nodes, child first.
	Walk(gm, nil, WalkPrepost{
		Post: func(node Ast, parent Ast) string {
			node.Prepare()
			return ""
		},
	})
	gm.wasInitialized = true
}

func (gm *Grammar) PrintRules() {
	fmt.Println("+ -- Grammar.Debug() --------")
	fmt.Println("| name         : " + gm.GNode.Name)
	fmt.Println("| label        : " + gm.GNode.Label)
	fmt.Println("| contentString: " + gm.ContentString())
	fmt.Println("| rules        : " + strconv.Itoa(gm.NumRules))
	fmt.Println("| ")
	fmt.Println("| ",
		helpers.PadLeft("key", 14),
		helpers.PadLeft("id", 3),
		helpers.PadLeft("type", 13),
		helpers.PadLeft("cap", 3),
		helpers.PadLeft("label", 7),
		helpers.PadLeft("labels()", 21),
		helpers.PadLeft("parent.name", 16),
		helpers.PadLeft("contentString", 30),
	)
	fmt.Println("|   -------------------------------------------------------------------------------------")
	for i := 0; i < gm.NumRules; i++ {
		v := gm.Id2Rule[i]
		sParentName := "-"
		if v.GetGNode().Parent != nil {
			sParentName = v.GetGNode().Parent.GetGNode().Name
		}
		fmt.Println("|  ",
			helpers.PadLeft(v.GetGNode().Name, 14),
			helpers.PadLeft(strconv.Itoa(v.GetGNode().Id), 3),
			helpers.PadLeft(helpers.TypeOfToString(v), 13),
			helpers.PadLeft(helpers.BoolToString(v.GetGNode().Capture), 3),
			helpers.PadLeft(v.GetGNode().Label, 7),
			helpers.PadLeft(strings.Join(v.GetGNode().Labels_.Get(), ","), 21),
			helpers.PadLeft(sParentName, 16),
			helpers.PadLeft(v.ContentString(), 30),
		)
	}
	fmt.Println("| ")
}
