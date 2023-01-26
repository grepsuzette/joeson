package joeson

import (
	"grepsuzette/joeson/helpers"
	"strings"
)

type choice struct {
	*GNode
	choices []Ast
}

func newEmptyChoice() *choice {
	ch := &choice{NewGNode(), []Ast{}}
	ch.GNode.Node = ch
	return ch
}
func newChoice(it Ast) *choice {
	if a, ok := it.(*NativeArray); ok {
		ch := &choice{NewGNode(), a.Array}
		ch.GNode.Node = ch
		return ch
	} else {
		panic("Choice expects a NativeArray")
	}
}

func (ch *choice) isMonoChoice() bool      { return len(ch.choices) == 1 }
func (ch *choice) Append(node Ast)         { ch.choices = append(ch.choices, node) }
func (ch *choice) GetGNode() *GNode        { return ch.GNode }
func (ch *choice) HandlesChildLabel() bool { return false }

func (ch *choice) Prepare() {
	for _, choice := range ch.choices {
		if !choice.GetGNode().Capture {
			ch.GNode.Capture = false
			return
		}
	}
	ch.GNode.Capture = true
}

func (ch *choice) Parse(ctx *ParseContext) Ast {
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

func (ch *choice) ContentString() string {
	var b strings.Builder
	b.WriteString(blue("("))
	a := helpers.AMap(ch.choices, func(x Ast) string { return String(x) })
	b.WriteString(strings.Join(a, blue(" | ")))
	b.WriteString(blue(")"))
	return b.String()
}

func (ch *choice) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	// we must first walk through rules, and then only through choices
	ch.GetGNode().Rules = ForEachChild_InRules(ch, f)
	ch.choices = ForEachChild_Array(ch.choices, f)
	return ch
}
