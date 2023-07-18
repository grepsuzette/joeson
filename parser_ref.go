package joeson

import (
	"fmt"
	"strconv"

	"github.com/grepsuzette/joeson/helpers"
)

type ref struct {
	Attributes
	*gnodeimpl
	ref   string // ref because joeson.coffee used @ref, because @name was reserved
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
		name = v.Str
	case *NativeArray:
		var na *NativeArray = v
		if na.Length() == 0 {
			panic("assert")
		}
		name = na.Get(0).(NativeString).Str
		if na.Length() > 1 {
			// fmt.Printf("ref param %s %T\n", na.Get(1).String(), na.Get(1))
			param = na.Get(1).(Parser)
		}
	default:
		panic(fmt.Sprintf("unexpected type for NewRef: %T %v\n", it, it))
	}
	ref := &ref{Attributes: Attributes{}, gnodeimpl: NewGNode(), ref: name, param: param}
	ref.gnodeimpl.node = ref
	if name[0:1] == "_" {
		ref.SetCapture(false)
	}
	ref.gnodeimpl.labels_ = helpers.NewLazyFromFunc(func() []string {
		if ref.Label() == "@" {
			referenced := ref.grammar.getRule(ref.ref)
			if referenced == nil {
				panic("ref " + ref.ref + " was not found in grammar.Rules")
			} else {
				return referenced.gnode().labels_.Get()
			}
		} else if ref.Label() != "" {
			return []string{ref.Label()}
		} else {
			return []string{}
		}
	})
	return ref
}

func (x *ref) gnode() *gnodeimpl       { return x.gnodeimpl }
func (x *ref) handlesChildLabel() bool { return false }
func (x *ref) prepare()                {}
func (x *ref) Parse(ctx *ParseContext) Ast {
	return wrap(func(ctx *ParseContext, _ Parser) Ast {
		node := x.grammar.getRule(x.ref)
		if node == nil {
			panic("Unknown reference " + x.ref + ". Grammar has " + strconv.Itoa(len(x.grammar.rules)) + " rules. ")
		}
		ctx.stackPeek(0).param = x.param
		return node.Parse(ctx)
	}, x)(ctx)
}

func (x *ref) String() string { return red(x.ref) }
func (x *ref) ForEachChild(f func(Parser) Parser) Parser {
	// no children defined for Ref, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	x.rules = ForEachChildInRules(x, f)
	return x
}
