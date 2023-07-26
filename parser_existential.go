package joeson

type existential struct {
	Attr
	*gnodeimpl
	it Parser
}

func newExistential(it Ast) *existential {
	ex := &existential{Attr: newAttr(), gnodeimpl: newGNode(), it: it.(Parser)}
	ex.gnodeimpl.node = ex
	return ex
}

// TODO handlesChildLabel$: get: -> @parent?.handlesChildLabel
// examine this case^
func (ex *existential) handlesChildLabel() bool {
	if ex.gnodeimpl.parent != nil {
		return ex.gnodeimpl.parent.handlesChildLabel()
	} else {
		return false
	}
}

func (ex *existential) gnode() *gnodeimpl { return ex.gnodeimpl }

func (ex *existential) prepare() {
	lbls := ex.calculateLabels()
	ex.cachedLabels.Set(lbls)
	if len(lbls) > 0 && ex.label == "" {
		ex.label = "@"
	}
	caps := ex.it.gnode().cachedCaptures.Get()
	ex.cachedCaptures.Set(caps)
	ex.capture = len(caps) > 0
}

func (ex *existential) calculateLabels() []string {
	lbl := ex.label
	if lbl != "" && lbl != "@" && lbl != "&" {
		return []string{lbl}
	} else {
		return ex.it.gnode().cachedLabels.Get()
	}
}

func (ex *existential) String() string {
	return String(ex.it) + Blue("?")
}

func (ex *existential) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
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
	ex.rules = ForEachChildInRules(ex, f)
	if ex.it != nil {
		ex.it = f(ex.it)
	}
	return ex
}
