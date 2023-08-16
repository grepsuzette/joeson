package joeson

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// what would it take to have nativemap
// handles order?
// - nm.Map would have to be unimported
//

type NativeMap struct {
	*Attr
	vals map[string]Ast
	keys []string // the order, for Concat() to work well with sequence parser
}

func NewEmptyNativeMap() *NativeMap { return NewNativeMap(map[string]Ast{}) }
func NewNativeMap(h map[string]Ast) *NativeMap {
	// since there may exist keys in h,
	// we decide an arbitrary but predictable order for them
	return &NativeMap{
		newAttr(),
		h,
		helpers.SortStringKeys(h),
	}
}

// NativeMap.Concat works provided all types within its tree are Native*.
// it will panic otherwise.
func (nm *NativeMap) Concat() string {
	var b strings.Builder
	for _, k := range nm.Keys() {
		if v, ok := nm.GetExists(k); ok {
			switch w := v.(type) {
			case *NativeMap:
				b.WriteString(w.Concat())
			case *NativeArray:
				b.WriteString(w.Concat())
			case NativeUndefined:
			case NativeString:
				b.WriteString(string(w))
			case NativeInt:
				b.WriteString(w.String())
			default:
				panic("NativeMap.Concat only works with Native* types")
			}
		}
	}
	return b.String()
}

func (nm *NativeMap) assertNode() {}
func (nm *NativeMap) String() string {
	var b strings.Builder
	b.WriteString("NativeMap{")
	first := true
	for _, k := range nm.Keys() {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(k + ":" + nm.GetOrPanic(k).String())
		first = false
	}
	b.WriteString("}")
	return b.String()
}

func (nm *NativeMap) Keys() []string {
	return nm.keys
}

func (nm *NativeMap) IsEmpty() bool {
	for range nm.vals {
		return true
	}
	return false
}

func (nm *NativeMap) GetExists(k string) (Ast, bool) {
	v, exists := nm.vals[k]
	return v, exists
}

// specialized getter when value is known to be a NativeString (or NativeArray
// of NativeString).
// It panics when it is NOT a NativeString. The returned bool is
// false whenever the given string is not a key.
func (nm *NativeMap) GetStringExists(k string) (string, bool) {
	if vv, exists := nm.vals[k]; exists {
		switch v := vv.(type) {
		case NativeString:
			return string(v), true
		case *NativeArray:
			return v.Concat(), true
		default:
			panic("unexpected type")
		}
	} else {
		return "", false
	}
}

// specialized getter when value is known to an int.
// It panics when the key exists but when it can not retrieve an int.
// The returned bool is false whenever the given string is not a key.
func (nm *NativeMap) GetIntExists(k string) (int, bool) {
	if ast, exists := nm.vals[k]; exists {
		switch v := ast.(type) {
		case NativeInt:
			return v.Int(), true
		case NativeString:
			if n, e := strconv.Atoi(string(v)); e == nil {
				return n, true
			} else {
				panic("Could not Atoi(" + string(v) + "): " + e.Error())
			}
		case *NativeArray:
			if n, e := strconv.Atoi(v.Concat()); e == nil {
				return n, true
			} else {
				panic("Could not Atoi(" + v.Concat() + "): " + e.Error())
			}
		default:
			panic("Could not get Int from " + ast.String())
		}
	} else {
		return 0, false
	}
}

// true when key is not defined or when its value is NativeUndefined.
// Note: successfully parsed Ast can't possibly return nil.
func (nm *NativeMap) IsUndefined(k string) bool {
	if v, exists := nm.vals[k]; exists {
		if _, ok := v.(NativeUndefined); ok {
			return true
		} else {
			return false
		}
	}
	return true
}

func (nm *NativeMap) GetOrPanic(k string) Ast {
	if r, ok := nm.vals[k]; ok {
		return r
	} else {
		panic(fmt.Sprintf("no key '%s'", k))
	}
}

// returns whichever of the first keys exist, or panic.
func (nm *NativeMap) GetWhicheverOrPanic(a []string) Ast {
	for _, k := range a {
		if v, yes := nm.GetExists(k); yes {
			return v
		}
	}
	panic("No key found")
}

func (nm *NativeMap) Get(k string) Ast {
	return nm.vals[k]
}

func (nm *NativeMap) Set(k string, v Ast) {
	nm.vals[k] = v
	// keep order of key if already in
	for _, k2 := range nm.keys {
		if k == k2 {
			return
		}
	}
	nm.keys = append(nm.keys, k)
}

func (nm *NativeMap) Exists(k string) bool {
	_, ok := nm.vals[k]
	return ok
}
