package joeson

import (
	"fmt"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

type sequenceRepr int

const (
	Single sequenceRepr = 0
	Array  sequenceRepr = 1
	Object sequenceRepr = 2
)

type sequence struct {
	*Attr
	*rule
	sequence []Parser
	lazyType *helpers.Lazy[sequenceRepr] // internal cache for internalType()
}

func newSequence(it Ast) *sequence {
	if a, ok := it.(*NativeArray); !ok {
		panic("Sequence expected a NativeArray")
	} else {
		if a == nil {
			panic("expecting non nil array")
		}
		gn := newRule()
		parsers := make([]Parser, 0)
		for _, v := range *a {
			parsers = append(parsers, v.(Parser))
		}
		seq := &sequence{
			Attr:     newAttr(),
			rule:     gn,
			sequence: parsers,
		}
		gn.node = seq
		gn.labels_ = helpers.LazyFromFunc(func() []string { return seq.calculateLabels() })
		gn.captures_ = helpers.LazyFromFunc(func() []Ast { return seq.calculateCaptures() })
		seq.lazyType = helpers.LazyFromFunc(func() sequenceRepr { return seq.calculateType() })
		return seq
	}
}

func (seq *sequence) gnode() *rule            { return seq.rule }
func (seq *sequence) HandlesChildLabel() bool { return true }
func (seq *sequence) prepare()                {}

func (seq *sequence) calculateLabels() []string {
	a := []string{}
	for _, child := range seq.sequence {
		a = append(a, child.gnode().labels_.Get()...)
	}
	return a
}

func (seq *sequence) calculateCaptures() []Ast {
	a := []Ast{}
	for _, child := range seq.sequence {
		a = append(a, child.gnode().captures_.Get()...)
	}
	return a
}

// as soon as there is >=1 label, it is Object
// otherwise, if at least 1 capture, it is Array
// otherwise a Single
func (seq *sequence) calculateType() sequenceRepr {
	if len(seq.labels_.Get()) == 0 {
		if len(seq.captures_.Get()) > 1 {
			return Array
		} else {
			return Single
		}
	} else {
		return Object
	}
}

func (seq *sequence) String() string {
	var b strings.Builder
	first := true
	for _, v := range seq.sequence {
		if !first {
			b.WriteString(" ")
		}
		b.WriteString(String(v))
		first = false
	}
	return Blue("(") + b.String() + Blue(")")
}

func (seq *sequence) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		switch seq.lazyType.Get() {
		case Array:
			return seq.parseAsArray(ctx)
		case Single:
			return seq.parseAsSingle(ctx)
		case Object:
			if len(seq.sequence) == 0 {
				// seems never called
				return NewNativeUndefined()
			} else {
				return seq.parseAsObject(ctx)
			}
		default:
			panic(fmt.Sprintf("Unexpected type %x", seq.lazyType.Get()))
		}
	}, seq)(ctx)
}

// OPTIMIZE
func (seq *sequence) parseAsSingle(ctx *ParseContext) Ast {
	var result Ast = nil // OPTIMIZE critical function (about 100k calls to parse the intention grammar), so avoid needless calls to NewNativeUndefined()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			return nil
		}
		if child.Capture() {
			result = res
		}
	}
	if result == nil {
		return NewNativeUndefined()
	} else {
		return result
	}
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

// OPTIMIZE
func (seq *sequence) parseAsObject(ctx *ParseContext) Ast {
	var results Ast
	results = nil // OPTIMIZE critical function (about 100k calls to parse the intention grammar), so avoid needless calls to NewNativeUndefined()
	for _, child := range seq.sequence {
		res := child.Parse(ctx)
		if res == nil {
			// fmt.Printf(Red("sequence %x %d parseAsObject childlabel=%s res==nil\n"), rnd, k, childLabel)
			return nil
		}
		label := child.GetRuleLabel()
		switch label {
		case "&":
			if IsUndefined(results) {
				results = res
			} else {
				results = merge(res, results)
			}
		case "@":
			if IsUndefined(results) {
				results = res
			} else {
				results = merge(results, res)
			}
		case "":
		default:
			if IsUndefined(results) {
				results = NewEmptyNativeMap()
				results.(*NativeMap).Set(label, res)
			} else {
				if h, isMap := results.(*NativeMap); isMap {
					h.Set(label, res)
				} else {
					panic("assert")
				}
			}
		}
	}
	if results == nil { // see comment above (critical function)
		return NewNativeUndefined()
	} else {
		return results
	}
}

func (seq *sequence) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   sequence:   {type:[type:GNode]}
	seq.rules = ForEachChildInRules(seq, f)
	if seq.sequence != nil {
		seq.sequence = ForEachChild_Array(seq.sequence, f)
	}
	return seq
}

// Port of lib/helpers.js:extend() in a less general way (Ast-specific)
// Extend a source object with the properties of a newcomer object (shallow copy).
// The modified `toExtend` object is returned.
func merge(toExtend Ast, newcomer Ast) Ast {
	// @extend = extend = (object, properties) ->
	//   for key, val of properties
	//     object[key] = val
	//   object
	if toExtend == nil || newcomer == nil {
		return toExtend
	}
	switch vnewcomer := newcomer.(type) {
	case NativeUndefined:
		return toExtend
	case *NativeMap:
		switch vToExtend := toExtend.(type) {
		case *NativeMap:
			for _, k := range vnewcomer.Keys() {
				vToExtend.Set(k, vnewcomer.GetOrPanic(k))
			}
		case Parser:
			for _, k := range vnewcomer.Keys() {
				switch k {
				case "label":
					value := vnewcomer.GetOrPanic(k)
					vToExtend.SetRuleLabel(string(value.(NativeString)))
				default:
					panic("unhandled property " + k + " in func (Ast) Merge(). toExtend=" + toExtend.String() + " \n withPropertiesOf=" + newcomer.String())
				}
			}
		default:
			panic("assert")
		}
		return toExtend
	case Parser:
		switch toExtend.(type) {
		case Parser:
			panic("Unhandled case in func (Ast) Merge()")
		default:
			panic("assert")
		}
	default:
		panic("assert")
	}
}
