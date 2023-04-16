package joeson

type existential struct {
	*GNodeImpl
	it Parser
}

func newExistential(it Ast) *existential {
	ex := &existential{GNodeImpl: NewGNode(), it: it.(Parser)}
	ex.GNodeImpl.node = ex
	return ex
}

// TODO handlesChildLabel$: get: -> @parent?.handlesChildLabel
// examine this case^
func (ex *existential) HandlesChildLabel() bool {
	if ex.GNodeImpl.parent != nil {
		return ex.GNodeImpl.parent.HandlesChildLabel()
	} else {
		return false
	}
}

func (ex *existential) GetGNode() *GNodeImpl { return ex.GNodeImpl }

func (ex *existential) Prepare() {
	gn := ex.GetGNode()
	var lbls = ex.calculateLabels()
	gn.labels_.Set(lbls)
	if len(lbls) > 0 && gn.label == "" {
		gn.label = "@"
	}
	var caps = ex.it.GetGNode().captures_.Get()
	gn.captures_.Set(caps)
	gn.capture = len(caps) > 0
}

func (ex *existential) calculateLabels() []string {
	gn := ex.GetGNode()
	lbl := gn.label
	if lbl != "" && lbl != "@" && lbl != "&" {
		return []string{lbl}
	} else {
		return ex.it.GetGNode().labels_.Get()
	}
}

func (ex *existential) ContentString() string {
	return String(ex.it) + blue("?")
}

func (ex *existential) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos
		result := ex.it.Parse(ctx)
		if result == nil {
			ctx.Code.Pos = pos
			return NewNativeUndefined()
		} else {
			return result
		}
	}, ex)(ctx)
}
func (ex *existential) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	ex.GetGNode().rules = ForEachChild_InRules(ex, f)
	if ex.it != nil {
		ex.it = f(ex.it)
	}
	return ex
}
