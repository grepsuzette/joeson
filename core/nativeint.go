package core

import (
	"fmt"
	"strconv"
)

// NativeInt and NativeString denote our need to express a terminal element
// that at the same time satisfies the `Astnode` interface. Joeson.coffee used
// Number or string directly, but javascript is different language
type NativeInt int

func NewNativeInt(n int) NativeInt { return NativeInt(n) }
func NewNativeIntFromString(s string) NativeInt {
	fmt.Println(s)
	if n, e := strconv.Atoi(s); e == nil {
		return NativeInt(n)
	} else {
		panic("can not convert string " + s + " to NativeInt")
	}
}

// func NewNativeIntFromNativeString(ns NativeString) NativeInt {
// 	if n, e := strconv.Atoi(ns.Str); e == nil {
// 		return NewNativeInt(n)
// 	} else {
// 		panic(e)
// 	}
// }
func NewNativeIntFrom(x Astnode) NativeInt {
	switch v := x.(type) {
	case NativeString:
		if n, e := strconv.Atoi(v.Str); e == nil {
			return NewNativeInt(n)
		} else {
			panic(e)
		}
	case *NativeArray:
		s := ""
		for i, _ := range v.Array {
			s += v.Array[i].(NativeString).Str
		}
		return NewNativeIntFromString(s)
	default:
		panic("Unable to make NativeInt from " + x.ContentString())
	}
}

func (n NativeInt) String() string                  { return strconv.Itoa(int(n)) }
func (n NativeInt) Int() int                        { return int(n) }
func (n NativeInt) ContentString() string           { return strconv.Itoa(int(n)) }
func (n NativeInt) HandlesChildLabel() bool         { return false }
func (n NativeInt) GetGNode() *GNode                { return nil }
func (n NativeInt) Labels() []string                { return []string{} }
func (n NativeInt) Captures() []Astnode             { return []Astnode{} }
func (n NativeInt) Prepare()                        {}
func (n NativeInt) Parse(ctx *ParseContext) Astnode { panic("uncallable") }

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (n NativeInt) ForEachChild(f func(Astnode) Astnode) Astnode { return n }
