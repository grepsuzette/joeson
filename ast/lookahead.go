package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
)

type Lookahead struct {
	*GNode
	expr Ast
}

func NewLookahead(it Ast) *Lookahead {
	gn := NewGNode()
	la := &Lookahead{gn, it}
	gn.Capture = false
	gn.Node = la
	return la
}

func (look *Lookahead) Prepare()                {}
func (look *Lookahead) GetGNode() *GNode        { return look.GNode }
func (look *Lookahead) HandlesChildLabel() bool { return false }
func (look *Lookahead) ContentString() string {
	return Blue("(?") + String(look.expr) + Blue(")")
}
func (look *Lookahead) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		pos := ctx.Code.Pos
		result := look.expr.Parse(ctx) // check whether it parses...
		ctx.Code.Pos = pos             // ...but revert to prev pos if so
		return result
	}, look)(ctx)
}
func (look *Lookahead) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   expr:       {type:GNode}
	look.GetGNode().Rules = ForEachChild_InRules(look, f)
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	return look
}
