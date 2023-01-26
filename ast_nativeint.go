package joeson

import (
	"fmt"
	"strconv"
)

// NativeInt and NativeString denote terminal nodes
// and at the same time satisfy the `Ast` interface. Joeson.coffee used
// Number or string directly.

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

func NewNativeIntFrom(x Ast) NativeInt {
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

func (n NativeInt) String() string              { return strconv.Itoa(int(n)) }
func (n NativeInt) Int() int                    { return int(n) }
func (n NativeInt) ContentString() string       { return strconv.Itoa(int(n)) }
func (n NativeInt) HandlesChildLabel() bool     { return false }
func (n NativeInt) GetGNode() *GNode            { return nil }
func (n NativeInt) Prepare()                    {}
func (n NativeInt) Parse(ctx *ParseContext) Ast { panic("uncallable") }

func (n NativeInt) ForEachChild(f func(Ast) Ast) Ast { return n }
