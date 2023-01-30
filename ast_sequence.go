package joeson

import (
	"fmt"
	"grepsuzette/joeson/helpers"
	"strings"
)

type sequenceRepr int

const (
	Single sequenceRepr = 0
	Array  sequenceRepr = 1
	Object sequenceRepr = 2
)

type sequence struct {
	*GNodeImpl
	sequence []Parser
	type_    *helpers.Lazy[sequenceRepr] // internal cache for internalType()
}

func newSequence(it Ast) *sequence {
	if a, ok := it.(*NativeArray); !ok {
		panic("Sequence expected a NativeArray")
	} else {
		if a == nil {
			panic("expecting non nil array")
		}
		gn := NewGNode()
		seq := &sequence{GNodeImpl: gn, sequence: helpers.AMap(a.Array, func(a Ast) Parser { return a.(Parser) })}
		gn.node = seq
		gn.Labels_ = helpers.NewLazyFromFunc(func() []string { return seq.calculateLabels() })
		gn.Captures_ = helpers.NewLazyFromFunc(func() []Ast { return seq.calculateCaptures() })
		seq.type_ = helpers.NewLazyFromFunc(func() sequenceRepr { return seq.calculateType() })
		return seq
	}
}

func (seq *sequence) GetGNode() *GNodeImpl    { return seq.GNodeImpl }
func (seq *sequence) HandlesChildLabel() bool { return true }
func (seq *sequence) Prepare()                {}

func (seq *sequence) calculateLabels() []string {
	a := []string{}
	for _, child := range seq.sequence {
		a = append(a, child.GetGNode().Labels_.Get()...)
	}
	return a
}
func (seq *sequence) calculateCaptures() []Ast {
	a := []Ast{}
	for _, child := range seq.sequence {
		a = append(a, child.GetGNode().Captures_.Get()...)
	}
	return a
}

// as soon as there is >=1 label, it is Object
// otherwise, if at least 1 capture, it is Array
// otherwise a Single
func (seq *sequence) calculateType() sequenceRepr {
	if len(seq.GetGNode().Labels_.Get()) == 0 {
		if len(seq.GetGNode().Captures_.Get()) > 1 {
			return Array
		} else {
			return Single
		}
	} else {
		return Object
	}
}

func (seq *sequence) ContentString() string {
	var b strings.Builder
	as := helpers.AMap(seq.sequence, func(x Parser) string { return String(x) })
	b.WriteString(strings.Join(as, " "))
	return blue("(") + b.String() + blue(")")
}

func (seq *sequence) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Parser) Ast {
		switch seq.type_.Get() {
		case Array:
			return seq.parseAsArray(ctx)
		case Single:
			return seq.parseAsSingle(ctx)
		case Object:
			if len(seq.sequence) == 0 {
				return NewNativeUndefined()
			} else {
				return seq.parseAsObject(ctx)
			}
		default:
			panic(fmt.Sprintf("Unexpected type %x", seq.type_.Get()))
		}
	}, seq)(ctx)
}

func (seq *sequence) parseAsSingle(ctx *ParseContext) Ast {
	var result Ast = NewNativeUndefined()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		if child.Capture() {
			result = res
		}
	}
	return result
}

func (seq *sequence) parseAsArray(ctx *ParseContext) Ast {
	results := make([]Ast, 0)
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		if child.Capture() {
			results = append(results, res)
		}
	}
	return NewNativeArray(results)
}

func (seq *sequence) parseAsObject(ctx *ParseContext) Ast {
	var results Ast
	results = NewNativeUndefined()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			// fmt.Printf(Red("sequence %x %d parseAsObject childlabel=%s res==nil\n"), rnd, k, childLabel)
			return nil
		}
		if child.Label() == "&" {
			if notNilAndNotNativeUndefined(results) {
				results = merge(res, results)
			} else {
				results = res
			}
		} else if child.Label() == "@" {
			if notNilAndNotNativeUndefined(results) {
				results = merge(results, res)
			} else {
				results = res
			}
		} else if child.Label() != "" {
			if notNilAndNotNativeUndefined(results) {
				if h, isMap := results.(NativeMap); isMap {
					h.Set(child.Label(), res)
				} else {
					panic("assert")
				}
			} else {
				results = NewEmptyNativeMap()
				results.(NativeMap).Set(child.Label(), res)
			}
		}
	}
	return results
}

func (seq *sequence) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   sequence:   {type:[type:GNode]}
	seq.GetGNode().Rules = ForEachChild_InRules(seq, f)
	if seq.sequence != nil {
		seq.sequence = ForEachChild_Array(seq.sequence, f)
	}
	return seq
}
