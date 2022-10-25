package ast

import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"

type Lookahead struct {
	*GNode
	expr Astnode
}

func NewLookahead(it Astnode) *Lookahead {
	g := NewGNode()
	g.Capture = false
	return &Lookahead{g, it}

}

func (look *Lookahead) Prepare()                {}
func (look *Lookahead) HandlesChildLabel() bool { return false }
func (look *Lookahead) Labels() []string        { return look.GNode.Labels() }
func (look *Lookahead) Captures() []Astnode     { return look.GNode.Captures() }
func (look *Lookahead) ContentString() string {
	return ShowLabelOrNameIfAny(look) + Blue("(?") + look.expr.ContentString() + Blue(")")
}
func (look *Lookahead) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
		pos := ctx.Code.Pos
		result := look.expr.Parse(ctx) // check whether it parses
		ctx.Code.Pos = pos             // but revert pos
		return result
	}, look)(ctx)
}
func (look *Lookahead) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   expr:       {type:GNode}
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	look.GetGNode().Rules = ForEachChild_MapString(look.GetGNode().Rules, f)
	return look
}
