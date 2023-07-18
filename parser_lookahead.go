package joeson

type lookahead struct {
	*gnodeimpl
	expr Parser
}

func newLookahead(it Ast) *lookahead {
	gn := NewGNode()
	la := &lookahead{gn, newStrFromAst(it)}
	gn.capture = false
	gn.node = la
	return la
}

func (look *lookahead) prepare()                {}
func (look *lookahead) gnode() *gnodeimpl       { return look.gnodeimpl }
func (look *lookahead) handlesChildLabel() bool { return false }
func (look *lookahead) ContentString() string {
	return blue("(?") + String(look.expr) + blue(")")
}

func (look *lookahead) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos
		result := look.expr.Parse(ctx) // check whether it parses...
		ctx.Code.Pos = pos             // ...but revert to prev pos if so
		return result
	}, look)(ctx)
}

func (look *lookahead) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   expr:       {type:GNode}
	look.rules = ForEachChildInRules(look, f)
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	return look
}
