package joeson

import (
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

type NativeMap struct {
	*Origin
	Map map[string]Ast
}

func NewEmptyNativeMap() NativeMap            { return NewNativeMap(map[string]Ast{}) }
func NewNativeMap(h map[string]Ast) NativeMap { return NativeMap{&Origin{}, h} }

func (nm NativeMap) assertNode() {}
func (nm NativeMap) String() string {
	var b strings.Builder
	if nm.Origin.RuleName != "" {
		b.WriteString(nm.Origin.RuleName + "=")
	}
	b.WriteString("NativeMap{")
	first := true
	for _, k := range helpers.SortStringKeys(nm.Map) {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(k + ":" + nm.GetOrPanic(k).String())
		first = false
	}
	b.WriteString("}")
	return b.String()
}

func (nm NativeMap) Keys() []string {
	a := []string{}
	for k := range nm.Map {
		a = append(a, k)
	}
	return a
}

func (nm NativeMap) IsEmpty() bool {
	for range nm.Map {
		return true
	}
	return false
}

func (nm NativeMap) GetExists(k string) (Ast, bool) {
	v, exists := nm.Map[k]
	return v, exists
}

// specialized getter when value is known to be a NativeString (or NativeArray
// of NativeString).
// It panics when it is NOT a NativeString. The returned bool is
// false whenever the given string is not a key.
func (nm NativeMap) GetStringExists(k string) (string, bool) {
	if vv, exists := nm.Map[k]; exists {
		switch v := vv.(type) {
		case NativeString:
			return v.Str, true
		case *NativeArray:
			return stringFromNativeArray(v), true
		default:
			panic("unexpected type")
		}
	} else {
		return "", false
	}
}

// specialized getter when value is known to be a NativeInt.
// It panics when it is NOT a NativeInt. The returned bool is
// false whenever the given string is not a key.
func (nm NativeMap) GetIntExists(k string) (int, bool) {
	if v, exists := nm.Map[k]; exists {
		return v.(NativeInt).Int(), true
	} else {
		return 0, false
	}
}

// true when key is not defined or when its value is NativeUndefined.
// Note: successfully parsed Ast can't possibly return nil.
func (nm NativeMap) IsUndefined(k string) bool {
	if v, exists := nm.Map[k]; exists {
		if _, ok := v.(NativeUndefined); ok {
			return true
		} else {
			return false
		}
	}
	return true
}

func (nm NativeMap) GetOrPanic(k string) Ast {
	if r, ok := nm.Map[k]; ok {
		return r
	} else {
		panic("assert")
	}
}

// returns whichever of the first keys exist, or panic.
func (nm NativeMap) GetWhicheverOrPanic(a []string) Ast {
	for _, k := range a {
		if v, yes := nm.GetExists(k); yes {
			return v
		}
	}
	panic("No key found")
}

func (nm NativeMap) Get(k string) Ast {
	return nm.Map[k]
}

// if key doesn't exist return nil, otherwise forces a Parser or panic.
func (nm NativeMap) GetParser(k string) Parser {
	switch x := nm.Map[k].(type) {
	case nil:
		return nil
	case Parser:
		return x
	case Ast:
		panic("assert Parser expected, not Ast: " + x.String())
	default:
		panic("assert")
	}
}

func (nm NativeMap) Set(k string, v Ast) {
	nm.Map[k] = v
}

func (nm NativeMap) Exists(k string) bool {
	_, ok := nm.Map[k]
	return ok
}
