package joeson

type not struct {
	*Attributes
	*gnodeimpl
	it Parser
}

func newNot(it Ast) *not {
	gn := NewGNode()
	x := &not{&Attributes{}, gn, it.(Parser)}
	gn.capture = false
	gn.node = x
	return x
}

func (no *not) gnode() *gnodeimpl       { return no.gnodeimpl }
func (no *not) prepare()                {}
func (no *not) handlesChildLabel() bool { return false }

func (no *not) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos
		res := no.it.Parse(ctx)
		ctx.Code.Pos = pos
		if res != nil {
			return nil
		} else {
			return NewNativeUndefined()
		}
	}, no)(ctx)
}

func (no *not) String() string {
	return yellow("!") + String(no.it)
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
