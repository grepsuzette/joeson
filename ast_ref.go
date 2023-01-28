package joeson

import (
	"fmt"
	"grepsuzette/joeson/helpers"
	"strconv"
)

type ref struct {
	*GNode
	ref   string // ref because joeson.coffee used @ref, because @name was reserved
	param Ast
}

// `it` can be a NativeString ("WORD")
//      or NativeArray with 1 element (["WORD"])
//      or NativeArray with 2 elements (["WORD", "EXPR"])
//         That last case is built with
// o(s(r("WORD"), st("("), r("EXPR"), st(")")), func(it Ast) Ast {
func newRef(it Ast) *ref {
	var name string
	var param Ast = nil
	switch v := it.(type) {
	case NativeString:
		name = v.Str
	case *NativeArray:
		var na *NativeArray = v
		if na.Length() == 0 {
			panic("assert")
		}
		name = na.Get(0).(NativeString).Str
		if na.Length() > 1 {
			// fmt.Printf("ref param %s %T\n", na.Get(1).ContentString(), na.Get(1))
			param = na.Get(1)
		}
	default:
		panic(fmt.Sprintf("unexpected type for NewRef: %T %v\n", it, it))
	}
	ref := &ref{GNode: NewGNode(), ref: name, param: param}
	ref.GNode.Node = ref
	if name[0:1] == "_" {
		ref.GNode.Capture = false
	}
	ref.GNode.Labels_ = helpers.NewLazyFromFunc(func() []string {
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

func (x *ref) GetGNode() *GNode        { return x.GNode }
func (x *ref) HandlesChildLabel() bool { return false }
func (x *ref) Prepare()                {}
func (x *ref) Parse(ctx *ParseContext) Ast {
	return Wrap(func(ctx *ParseContext, _ Ast) Ast {
		node := x.GNode.Grammar.GetGNode().Rules[x.ref]
		if node == nil {
			panic("Unknown reference " + x.ref + ". Grammar has " + strconv.Itoa(len(x.GNode.Grammar.GetGNode().Rules)) + " rules. ")
		}
		ctx.stackPeek(0).Param = x.param
		return node.Parse(ctx)
	}, x)(ctx)
}

func (x *ref) ContentString() string { return red(x.ref) }
func (x *ref) ForEachChild(f func(Ast) Ast) Ast {
	// no children defined for Ref, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	x.GetGNode().Rules = ForEachChild_InRules(x, f)
	return x
}
