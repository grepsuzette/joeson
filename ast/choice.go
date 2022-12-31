package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/lambda"
	"strings"
)

type Choice struct {
	*GNode
	choices []Astnode
}

func NewEmptyChoice() *Choice {
	ch := &Choice{NewGNode(), []Astnode{}}
	ch.GNode.Node = ch
	return ch
}
func NewChoice(it Astnode) *Choice {
	if a, ok := it.(*NativeArray); ok {
		ch := &Choice{NewGNode(), a.Array}
		ch.GNode.Node = ch
		return ch
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
		for _, choice := range ch.choices {
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
	b.WriteString(Blue("("))
	a := lambda.Map(ch.choices, func(x Astnode) string {
		return Prefix(x) + x.ContentString()
	})
	b.WriteString(strings.Join(a, Blue(" | ")))
	b.WriteString(Blue(")"))
	return b.String()
}

func (ch *Choice) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	// we must first walk through rules, and then only through choices
	ch.GetGNode().Rules = ForEachChild_InRules(ch, f)
	ch.choices = ForEachChild_Array(ch.choices, f)
	return ch
}
