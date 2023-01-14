package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
)

type Existential struct {
	*GNode
	it Ast
	// moved to GNode
	// _labels   helpers.Lazy[[]string]  // internal cache for Labels()
	// _captures helpers.Lazy[[]Astnode] // internal cache for Captures()
}

func NewExistential(it Ast) *Existential {
	ex := &Existential{GNode: NewGNode(), it: it}
	ex.GNode.Node = ex
	return ex
}

// TODO handlesChildLabel$: get: -> @parent?.handlesChildLabel
// examine this case^
func (ex *Existential) HandlesChildLabel() bool {
	if ex.GNode.Parent != nil {
		return ex.GNode.Parent.HandlesChildLabel()
	} else {
		return false
	}
}

func (ex *Existential) GetGNode() *GNode { return ex.GNode }

func (ex *Existential) Labels() []string { panic("z") }
func (ex *Existential) Captures() []Ast  { panic("z") }

func (ex *Existential) Prepare() {
	gn := ex.GetGNode()
	var lbls = ex.calculateLabels()
	gn.Labels_.Set(lbls)
	if len(lbls) > 0 && gn.Label == "" {
		gn.Label = "@"
	}
	var caps = ex.it.GetGNode().Captures_.Get()
	gn.Captures_.Set(caps)
	gn.Capture = len(caps) > 0
}

func (ex *Existential) calculateLabels() []string {
	gn := ex.GetGNode()
	lbl := gn.Label
	if lbl != "" && lbl != "@" && lbl != "&" {
		return []string{lbl}
	} else {
		return ex.it.GetGNode().Labels_.Get()
	}
}

func (ex *Existential) ContentString() string {
	return Prefix(ex.it) + ex.it.ContentString() + Blue("?")
}

func (ex *Existential) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
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
func (ex *Existential) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	ex.GetGNode().Rules = ForEachChild_InRules(ex, f)
	if ex.it != nil {
		ex.it = f(ex.it)
	}
	return ex
}
