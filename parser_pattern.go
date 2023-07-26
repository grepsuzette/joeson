package joeson

import (
	"reflect"
	"strconv"
	"strings"
)

type pattern struct {
	Attr
	*gnodeimpl
	value Parser
	join  Parser
	min   int // -1 for unspec.
	max   int // -1 for unspec.
}

// `it` must be a NativeMap with keys like 'value', 'join', 'min', 'max'
func newPattern(it Ast) *pattern {
	patt := &pattern{newAttr(), newGNode(), nil, nil, -1, -1}
	patt.node = patt
	if nativemap, ok := it.(*NativeMap); !ok {
		panic("Pattern expecting a map with value, join")
	} else {
		patt.value = nativemap.Get("value").(Parser)
		if patt.value == nil {
			panic("Pattern must have a value")
		} else {
			patt.SetCapture(patt.value.Capture())
		}
		if join, exists := nativemap.GetExists("join"); exists {
			if join == nil {
				patt.join = nil
			} else {
				patt.join = join.(Parser)
			}
		} else {
			patt.join = nil
		}
		patt.min = -1
		patt.max = -1
		if min, exists := nativemap.GetExists("min"); exists {
			switch v := min.(type) {
			case NativeUndefined:
				patt.min = -1
			case NativeInt:
				patt.min = v.Int()
			case NativeString:
				patt.min = NewNativeIntFromString(v.Str).Int()
			default:
				panic("NewPattern unhandled type for min: " + reflect.TypeOf(min).String())
			}
		}
		if max, exists := nativemap.GetExists("max"); exists {
			switch v := max.(type) {
			case NativeUndefined:
				patt.max = -1
			case NativeInt:
				patt.max = v.Int()
			case NativeString:
				patt.max = NewNativeIntFromString(v.Str).Int()
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
		resValue := patt.value.Parse(ctx)
		if resValue == nil {
			ctx.Code.Pos = pos
			if patt.min > 0 {
				return nil
			}
			return NewNativeArray([]Ast{})
		}
		var matches []Ast = []Ast{resValue}
		for {
			pos2 := ctx.Code.Pos
			if !isUndefined(patt.join) {
				resJoin := patt.join.Parse(ctx)
				// return nil to revert pos
				if resJoin == nil {
					ctx.Code.Pos = pos2
					break
				}
			}
			resValue = patt.value.Parse(ctx)
			// return nil to revert pos
			if resValue == nil {
				ctx.Code.Pos = pos2
				break
			}
			matches = append(matches, resValue)
			if patt.max > -1 && len(matches) >= int(patt.max) {
				break
			}
		}
		if patt.min > -1 && int(patt.min) > len(matches) {
			ctx.Code.Pos = pos
			return nil
		}
		return NewNativeArray(matches)
	}, patt)(ctx)
}

func (patt *pattern) handlesChildLabel() bool { return false }
func (patt *pattern) prepare()                {}
func (patt *pattern) String() string {
	var b strings.Builder
	b.WriteString(String(patt.value))
	b.WriteString(Cyan("*"))
	if !isUndefined(patt.join) {
		b.WriteString(String(patt.join))
	}
	if patt.min < 0 && patt.max < 0 {
		// Cyan("") so output is identical to coffee
		return b.String() + Cyan("")
	} else {
		sCyan := "{"
		if patt.min > -1 {
			sCyan += strconv.Itoa(patt.min)
		}
		sCyan += ","
		if patt.max > -1 {
			sCyan += strconv.Itoa(patt.max)
		}
		sCyan += "}"
		return b.String() + Cyan(sCyan)
	}
}

func (patt *pattern) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   value:      {type:GNode}
	//   join:       {type:GNode}
	patt.rules = ForEachChildInRules(patt, f)
	if patt.value != nil {
		patt.value = f(patt.value)
	}
	if patt.join != nil {
		patt.join = f(patt.join)
	}
	return patt
}
