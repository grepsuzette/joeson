package joeson

import (
	"grepsuzette/joeson/helpers"
	"strings"
)

type Choice struct {
	*GNode
	choices []Ast
}

func NewEmptyChoice() *Choice {
	ch := &Choice{NewGNode(), []Ast{}}
	ch.GNode.Node = ch
	return ch
}
func NewChoice(it Ast) *Choice {
	if a, ok := it.(*NativeArray); ok {
		ch := &Choice{NewGNode(), a.Array}
		ch.GNode.Node = ch
		return ch
	} else {
		panic("Choice expects a NativeArray")
	}
}

func (ch *Choice) IsMonoChoice() bool      { return len(ch.choices) == 1 }
func (ch *Choice) Append(node Ast)         { ch.choices = append(ch.choices, node) }
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

func (ch *Choice) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
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
	b.WriteString(blue("("))
	a := helpers.AMap(ch.choices, func(x Ast) string { return String(x) })
	b.WriteString(strings.Join(a, blue(" | ")))
	b.WriteString(blue(")"))
	return b.String()
}

func (ch *Choice) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	// we must first walk through rules, and then only through choices
	ch.GetGNode().Rules = ForEachChild_InRules(ch, f)
	ch.choices = ForEachChild_Array(ch.choices, f)
	return ch
}
