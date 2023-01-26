package joeson

type lookahead struct {
	*GNode
	expr Ast
}

func newLookahead(it Ast) *lookahead {
	gn := NewGNode()
	la := &lookahead{gn, it}
	gn.Capture = false
	gn.Node = la
	return la
}

func (look *lookahead) Prepare()                {}
func (look *lookahead) GetGNode() *GNode        { return look.GNode }
func (look *lookahead) HandlesChildLabel() bool { return false }
func (look *lookahead) ContentString() string {
	return blue("(?") + String(look.expr) + blue(")")
}
func (look *lookahead) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		pos := ctx.Code.Pos
		result := look.expr.Parse(ctx) // check whether it parses...
		ctx.Code.Pos = pos             // ...but revert to prev pos if so
		return result
	}, look)(ctx)
}
func (look *lookahead) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   expr:       {type:GNode}
	look.GetGNode().Rules = ForEachChild_InRules(look, f)
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	return look
}
