package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"reflect"
	"strings"
)

type Pattern struct {
	*GNode
	Value Ast       // in coffee, declared as `value: {type:GNode}`
	Join  Ast       // in coffee, declared as `join:  {type:GNode}`
	Min   NativeInt // -1 if unspec.
	Max   NativeInt // -1 if unspec.
}

// `it` must be a NativeMap with keys like 'value', 'join', 'min', 'max'
func NewPattern(it Ast) *Pattern {
	patt := &Pattern{NewGNode(), nil, nil, -1, -1}
	patt.Node = patt
	// {value Astnode, join Astnode, min int, max int}
	if nativemap, ok := it.(NativeMap); !ok {
		panic("Pattern expecting a map with value, join")
	} else {
		patt.Value = nativemap.Get("value")
		if patt.Value == nil {
			panic("Pattern must have a value")
		} else if _, is := patt.Value.(NativeString); is {
			// according to grammars at any rate it should be a Str,
			// which has a GNode
			panic("Pattern.Value can not be NativeString")
		}
		if patt.Value.GetGNode() != nil {
			patt.GetGNode().Capture = patt.Value.GetGNode().Capture
		} else { // Native* types don't have a GNode.
			// TODO coffee has @capture = @value.capture
			//  and @GNode has capture: yes
			// What's the right way?
			// This is probably just a theoritical case, so let's panic for now
			panic("patt.Value is a Native* type")
			patt.GetGNode().Capture = true
		}
		patt.Join = nativemap.Get("join") // can be nil
		patt.Min = NewNativeInt(-1)
		patt.Max = NewNativeInt(-1)
		if min, exists := nativemap.GetExists("min"); exists {
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
		}
		if max, exists := nativemap.GetExists("max"); exists {
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
	return patt
}
func (patt *Pattern) GetGNode() *GNode { return patt.GNode }
func (patt *Pattern) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		pos := ctx.Code.Pos
		resValue := patt.Value.Parse(ctx)
		if resValue == nil {
			ctx.Code.Pos = pos
			if patt.Min > 0 {
				return nil
			}
			return NewNativeArray([]Ast{})
		}
		var matches []Ast = []Ast{resValue}
		for true {
			pos2 := ctx.Code.Pos
			if NotNilAndNotNativeUndefined(patt.Join) {
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
			// fmt.Printf("Pattern matches = append(matches, resValue='%s')\n", resValue.ContentString())
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
func (patt *Pattern) Labels() []string        { panic("z") }
func (patt *Pattern) Captures() []Ast         { panic("z") }
func (patt *Pattern) Prepare()                {}
func (patt *Pattern) ContentString() string {
	var b strings.Builder
	b.WriteString(Prefix(patt.Value) + patt.Value.ContentString())
	b.WriteString(Cyan("*"))
	if NotNilAndNotNativeUndefined(patt.Join) {
		b.WriteString(Prefix(patt.Join) + patt.Join.ContentString())
	}
	if patt.Min < 0 && patt.Max < 0 {
		// Cyan("") so output is identical to coffee
		return b.String() + Cyan("")
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
func (patt *Pattern) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   value:      {type:GNode}
	//   join:       {type:GNode}
	patt.GetGNode().Rules = ForEachChild_InRules(patt, f)
	if patt.Value != nil {
		patt.Value = f(patt.Value)
	}
	if patt.Join != nil {
		patt.Join = f(patt.Join)
	}
	return patt
}
