package joeson

// foo? -> existential{foo}
// (?foo) -> lookahead{foo}
// ?foo -> lookahead{foo}
type existential struct {
	*Attr
	*rule
	it Parser
}

func newExistential(it Ast) *existential {
	ex := &existential{Attr: newAttr(), rule: newRule(), it: it.(Parser)}
	ex.rule.node = ex
	return ex
}

func (ex *existential) handlesChildLabel() bool {
	if ex.rule.parent != nil {
		return ex.rule.parent.handlesChildLabel()
	} else {
		return false
	}
}

func (ex *existential) getRule() *rule { return ex.rule }

func (ex *existential) prepare() {
	lbls := ex.calculateLabels()
	ex.labels_.Set(lbls)
	if len(lbls) > 0 && ex.label == "" {
		ex.label = "@"
	}
	caps := ex.it.getRule().captures_.Get()
	ex.captures_.Set(caps)
	ex.capture = len(caps) > 0
}

func (ex *existential) calculateLabels() []string {
	lbl := ex.label
	if lbl != "" && lbl != "@" && lbl != "&" {
		return []string{lbl}
	} else {
		return ex.it.getRule().labels_.Get()
	}
}

func (ex *existential) String() string {
	return String(ex.it) + Blue("?")
}

func (ex *existential) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		pos := ctx.Code.Pos()
		result := ex.it.Parse(ctx)
		if result == nil {
			ctx.Code.SetPos(pos)
			return NewNativeUndefined()
		} else {
			return result
		}
	}, ex)(ctx)
}

func (ex *existential) forEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	ex.rules = ForEachChildInRules(ex, f)
	if ex.it != nil {
		ex.it = f(ex.it)
	}
	return ex
}
