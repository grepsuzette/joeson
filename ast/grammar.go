package ast

import (
	"errors"
	"fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"reflect"
	"strconv"
)

/*
# Main external access.
# I dunno if Grammar should be a GNode or not. It
# might come in handy when embedding grammars
# in some glue language.
  jae 2012-04-18 in original joeson.coffee
*/
type Grammar struct {
	*GNode
	rank *Rank

	// each ast node can have rules, recursively.
	// in the Postinit below, Grammar will however
	// collect all children rules in its own GNode.Rules
	// and NumRules will be computed. Each rule node
	// will be given one incremental int Id.
	NumRules int

	// id2Rule: slow lookup for debugging...
	id2Rule map[int]Astnode // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605

	wasInitialized bool
}

func NewEmptyGrammarNamed(name string) *Grammar {
	gm := Grammar{NewGNode(), nil, 0, map[int]Astnode{}, false}
	gm.GNode.Name = name
	return &gm
}

// after Rank has already been set,
// collect and collect rules, simplify the rule tree etc.
func (gm *Grammar) Postinit() {
	if gm.rank == nil {
		panic("You can only call grammar.Postinit() after some rank has been set")
	}
	// from joeson.coffee
	// "TODO refactor into translation passes."

	// TODO Optimization that can wait  it's about:
	// Merge Choices with just a single choice.
	Walk(gm, nil, WalkPrepost{
		/*
					Pre: func(node Astnode, parent Astnode) Astnode {
						if choice, ok := node.(Choice); ok && len(choice.choices) == 1 {
							// Merge label
							if choice.choices[0].GetGNode().Label == "" {
								choice.choices[0].GetGNode().Label = choice.GetGNode().Label
							}
							// Merge included rules
							if len(choice.GetGNode().Rules) > 0 {
								for k, v := range choice.GetGNode().Rules {
									choice.choices[0].GetGNode().Rules[k] = v
								}
							}
							// Replace with grandchild
							// hum the key, index are not available in our implementation
							// TODO finish
							// ANOTHER WAY would be, when descending, if
							// child is sequence and
							// grandchildren are "choices", annd only 1, then
							// set it directly to grandchild[0]
							//
							// There is also a usage however in javascript.joeson
							// wait for now
							if index  > -1 {
								// TODO check hypothesis
								// hypothesis is if index is provided, it is an array
								reflect.ValueOf(parent).FieldByName(key).FieldByIndex(index)
							} else {
			   i					parent.GetGNode().SetArbitraryField(key, choice.Choice.choices[0])
								// tricky, needs & and Elem(). Also, field name must be
								// exported (capitalized first letter)
								// https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field

								reflect.ValueOf(&parent).Elem().FieldByName(key).se
							}
							parent
							reflect.S
						}
					}
		*/
	})

	// Connect all the nodes and collect dereferences into @rules
	Walk(gm, nil, WalkPrepost{
		Pre: func(node Astnode, parent Astnode) string {
			fmt.Println("grammar PRE: typeof node= " + reflect.TypeOf(node).String())
			gnode := node.GetGNode()
			// sanity check: it must have no parent yet if it's not a rule
			if !gnode.IsRule() && gnode.Parent != nil {
				panic("Grammar tree should be a DAG, nodes should not be referenced more than once.")
			}
			gnode.Grammar = gm
			gnode.Parent = parent
			if false {
				// TODO "inline rules are special" in original
			} else {
				// set node.rule, the root node for this rule
				if gnode.Rule == nil {
					fmt.Println("PRE looking if parent has rule")
					var r Astnode
					if parent != nil {
						r = parent.GetGNode().Rule
						fmt.Println("PRE  yes, parent has rule, returning it")
					} else {
						// must we return nil or undefined?
						fmt.Println("PRE  no, parent has no rule, returning NewNativeUndefined")
						r = NewNativeUndefined() // solution 1
						//r = nil // solution 2
					}
					gnode.Rule = r
				}
			}
			return ""
		},
		Post: func(node Astnode, parent Astnode) string {
			fmt.Println("grammar POST: typeof node= " + reflect.TypeOf(node).String() + " cs:" + node.ContentString())
			gnode := node.GetGNode()
			if gnode == nil {
				fmt.Println("ignoring gnode==nil type " + reflect.TypeOf(node).String())
			}
			if gnode.IsRule() {
				fmt.Println(Green("is rule!!!!!!!!!!!!!!!!!!!!!!"))
				gm.GetGNode().Rules[gnode.Name] = node
				gm.NumRules++
				gnode.Id = gm.NumRules
				gm.id2Rule[gnode.Id] = node
				if Trace.Loop { // print out id->rulename for convenience
					fmt.Println(Red(strconv.Itoa(gnode.Id)) + ":\t" + node.ContentString())
				}
			} else {
				fmt.Println("notta rul")
			}
			return ""
		},
	})
	// Prepare all the nodes, child first.
	Walk(gm, nil, WalkPrepost{
		Post: func(node Astnode, parent Astnode) string {
			node.Prepare()
			return ""
		},
	})
	gm.wasInitialized = true
}

