package joeson

import (
	"strconv"
)

// NativeInt is an `int` that implements `Ast`.
type NativeInt int

func NewNativeInt(n int) NativeInt { return NativeInt(n) }
func NewNativeIntFromBool(b bool) NativeInt {
	n := 0
	if b {
		n = 1
	}
	return NativeInt(n)
}

func NewNativeIntFromString(s string) NativeInt {
	if n, e := strconv.Atoi(s); e == nil {
		return NativeInt(n)
	} else {
		panic("can not convert string " + s + " to NativeInt")
	}
}

func NewNativeIntFrom(x Ast) NativeInt {
	return NativeIntFrom(x).(NativeInt)
}

// Create a NativeInt from an Ast.
// Return an Ast.
// It panics if necessary.
func NativeIntFrom(x Ast) Ast {
	switch v := x.(type) {
	case NativeInt:
		return v
	case NativeString:
		if n, e := strconv.Atoi(string(v)); e == nil {
			return NewNativeInt(n)
		} else {
			panic(e)
		}
	case *NativeArray:
		// OPTIM
		return NewNativeIntFromString(v.Concat().String())
	default:
		panic("Unable to make NativeInt from " + x.String())
	}
}

func (n NativeInt) assertNode()    {}
func (n NativeInt) Int() int       { return int(n) }
func (n NativeInt) Bool() bool     { return int(n) != 0 }
func (n NativeInt) String() string { return strconv.Itoa(int(n)) }

func (n NativeInt) SetLine(m int)                                   {}
func (n NativeInt) GetLine() int                                    { return 1 }
func (n NativeInt) SetOrigin(o Origin)                              {}
func (n NativeInt) GetOrigin() Origin                               { return Origin{} }
func (n NativeInt) HasAttribute(key interface{}) bool               { return false }
func (n NativeInt) GetAttribute(key interface{}) interface{}        { return nil }
func (n NativeInt) SetAttribute(key interface{}, value interface{}) { panic("N/A") }
