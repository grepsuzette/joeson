package ast

import "grepsuzette/joeson/helpers"
import . "grepsuzette/joeson/colors"
import . "grepsuzette/joeson/core"

type Existential struct {
	*GNode
	it        Astnode
	_labels   helpers.Varcache[[]string]  // internal cache for Labels()
	_captures helpers.Varcache[[]Astnode] // internal cache for Captures()
}

func NewExistential(it Astnode) *Existential {
	ex := Existential{GNode: NewGNode(), it: it}
	return &ex
}

func (ex *Existential) HandlesChildLabel() bool {
	if ex.GNode.Parent != nil {
		return ex.GNode.Parent.HandlesChildLabel()
	} else {
		return false
	}
}

// their cache have been written in Prepare()
func (ex *Existential) Labels() []string    { return ex._labels.GetCache() }
func (ex *Existential) Captures() []Astnode { return ex._captures.GetCache() }

func (ex *Existential) Prepare() {
	var labels []string
	if ex.GNode.Label != "" && ex.GNode.Label != "@" && ex.GNode.Label != "&" {
		labels = []string{ex.GNode.Label}
	} else {
		labels = ex.it.Labels()
	}
	if len(labels) > 0 {
		ex.GNode.Label = "@"
	}
	captures := ex.it.Captures()
	ex.GNode.Capture = captures != nil && len(captures) > 0
	ex._labels.SetCache(labels)
	ex._captures.SetCache(captures)
}

func (ex *Existential) ContentString() string {
	return ShowLabelOrNameIfAny(ex) + ex.it.ContentString() + Blue("?")
}

func (ex *Existential) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
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
func (ex *Existential) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	if ex.it != nil {
		ex.it = f(ex.it)
	}
	ex.GetGNode().Rules = ForEachChild_MapString(ex.GetGNode().Rules, f)
	return ex
}