// MAIN GRAMMAR PARSE FUNCTION
func (gm *Grammar) ParseString(sCode string, attrs ParseOptions) Astnode {
	return gm.ParseCode(NewCodeStream(sCode), attrs)
}

func (gm *Grammar) ParseCode(code *CodeStream, attrs ParseOptions) Astnode {
	return gm.Parse(NewParseContext(code, gm, attrs))
}

func (gm *Grammar) Parse(ctx *ParseContext) Astnode {
	var oldTrace TraceSettings
	if ctx.Debug {
		// temporarily enable stack tracing
		oldTrace = Trace
		Trace.Stack = true
	}
	// parse
	if gm.rank == nil {
		panic("Grammar.rank is nil")
	}
	result := gm.rank.Parse(ctx)
	// TODO
	// $.result = @rank.parse $
	// $.result?.code = code  <---- i dodoubt this line (joeson.coffee:625)
	// if ctx.Result != nil {
	// 	ctx.Result.code = code
	// }
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
		sErr += "\nDetails:\n"
		sErr += Green("OK") + "/"
		sErr += Yellow("Parsing") + "/"
		sErr += Red("Suspect") + "/"
		sErr += White("Unknown") + "\n\n"
		sErr += Green(ctx.Code.Peek(Peek{BeforeLines: helpers.NewNullInt(2)}))
		sErr += Yellow(ctx.Code.Peek(Peek{AfterChars: helpers.NewNullInt(maxSuccess - ctx.Code.Pos)}))
		ctx.Code.Pos = maxSuccess
		sErr += Red(ctx.Code.Peek(Peek{AfterChars: helpers.NewNullInt(maxAttempt - ctx.Code.Pos)})) + "/"
		ctx.Code.Pos = maxAttempt
		sErr += White(ctx.Code.Peek(Peek{AfterLines: helpers.NewNullInt(2)})) + "\n"
		panic(errors.New(sErr))
	}
	// joeson.coffee has a opts.returnContext but won't implement it
	return result
}

func (gm *Grammar) Prepare()                {}
func (gm *Grammar) HandlesChildLabel() bool { return false }
func (gm *Grammar) Labels() []string        { return gm.GNode.Labels() }
func (gm *Grammar) Captures() []Astnode     { return gm.GNode.Captures() }
func (gm *Grammar) ContentString() string {
	return Magenta("GRAMMAR{") +
		ShowLabelOrNameIfAny(gm) +
		gm.rank.ContentString() + Magenta("}")
}

// satisfy GrammarRuleCounter
func (gm *Grammar) CountRules() int { return gm.NumRules }
func (gm *Grammar) IsReady() bool   { return gm.rank != nil && gm.wasInitialized }
func (gm *Grammar) SetRankIfEmpty(rank *Rank) {
	if gm.rank != nil {
		return
	}
	if gm.IsReady() {
		panic("Grammar is already defined and can not be changed on the fly at the moment")
	}
	gm.rank = rank
}
func (gm *Grammar) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren rank: {type:Rank}
	if gm.rank != nil {
		gm.rank = f(gm.rank).(*Rank)
	}
	return gm
}
