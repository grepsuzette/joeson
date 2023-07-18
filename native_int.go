package joeson

import (
	"strconv"
)

// NativeInt is an `int` that implements `Ast`.
type NativeInt struct {
	*Attributes
	int
}

func NewNativeInt(n int) NativeInt { return NativeInt{&Attributes{}, n} }
func NewNativeIntFromString(s string) NativeInt {
	if n, e := strconv.Atoi(s); e == nil {
		return NativeInt{&Attributes{}, n}
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
