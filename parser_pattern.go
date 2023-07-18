package joeson

import (
	"reflect"
	"strings"
)

type pattern struct {
	*gnodeimpl
	Value Parser
	Join  Parser
	Min   NativeInt // -1 for unspec.
	Max   NativeInt // -1 for unspec.
}

// `it` must be a NativeMap with keys like 'value', 'join', 'min', 'max'
func newPattern(it Ast) *pattern {
	patt := &pattern{NewGNode(), nil, nil, -1, -1}
	patt.node = patt
	if nativemap, ok := it.(NativeMap); !ok {
		panic("Pattern expecting a map with value, join")
	} else {
		patt.Value = nativemap.Get("value").(Parser)
		if patt.Value == nil {
			panic("Pattern must have a value")
		} else {
			patt.SetCapture(patt.Value.Capture())
		}
		patt.Join = nativemap.GetParser("join") // can be nil
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
func (patt *pattern) gnode() *gnodeimpl { return patt.gnodeimpl }
func (patt *pattern) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
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
		for {
			pos2 := ctx.Code.Pos
			if notNilAndNotNativeUndefined(patt.Join) {
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

func (patt *pattern) handlesChildLabel() bool { return false }
func (patt *pattern) prepare()                {}
func (patt *pattern) ContentString() string {
	var b strings.Builder
	b.WriteString(String(patt.Value))
	b.WriteString(cyan("*"))
	if notNilAndNotNativeUndefined(patt.Join) {
		b.WriteString(String(patt.Join))
	}
	if patt.Min < 0 && patt.Max < 0 {
		// Cyan("") so output is identical to coffee
		return b.String() + cyan("")
	} else {
		sCyan := "{"
		if patt.Min > -1 {
			sCyan += patt.Min.String()
		}
		sCyan += ","
		if patt.Max > -1 {
			sCyan += patt.Max.String()
		}
		sCyan += "}"
		return b.String() + cyan(sCyan)
	}
}

func (patt *pattern) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   value:      {type:GNode}
	//   join:       {type:GNode}
	patt.rules = ForEachChildInRules(patt, f)
	if patt.Value != nil {
		patt.Value = f(patt.Value)
	}
	if patt.Join != nil {
		patt.Join = f(patt.Join)
	}
	return patt
}
