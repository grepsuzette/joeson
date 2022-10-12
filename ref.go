package main

import . "grepsuzette/joeson/colors"
import "grepsuzette/joeson/helpers"

type Ref struct {
	GNode
	ref     string                     // ref because joeson.coffee used @ref, because @name was reserved
	param   any                        // see frame.go (`param` field) and joeson.coffee:67 TODO is it any?
	_labels helpers.Varcache[[]string] // internal cache for labels()
}

func NewRef(name string, param any) Ref {
	ref := Ref{newGNode(), name, param}
	if name[0:1] == "_" {
		ref.GNode.capture = false
	}
	return ref
}

func (ref Ref) GetGNode() GNode    { return ref.GNode }
func (ref Ref) HandlesChildLabel() { return false }
func (ref Ref) Prepare()           {}
func (ref Ref) Parse(ctx *ParseContext) astnode {
	return ref.GNode._wrap(func(_, _) astnode {
		var x astnode = ref.GNode.grammar.rules[ref.ref]
		if x == nil {
			panic("Unknown reference " + ref.ref)
		}
		ctx.stackPeek(0).param = ref.param
		return x.Parse(ctx)
	})(ref, ctx)
}

func (ref Ref) Captures() []astnode { return ref.GNode.Captures() }
func (ref Ref) Labels() []string {
	return ref._labels.GetCacheOrSet(func() []string {
		if ref.GNode.label == "@" {
			return ref.GNode.grammar.rules[ref.ref].Labels()
		} else if ref.GNode.label != "" {
			return []string{ref.GNode.label}
		} else {
			return []string{}
		}
	})
}

func (ref Ref) ContentString() string {
	return Red(ref.ref)
}
