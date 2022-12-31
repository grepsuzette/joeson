package ast

import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"

type Lookahead struct {
	*GNode
	expr Astnode
}

func NewLookahead(it Astnode) *Lookahead {
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
	return Blue("(?") + Prefix(look.expr) + look.expr.ContentString() + Blue(")")
}
func (look *Lookahead) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
		pos := ctx.Code.Pos
		result := look.expr.Parse(ctx) // check whether it parses...
		ctx.Code.Pos = pos             // ...but revert to prev pos if so
		return result
	}, look)(ctx)
}
func (look *Lookahead) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   expr:       {type:GNode}
	look.GetGNode().Rules = ForEachChild_InRules(look, f)
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	return look
}
