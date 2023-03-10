package joeson

import (
	"github.com/grepsuzette/joeson/helpers"
	"strings"
)

type choice struct {
	*GNodeImpl
	choices []Parser
}

func newEmptyChoice() *choice {
	ch := &choice{NewGNode(), []Parser{}}
	ch.GNodeImpl.node = ch
	return ch
}
func newChoice(it Ast) *choice {
	if a, ok := it.(*NativeArray); ok {
		ch := &choice{NewGNode(), helpers.AMap(a.Array, func(ast Ast) Parser { return ast.(Parser) })}
		ch.GNodeImpl.node = ch
		return ch
	} else {
		panic("Choice expects a NativeArray")
	}
}

func (ch *choice) isMonoChoice() bool      { return len(ch.choices) == 1 }
func (ch *choice) Append(node Parser)      { ch.choices = append(ch.choices, node) }
func (ch *choice) GetGNode() *GNodeImpl    { return ch.GNodeImpl }
func (ch *choice) HandlesChildLabel() bool { return false }

func (ch *choice) Prepare() {
	for _, choice := range ch.choices {
		if !choice.Capture() {
			ch.SetCapture(false)
			return
		}
	}
	ch.SetCapture(true)
}

func (ch *choice) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Parser) Ast {
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
	a := helpers.AMap(ch.choices, func(x Parser) string { return String(x) })
	b.WriteString(strings.Join(a, blue(" | ")))
	b.WriteString(blue(")"))
	return b.String()
}

func (ch *choice) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	// we must first walk through rules, and then only through choices
	ch.GetGNode().rules = ForEachChild_InRules(ch, f)
	ch.choices = ForEachChild_ArrayParser(ch.choices, f)
	return ch
}
