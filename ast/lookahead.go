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
func (look *Lookahead) GetGNode() *GNode        { return look.GNode }
func (look *Lookahead) HandlesChildLabel() bool { return false }
func (look *Lookahead) Labels() []string        { return MyLabelIfDefinedOrEmpty(look) }
func (look *Lookahead) Captures() []Astnode     { return MeIfCaptureOrEmpty(look) }
func (look *Lookahead) ContentString() string {
	return LabelOrName(look) + Blue("(?") + look.expr.ContentString() + Blue(")")
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
	look.GetGNode().Rules = ForEachChild_MapString(look.GetGNode().Rules, f)
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	return look
}
