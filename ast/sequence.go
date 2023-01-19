package ast

import (
	"fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"strings"
)

type sequenceRepr int

const (
	Single sequenceRepr = 0
	Array               = 1
	Object              = 2
)

type Sequence struct {
	*GNode
	sequence []Ast
	type_    *helpers.Lazy[sequenceRepr] // internal cache for internalType()
}

func NewSequence(it Ast) *Sequence {
	if a, ok := it.(*NativeArray); !ok {
		panic("Sequence expected a NativeArray")
	} else {
		if a == nil {
			panic("expecting non nil array")
		}
		gn := NewGNode()
		seq := &Sequence{GNode: gn, sequence: a.Array}
		gn.Node = seq
		gn.Labels_ = helpers.NewLazyFromFunc[[]string](func() []string { return seq.calculateLabels() })
		gn.Captures_ = helpers.NewLazyFromFunc[[]Ast](func() []Ast { return seq.calculateCaptures() })
		seq.type_ = helpers.NewLazyFromFunc[sequenceRepr](func() sequenceRepr { return seq.calculateType() })
		return seq
	}
}

func (seq *Sequence) GetGNode() *GNode        { return seq.GNode }
func (seq *Sequence) HandlesChildLabel() bool { return true }
func (seq *Sequence) Prepare()                {}

func (seq *Sequence) calculateLabels() []string {
	a := []string{}
	for _, child := range seq.sequence {
		for _, label := range child.GetGNode().Labels_.Get() {
			a = append(a, label)
		}
	}
	return a
}
func (seq *Sequence) calculateCaptures() []Ast {
	a := []Ast{}
	for _, child := range seq.sequence {
		for _, captured := range child.GetGNode().Captures_.Get() {
			a = append(a, captured)
		}
	}
	return a
}

// as soon as there is >=1 label, it is Object
// otherwise, if at least 1 capture, it is Array
// otherwise a Single
func (seq *Sequence) calculateType() sequenceRepr {
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

func (seq *Sequence) ContentString() string {
	var b strings.Builder
	as := helpers.AMap(seq.sequence, func(x Ast) string { return String(x) })
	b.WriteString(strings.Join(as, " "))
	return Blue("(") + b.String() + Blue(")")
}

func (seq *Sequence) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
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
		panic("assert")
	}, seq)(ctx)
}

func (seq *Sequence) parseAsSingle(ctx *ParseContext) Ast {
	var result Ast = NewNativeUndefined()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		if child.GetGNode().Capture {
			result = res
		}
	}
	return result
}

func (seq *Sequence) parseAsArray(ctx *ParseContext) Ast {
	results := make([]Ast, 0)
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		if child.GetGNode().Capture {
			results = append(results, res)
		}
	}
	return NewNativeArray(results)
}

func (seq *Sequence) parseAsObject(ctx *ParseContext) Ast {
	var results Ast
	results = NewNativeUndefined()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			// fmt.Printf(Red("sequence %x %d parseAsObject childlabel=%s res==nil\n"), rnd, k, childLabel)
			return nil
		}
		if child.GetGNode().Label == "&" {
			if NotNilAndNotNativeUndefined(results) {
				results = Merge(res, results)
			} else {
				results = res
			}
		} else if child.GetGNode().Label == "@" {
			if NotNilAndNotNativeUndefined(results) {
				results = Merge(results, res)
			} else {
				results = res
			}
		} else if child.GetGNode().Label != "" {
			if NotNilAndNotNativeUndefined(results) {
				if h, isMap := results.(NativeMap); isMap {
					h.Set(child.GetGNode().Label, res)
				} else {
					panic("assert")
				}
			} else {
				results = NewEmptyNativeMap()
				results.(NativeMap).Set(child.GetGNode().Label, res)
			}
		}
	}
	return results
}

func (seq *Sequence) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   sequence:   {type:[type:GNode]}
	seq.GetGNode().Rules = ForEachChild_InRules(seq, f)
	if seq.sequence != nil {
		seq.sequence = ForEachChild_Array(seq.sequence, f)
	}
	return seq
}
