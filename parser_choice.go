package joeson

import (
	"strings"
)

type choice struct {
	*Attr
	*gnodeimpl
	choices []Parser
}

func newEmptyChoice() *choice {
	ch := &choice{newAttr(), newGNode(), []Parser{}}
	ch.gnodeimpl.node = ch
	return ch
}

func newChoice(it Ast) *choice {
	if a, ok := it.(*NativeArray); ok {
		var parsers []Parser
		for _, ast := range a.Array {
			parsers = append(parsers, ast.(Parser))
		}
		ch := &choice{
			newAttr(),
			newGNode(),
			parsers,
		}
		ch.gnodeimpl.node = ch
		return ch
	} else {
		panic("Choice expects a NativeArray")
	}
}

func (ch *choice) isMonoChoice() bool      { return len(ch.choices) == 1 }
func (ch *choice) Append(node Parser)      { ch.choices = append(ch.choices, node) }
func (ch *choice) gnode() *gnodeimpl       { return ch.gnodeimpl }
func (ch *choice) HandlesChildLabel() bool { return false }

func (ch *choice) prepare() {
	for _, choice := range ch.choices {
		if !choice.Capture() {
			ch.SetCapture(false)
			return
		}
	}
	ch.SetCapture(true)
}

func (ch *choice) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		for _, choice := range ch.choices {
			pos := ctx.Code.Pos()
			result := choice.Parse(ctx)
			if result == nil {
				ctx.Code.SetPos(pos)
			} else {
				return result
			}
		}
		return nil
	}, ch)(ctx)
}

func (ch *choice) String() string {
	var b strings.Builder
	b.WriteString(Blue("("))
	first := true
	for _, x := range ch.choices {
		if !first {
			b.WriteString(Blue(" | "))
		}
		b.WriteString(String(x))
		first = false
	}
	b.WriteString(Blue(")"))
	return b.String()
}

func (ch *choice) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	// we must first walk through rules, and then only through choices
	ch.rules = ForEachChildInRules(ch, f)
	ch.choices = ForEachChild_Array(ch.choices, f)
	return ch
}
