package ast

import (
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/helpers"

	. "grepsuzette/joeson/core"
)

type Ref struct {
	*GNode
	ref     string                     // ref because joeson.coffee used @ref, because @name was reserved
	param   Astnode                    // thought it was `any`, see frame.go (`param` field) and joeson.coffee:67. But Astnode must be good
	_labels helpers.Varcache[[]string] // internal cache for labels()
}

func NewRef(it Astnode) *Ref {
	// it -> name string, param any
	var name string
	var param Astnode = nil
	switch v := it.(type) {
	case NativeString:
		var ns NativeString = v
		name = ns.Str
	case *NativeArray:
		var na *NativeArray = v
		if na.Length() == 0 {
			panic("assert")
		}
		name = na.Get(0).(*NativeString).Str
		if na.Length() > 1 {
			param = na.Get(1)
		}
	default:
		panic("unexpected type for NewRef")
	}
	ref := Ref{GNode: NewGNode(), ref: name, param: param}
	if name[0:1] == "_" {
		ref.GNode.Capture = false
	}
	return &ref
}

func (ref *Ref) GetGNode() *GNode        { return ref.GNode }
func (ref *Ref) HandlesChildLabel() bool { return false }
func (ref *Ref) Prepare()                {}
func (ref *Ref) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(ctx *ParseContext, _ Astnode) Astnode {
		var x Astnode = ref.GNode.Grammar.(*Grammar).Rules[ref.ref]
		if x == nil {
			panic("Unknown reference " + ref.ref)
		}
		ctx.StackPeek(0).Param = ref.param
		return x.Parse(ctx)
	}, ref)(ctx)
}

func (ref *Ref) Captures() []Astnode { return MeIfCaptureOrEmpty(ref) }
func (ref *Ref) Labels() []string {
	return ref._labels.GetCacheOrSet(func() []string {
		if ref.GNode.Label == "@" {
			return ref.GNode.Grammar.(*Grammar).Rules[ref.ref].Labels()
		} else if ref.GNode.Label != "" {
			return []string{ref.GNode.Label}
		} else {
			return []string{}
		}
	})
}

func (ref *Ref) ContentString() string {
	return LabelOrName(ref) + Yellow(ref.ref)
}
func (ref *Ref) ForEachChild(f func(Astnode) Astnode) Astnode {
	// no children defined in coffee
	return ref
}
