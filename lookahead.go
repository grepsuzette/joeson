package main

type Lookahead struct {
	GNode
	expr astnode
}

func NewLookahead(expr Expr) {
	g := NewGNode()
	g.capture = false
	look := Lookahead{g, expr}
}

func (look Lookahead) Prepare()            {}
func (look Lookahead) HandlesChildLabel()  { return false }
func (look Lookahead) Labels() []string    { return look.GNode.Labels() }
func (look Lookahead) Captures() []astnode { return look.GNode.Captures() }
func (look Lookahead) ContentString() string {
	return Blue("(?") + look.expr.ContentString() + Blue(")")
}
func (look Lookahead) Parse(ctx *ParseContext) astnode {
	return _wrap(func(ctx, _) astnode {
		pos := ctx.code.pos
		result := look.expr.Parse(ctx)
		ctx.code.pos = pos
		return result
	})(ctx, look)
}
