package ast

import (
	"fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/lambda"
	"strings"
)

type Choice struct {
	*GNode
	choices []Astnode
}

func NewEmptyChoice() *Choice { return &Choice{NewGNode(), []Astnode{}} }
func NewChoice(it Astnode) *Choice {
	if a, ok := it.(*NativeArray); ok {
		return &Choice{NewGNode(), a.Array}
	} else {
		panic("Choice expects a NativeArray")
	}
}

func (ch *Choice) IsThereOnlyOneChoice() (bool, Astnode) {
	if len(ch.choices) == 1 {
		return true, ch.choices[0]
	} else {
		return false, nil
	}
}

func (ch *Choice) Append(node Astnode)     { ch.choices = append(ch.choices, node) }
func (ch *Choice) GetGNode() *GNode        { return ch.GNode }
func (ch *Choice) HandlesChildLabel() bool { return false }
func (ch *Choice) Labels() []string        { return MyLabelIfDefinedOrEmpty(ch) }
func (ch *Choice) Captures() []Astnode     { return MeIfCaptureOrEmpty(ch) }
func (ch *Choice) Prepare() {
	for _, choice := range ch.choices {
		if !choice.GetGNode().Capture {
			ch.GNode.Capture = false
			return
		}
	}
	ch.GNode.Capture = true
}

func (ch *Choice) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
		for i, choice := range ch.choices {
			fmt.Printf("Choice n=%d %s\n", i, choice.ContentString())
			pos := ctx.Code.Pos
			result := choice.Parse(ctx)
			if result == nil {
				ctx.Code.Pos = pos
			} else {
				return result
			}
		}
		return nil
	}, ch)(ctx)
}

func (ch *Choice) ContentString() string {
	var b strings.Builder
	b.WriteString(LabelOrName(ch))
	b.WriteString(Blue("("))
	a := lambda.Map(ch.choices, func(nn Astnode) string {
		return nn.ContentString()
	})
	b.WriteString(strings.Join(a, Blue(" | ")))
	b.WriteString(Blue(")"))
	return b.String()
}

func (ch *Choice) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	// we must first wall through rules, and then only through choices
	ch.GetGNode().Rules = ForEachChild_MapString(ch.GetGNode().Rules, f)
	ch.choices = ForEachChild_Array(ch.choices, f)
	return ch
}
func (ch *Choice) String() string { return "AFAKFAKg" }
