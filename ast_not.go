package joeson

type not struct {
	*GNodeImpl
	it Parser
}

func newNot(it Ast) *not {
	gn := NewGNode()
	x := &not{gn, it.(Parser)}
	gn.capture = false
	gn.node = x
	return x
}

func (no *not) GetGNode() *GNodeImpl    { return no.GNodeImpl }
func (no *not) Prepare()                {}
func (no *not) HandlesChildLabel() bool { return false }

func (no *not) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Parser) Ast {
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

func (no *not) ContentString() string {
	return yellow("!") + String(no.it)
}
func (no *not) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	no.GetGNode().Rules = ForEachChild_InRules(no, f)
	if no.it != nil {
		no.it = f(no.it)
	}
	return no
}
