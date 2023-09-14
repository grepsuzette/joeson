package joeson

import (
	"fmt"

	"github.com/grepsuzette/joeson/helpers"
)

type ref struct {
	*Attr
	*rule
	ref   string
	param Parser
}

// `it` can be a NativeString ("WORD")
//
//	or NativeArray with 1 element (["WORD"])
//	or NativeArray with 2 elements (["WORD", "EXPR"])
//	   That last case is built with
//
// o(s(r("WORD"), st("("), r("EXPR"), st(")")), func(it Ast) Ast {
func newRef(it Ast) *ref {
	var name string
	var param Parser = nil
	switch v := it.(type) {
	case NativeString:
		name = string(v)
	case *NativeArray:
		var na *NativeArray = v
		if na.Length() == 0 {
			panic("assert")
		}
		name = string(na.Get(0).(NativeString))
		if na.Length() > 1 {
			// fmt.Printf("ref param %s %T\n", na.Get(1).String(), na.Get(1))
			param = na.Get(1).(Parser)
		}
	default:
		panic(fmt.Sprintf("unexpected type for NewRef: %T %v\n", it, it))
	}
	ref := &ref{Attr: newAttr(), rule: newRule(), ref: name, param: param}
	ref.rule.node = ref
	if name[0:1] == "_" {
		ref.getRule().capture = false
	}
	ref.rule.labels_ = helpers.LazyFromFunc(func() []string {
		if ref.getRule().label == "@" {
			referenced := ref.grammar.getRuleRef(ref.ref)
			if referenced == nil {
				panic("ref " + ref.ref + " was not found in grammar.Rules")
			} else {
				return referenced.getRule().labels_.Get()
			}
		} else if ref.getRule().label != "" {
			return []string{ref.getRule().label}
		} else {
			return []string{}
		}
	})
	return ref
}

func (x *ref) getRule() *rule          { return x.rule }
func (x *ref) handlesChildLabel() bool { return false }
func (x *ref) prepare()                {}
func (x *ref) parse(ctx *ParseContext) Ast {
	return wrap(func(ctx *ParseContext, _ Parser) Ast {
		node := x.grammar.getRuleRef(x.ref)
		if node == nil {
			panic("Grammar has a reference to a type '" + x.ref + "' which is NOT defined")
		} else {
			ctx.stackPeek(0).param = x.param
			return node.parse(ctx)
		}
	}, x)(ctx)
}

func (x *ref) String() string { return Red(x.ref) }
func (x *ref) forEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	x.rules = ForEachChildInRules(x, f)
	return x
}
