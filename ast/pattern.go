package ast

import (
	"fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"reflect"
	"strings"
)

type Pattern struct {
	*GNode
	Value Astnode   // in coffee, declared as `value: {type:GNode}`
	Join  Astnode   // in coffee, declared as `join:  {type:GNode}`
	Min   NativeInt // -1 if unspec.
	Max   NativeInt // -1 if unspec.
}

// it is a NativeMap with keys 'value', 'join' and '@'
func NewPattern(it Astnode) *Pattern {
	patt := Pattern{NewGNode(), nil, nil, -1, -1}
	// (value Astnode, join Astnode, min int, max int)
	if nativemap, ok := it.(NativeMap); !ok {
		panic("Pattern expecting a map with value, join")
	} else {
		patt.Value = nativemap.Get("value")
		patt.GetGNode().Capture = patt.Value.GetGNode().Capture
		patt.Join = nativemap.Get("join") // can be nil
		patt.Min = NewNativeInt(-1)
		patt.Max = NewNativeInt(-1)
		if min, exist := nativemap.GetExist("min"); exist {
			// can also be NativeUndefined, when Existential returns it
			switch v := min.(type) {
			case NativeUndefined:
				patt.Min = NewNativeInt(-1)
			case NativeInt:
				patt.Min = v
			case NativeString:
				patt.Min = NewNativeIntFromString(v.Str)
			default:
				panic("NewPattern unhandled type for min: " + reflect.TypeOf(min).String())
			}
		} else if max, exist := nativemap.GetExist("max"); exist {
			// can also be NativeUndefined, when Existential returns it
			switch v := max.(type) {
			case NativeUndefined:
				patt.Max = NewNativeInt(-1)
			case NativeInt:
				patt.Max = v
			case NativeString:
				patt.Max = NewNativeIntFromString(v.Str)
			default:
				panic("NewPattern unhandled type for max: " + reflect.TypeOf(max).String())
			}
		}
	}
	return &patt
}
func (patt *Pattern) GetGNode() *GNode { return patt.GNode }
func (patt *Pattern) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
		pos := ctx.Code.Pos
		resValue := patt.Value.Parse(ctx)
		if resValue == nil {
			ctx.Code.Pos = pos
			if patt.Min > 0 {
				return nil
			}
			return NewNativeArray([]Astnode{})
		}
		var matches []Astnode = []Astnode{resValue}
		for true {
			pos2 := ctx.Code.Pos
			if patt.Join != nil {
				resJoin := patt.Join.Parse(ctx)
				// return nil to revert pos
				if resJoin == nil {
					ctx.Code.Pos = pos2
					break
				}
			}
			resValue = patt.Value.Parse(ctx)
			// return nil to revert pos
			if resValue == nil {
				ctx.Code.Pos = pos2
				break
			}
			fmt.Printf("Pattern matches = append(matches, resValue='%s')\n", resValue.ContentString())
			matches = append(matches, resValue)
			if patt.Max > -1 && len(matches) >= int(patt.Max) {
				break
			}
		}
		if patt.Min > -1 && int(patt.Min) > len(matches) {
			ctx.Code.Pos = pos
			return nil
		}
		return NewNativeArray(matches)
	}, patt)(ctx)
}

func (patt *Pattern) HandlesChildLabel() bool { return false }
func (patt *Pattern) Labels() []string        { return MyLabelIfDefinedOrEmpty(patt) }
func (patt *Pattern) Captures() []Astnode     { return MeIfCaptureOrEmpty(patt) }
func (patt *Pattern) Prepare()                {}
func (patt *Pattern) ContentString() string {
	var b strings.Builder
	b.WriteString(LabelOrName(patt))
	b.WriteString(patt.Value.ContentString() + Cyan("*"))
	if patt.Join != nil {
		b.WriteString(patt.Join.ContentString())
	}
	if patt.Min < 0 && patt.Max < 0 {
		return b.String()
	} else {
		cyan := "{"
		if patt.Min > -1 {
			cyan += patt.Min.String()
		}
		cyan += ","
		if patt.Max > -1 {
			cyan += patt.Max.String()
		}
		cyan += "}"
		return b.String() + Cyan(cyan)
	}
}
func (patt *Pattern) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   value:      {type:GNode}
	//   join:       {type:GNode}
	patt.GetGNode().Rules = ForEachChild_MapString(patt.GetGNode().Rules, f)
	if patt.Value != nil {
		patt.Value = f(patt.Value)
	}
	if patt.Join != nil {
		patt.Join = f(patt.Join)
	}
	return patt
}
