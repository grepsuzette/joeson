package ast

import (
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
	sequence      []Astnode
	_internalType helpers.Varcache[sequenceRepr] // internal cache for internalType()
	_labels       helpers.Varcache[[]string]     // internal cache for Labels()
	_captures     helpers.Varcache[[]Astnode]    // internal cache for Captures()
}

func NewSequence(it Astnode) *Sequence {
	if a, ok := it.(*NativeArray); ok {
		return &Sequence{GNode: NewGNode(), sequence: a.Array}
	} else {
		panic("Sequence expected a NativeArray")
	}
}

func (seq *Sequence) GetGNode() *GNode        { return seq.GNode }
func (seq *Sequence) HandlesChildLabel() bool { return true }
func (seq *Sequence) Prepare()                {}

func (seq *Sequence) Labels() []string {
	return seq._labels.GetCacheOrSet(func() []string {
		a := []string{}
		for _, child := range seq.sequence {
			for _, label := range child.Labels() {
				a = append(a, label)
			}
		}
		return a
	})
}

func (seq *Sequence) Captures() []Astnode {
	return seq._captures.GetCacheOrSet(func() []Astnode {
		a := []Astnode{}
		for _, child := range seq.sequence {
			for _, capture := range child.Captures() {
				a = append(a, capture)
			}
		}
		return a
	})
}

// as soon as there is >=1 label, it is Object
// otherwise, if at least 1 capture, it is Array
// otherwise a Single
func (seq *Sequence) internalType() sequenceRepr {
	return seq._internalType.GetCacheOrSet(func() sequenceRepr {
		if len(seq.Labels()) == 0 {
			if len(seq.Captures()) > 1 {
				return Array
			} else {
				return Single
			}
		} else {
			return Object
		}
	})
}

func (seq *Sequence) ContentString() string {
	var b strings.Builder
	b.WriteString(ShowLabelOrNameIfAny(seq))
	for _, x := range seq.sequence {
		b.WriteString(x.ContentString() + " ")
	}
	return Blue("(") + b.String() + Blue(")")
}

func (seq *Sequence) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
		switch seq.internalType() {
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
			panic("Unexpected type " + string(seq.internalType()))
		}
		panic("Error")
	}, seq)(ctx)
}

func (seq *Sequence) parseAsSingle(ctx *ParseContext) Astnode {
	var result Astnode = NewNativeUndefined()
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

func (seq *Sequence) parseAsArray(ctx *ParseContext) *NativeArray {
	results := make([]Astnode, 0)
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

func (seq *Sequence) parseAsObject(ctx *ParseContext) Astnode {
	results := NewEmptyNativeMap()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		// because it has labels, child res is basically guaranteed to be
		// a NativeMap both for "&" and "@"
		if child.GetGNode().Label == "&" {
			if results == nil {
				results = res.(NativeMap)
			} else {
				resMap := res.(NativeMap)
				for _, k := range results.Keys() {
					resMap.Set(k, results.Get(k))
				}
				results = resMap
			}
		} else if child.GetGNode().Label == "@" {
			if results == nil {
				results = res.(NativeMap)
			} else {
				h := res.(NativeMap)
				for _, k := range h.Keys() {
					results.Set(k, h.Get(k))
				}
			}
		} else if child.GetGNode().Label != "" {
			// TODO ^ i have a doubt that maybe label should be *string because of joeson.go:832...
			results.Set(child.GetGNode().Label, res)
		}
	}
	return results
}
func (seq *Sequence) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   sequence:   {type:[type:GNode]}
	// this one seems tricky,
	//  but think can recursively work with Native*.ForEachChild
	if seq.sequence != nil {
		seq.sequence = ForEachChild_Array(seq.sequence, f)
	}
	seq.GetGNode().Rules = ForEachChild_MapString(seq.GetGNode().Rules, f)
	return seq
}
