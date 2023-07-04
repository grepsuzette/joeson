package joeson

import (
	"fmt"
	"strconv"
)

// NativeInt is an `int` but also an object that satisfies the `Ast` interface.
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
		for i := range v.Array {
			s += v.Array[i].(NativeString).Str
		}
		return NewNativeIntFromString(s)
	default:
		panic("Unable to make NativeInt from " + x.ContentString())
	}
}

func (n NativeInt) String() string        { return strconv.Itoa(int(n)) }
func (n NativeInt) Int() int              { return int(n) }
func (n NativeInt) ContentString() string { return strconv.Itoa(int(n)) }
