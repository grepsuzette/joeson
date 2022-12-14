package ast

import (
	// "fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"reflect"
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
		if a == nil {
			panic("expecting non nil array")
		}
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
	b.WriteString(LabelOrName(seq))
	for _, x := range seq.sequence {
		b.WriteString(x.ContentString() + " ")
	}
	return Blue("(") + b.String() + Blue(")")
}

func (seq *Sequence) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
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

func (seq *Sequence) parseAsArray(ctx *ParseContext) Astnode {
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
	for k, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		// TODO dubious comment below, edit soon
		// if there is a label, child res is normally a NativeMap
		// otherwise, child res is a Ref
		if child.GetGNode().Label == "&" {
			switch v := res.(type) {
			case NativeMap:
				// TODO seems this case never happens after all,
				// it seems to be a Ref instead
				panic("AGAGAGA")
				resMap := v
				for _, k := range results.Keys() {
					resMap.Set(k, results.Get(k))
				}
				results = resMap
			case *Ref:
				if k == len(seq.sequence)-1 {
					return v
				} else {
					panic("unhandled case, where Ref in & is not the final element in a sequence, study how to merge")
				}
			case *Choice:
				return v
			case *Not:
				return v
			case *Regex:
				return v
			case Str:
				if k == len(seq.sequence)-1 {
					return v
				} else {
					panic("unhandled case, where Str in & is not the final element in a sequence, study how to merge")
				}
			case *Pattern:
				// fmt.Println("Pattern:")
				// fmt.Println("Value: " + v.Value.ContentString())
				// fmt.Println("Join: " + v.Join.ContentString())
				// fmt.Println("Min: " + v.Min.ContentString())
				// fmt.Println("Max: " + v.Max.ContentString())
				return v
			case *Existential:
				return v
			default:
				panic("unhandled type in parseAsObject: " + reflect.TypeOf(v).String() + "\n" + ctx.Code.Print())
			}
		} else if child.GetGNode().Label == "@" {
			if _, isUndefined := res.(NativeUndefined); !isUndefined {
				h := res.(NativeMap)
				for _, k := range h.Keys() {
					results.Set(k, h.Get(k))
				}
			}
		} else if child.GetGNode().Label != "" {
			results.Set(child.GetGNode().Label, res)
		}
	}
	return results
}
func (seq *Sequence) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   sequence:   {type:[type:GNode]}
	seq.GetGNode().Rules = ForEachChild_MapString(seq.GetGNode().Rules, f)
	if seq.sequence != nil {
		seq.sequence = ForEachChild_Array(seq.sequence, f)
	}
	return seq
}
