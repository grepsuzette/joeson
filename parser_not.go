package joeson

// Negative lookahead
//
// !foo -> not{foo}
// !(foo) -> not{foo}
//
// it does not capture.
// use parens when using this with alternation
// e.g. "a !c | b" is grouped as "a (!c | b)".
type not struct {
	*Attr
	*rule
	it Parser
}

func newNot(it Ast) *not {
	gn := newRule()
	x := &not{newAttr(), gn, it.(Parser)}
	gn.capture = false
	gn.node = x
	return x
}

func (no *not) gnode() *rule            { return no.rule }
func (no *not) prepare()                {}
func (no *not) HandlesChildLabel() bool { return false }

func (no *not) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos()
		res := no.it.Parse(ctx)
		ctx.Code.SetPos(pos)
		if res != nil {
			return nil
		} else {
			return NewNativeUndefined()
		}
	}, no)(ctx)
}

func (no *not) String() string {
	return Yellow("!") + String(no.it)
}

func (no *not) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	no.rules = ForEachChildInRules(no, f)
	if no.it != nil {
		no.it = f(no.it)
	}
	return no
}
