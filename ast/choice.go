package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/lambda"
	"strings"
)

type Choice struct {
	*GNode
	a []Astnode
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
	if len(ch.a) == 1 {
		return true, ch.a[0]
	} else {
		return false, nil
	}
}

func (ch *Choice) Append(node Astnode)     { ch.a = append(ch.a, node) }
func (ch *Choice) GetGNode() *GNode        { return ch.GNode }
func (ch *Choice) HandlesChildLabel() bool { return false }
func (ch *Choice) Labels() []string        { return ch.GNode.Labels() }
func (ch *Choice) Captures() []Astnode     { return ch.GNode.Captures() }
func (ch *Choice) Prepare() {
	for _, choice := range ch.a {
		if !choice.GetGNode().Capture {
			ch.GNode.Capture = false
			return
		}
	}
	ch.GNode.Capture = true
}

func (ch *Choice) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
		for _, choice := range ch.a {
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
	b.WriteString(ShowLabelOrNameIfAny(ch))
	b.WriteString(Blue("("))
	a := lambda.Map(ch.a, func(nn Astnode) string {
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
	ch.a = ForEachChild_Array(ch.a, f)
	ch.GetGNode().Rules = ForEachChild_MapString(ch.GetGNode().Rules, f)
	return ch
}
func (ch *Choice) String() string { return "AFAKFAKg" }
