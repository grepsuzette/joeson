package joeson

// Lookahead does not consume or capture characters.
// It is meant as an optimizer. See test file for an example.
//
// (?foo) -> lookahead{foo}
// ?foo -> lookahead{foo}
// foo? -> existential{foo}
//
// lookahead does not capture.
// use parens when using this with alternation
// e.g. "a ?c | b" is grouped as "a (?c | b)".
type lookahead struct {
	*Attr
	*rule
	expr Parser
}

func newLookahead(it Ast) *lookahead {
	gn := newRule()
	la := &lookahead{newAttr(), gn, it.(*NativeMap).GetOrPanic("expr").(Parser)}
	gn.capture = false
	gn.node = la
	return la
}

func (look *lookahead) prepare()                {}
func (look *lookahead) gnode() *rule            { return look.rule }
func (look *lookahead) HandlesChildLabel() bool { return false }
func (look *lookahead) String() string {
	return Blue("(?") + String(look.expr) + Blue(")")
}

func (look *lookahead) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos()
		result := look.expr.Parse(ctx) // check whether it parses...
		ctx.Code.SetPos(pos)           // ...but revert to prev pos if so
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
