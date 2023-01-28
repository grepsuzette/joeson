package joeson

type not struct {
	*GNode
	it Ast
}

func newNot(it Ast) *not {
	gn := NewGNode()
	x := &not{gn, it}
	gn.Capture = false
	gn.Node = x
	return x
}

func (no *not) GetGNode() *GNode        { return no.GNode }
func (no *not) Prepare()                {}
func (no *not) HandlesChildLabel() bool { return false }

func (no *not) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
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
func (no *not) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	no.GetGNode().Rules = ForEachChild_InRules(no, f)
	if no.it != nil {
		no.it = f(no.it)
	}
	return no
}