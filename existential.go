package main

import "grepsuzette/joeson/helpers"

type Existential struct {
	GNode
	it        astnode
	_labels   helpers.Varcache[[]string]  // internal cache for Labels()
	_captures helpers.Varcache[[]astnode] // internal cache for Captures()
}

func NewExistential(it astnode) {
	ex := Existential{NewGNode(), it}
	return ex
}

func (ex Existential) HandlesChildLabel() bool {
	if ex.GNode.parent != nil {
		return ex.GNode.parent.HandlesChildLabel()
	} else {
		return false
	}
}

func (ex Existential) Labels() []string    { return ex._labels.GetCache() }
func (ex Existential) Captures() []astnode { return ex._captures.GetCache() }

func (ex *Existential) Prepare() {
	ex._labels.GetCacheOrSet(func() []string {
		labels := []string{}
		if ex.GNode.label != "" && ex.GNode.label != "@" && ex.GNode.label != "&" {
			labels = []string{ex.GNode.label}
		} else {
			labels = ex.it.Labels()
		}
		return labels
	})
	ex._captures.GetCacheOrSet(func() []astnode {
		caps := ex.it.Captures()
		ex.GNode.capture = len(caps) > 0
		return caps
	})
}

// TODO parse, contentString
