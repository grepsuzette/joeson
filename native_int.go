package joeson

import (
	"strconv"
)

// NativeInt is an `int` that implements `Ast`.
type NativeInt struct {
	*Attr
	int
}

func NewNativeInt(n int) NativeInt { return NativeInt{newAttr(), n} }
func NewNativeIntFromString(s string) NativeInt {
	if n, e := strconv.Atoi(s); e == nil {
		return NativeInt{newAttr(), n}
	} else {
		panic("can not convert string " + s + " to NativeInt")
	}
}

// creates a NativeInt from, possibly: 1. a NativeString.
// 2. a NativeArray of NativeString.
// It panics if necessary.
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
		for i := range v.Array {
			s += v.Array[i].(NativeString).Str
		}
		return NewNativeIntFromString(s)
	default:
		panic("Unable to make NativeInt from " + x.String())
	}
}

func (n NativeInt) assertNode()    {}
func (n NativeInt) Int() int       { return n.int }
func (n NativeInt) String() string { return strconv.Itoa(n.int) }
