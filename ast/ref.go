package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"strconv"
)

type Ref struct {
	*GNode
	ref   string // ref because joeson.coffee used @ref, because @name was reserved
	param Ast    // thought it was `any`, see frame.go (`param` field) and joeson.coffee:67. But Astnode must be good
}

func NewRef(it Ast) *Ref {
	var name string
	var param Ast = nil
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
	ref := &Ref{GNode: NewGNode(), ref: name, param: param}
	ref.GNode.Node = ref
	if name[0:1] == "_" {
		ref.GNode.Capture = false
	}
	ref.GNode.Labels_ = helpers.NewLazy0[[]string](func() []string {
		if ref.GNode.Label == "@" {
			referenced := ref.GNode.Grammar.GetGNode().Rules[ref.ref]
			if referenced == nil {
				panic("ref " + ref.ref + " was not found in grammar.Rules")
			} else if referenced.GetGNode() == nil {
				panic("assert")
			} else {
				return referenced.GetGNode().Labels_.Get()
			}
		} else if ref.GNode.Label != "" {
			return []string{ref.GNode.Label}
		} else {
			return []string{}
		}
	})
	return ref
}

func (ref *Ref) GetGNode() *GNode        { return ref.GNode }
func (ref *Ref) HandlesChildLabel() bool { return false }
func (ref *Ref) Prepare()                {}
func (ref *Ref) Parse(ctx *ParseContext) Ast {
	return Wrap(func(ctx *ParseContext, _ Ast) Ast {
		node := ref.GNode.Grammar.GetGNode().Rules[ref.ref]
		if node == nil {
			panic("Unknown reference " + ref.ref + ". Grammar has " + strconv.Itoa(len(ref.GNode.Grammar.GetGNode().Rules)) + " rules. ")
		}
		ctx.StackPeek(0).Param = ref.param
		return node.Parse(ctx)
	}, ref)(ctx)
}

func (ref *Ref) ContentString() string { return Red(ref.ref) }
func (ref *Ref) ForEachChild(f func(Ast) Ast) Ast {
	// no children defined for Ref, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	ref.GetGNode().Rules = ForEachChild_InRules(ref, f)
	return ref
}
