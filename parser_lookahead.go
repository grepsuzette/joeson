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
	rule := newRule()
	la := &lookahead{newAttr(), rule, it.(*NativeMap).GetOrPanic("expr").(Parser)}
	rule.capture = false
	rule.node = la
	return la
}

func (look *lookahead) prepare()                {}
func (look *lookahead) getRule() *rule          { return look.rule }
func (look *lookahead) handlesChildLabel() bool { return false }
func (look *lookahead) String() string {
	return Blue("(?") + String(look.expr) + Blue(")")
}

func (look *lookahead) parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos()
		result := look.expr.parse(ctx) // check whether it parses...
		ctx.Code.SetPos(pos)           // ...but revert to prev pos if so
		return result
	}, look)(ctx)
}

func (look *lookahead) forEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   expr:       {type:GNode}
	look.rules = ForEachChildInRules(look, f)
	if look.expr != nil {
		look.expr = f(look.expr)
	}
	return look
}
