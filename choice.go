package main

import . "grepsuzette/joeson/colors"
import "grepsuzette/joeson/lambda"
import "strings"

type Choice struct {
	GNode
	choices []astnode
}

func NewEmptyChoice() Choice {
	ch := Choice{newGNode(), []astnode{}}
	return ch
}

func (ch Choice) GetGNode() GNode     { return ch.GNode }
func (ch Choice) HandlesChildLabel()  { return false }
func (ch Choice) Labels() []string    { return ch.GNode.Labels() }
func (ch Choice) Captures() []astnode { return ch.GNode.Captures() }
func (ch Choice) Prepare() {
	for _, choice := range ch.choices {
		if !choice.GetGNode().capture {
			ch.GNode.capture = false
			return
		}
	}
	ch.GNode.capture = true
}

func (ch Choice) Parse(ctx *ParseContext) astnode {
	// return ch.GNode._wrap(func(_, _) astnode {
	return _wrap(func(ctx, _) astnode {
		for _, choice := range ch.choices {
			pos := ctx.code.pos
			result := choice.Parse(ctx)
			if result == nil {
				ctx.code.pos = pos
			} else {
				return result
			}
		}
		return nil
	})(ctx, ch)
}

func (ch Choice) ContentString() string {
	var b strings.Builder
	b.WriteString(Blue("("))
	a := lambda.Map(ch.choices, func(nn astnode) {
		return nn.ContentString()
	})
	b.WriteString(strings.Join(a, " | "))
	b.WriteString(Blue(")"))
	return b.String()
}
