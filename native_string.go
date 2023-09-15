package joeson

import (
	"strconv"
)

// NativeString is an alias for`string` but implements Ast.
type NativeString string

func NewNativeString(s string) NativeString { return NativeString(s) }
func (ns NativeString) String() string      { return string(ns) }
func (ns NativeString) assertNode()         {}

func (ns NativeString) SetLine(n int)                                   {}
func (ns NativeString) GetLine() int                                    { return 1 }
func (ns NativeString) SetOrigin(o Origin)                              {}
func (ns NativeString) GetOrigin() Origin                               { return Origin{} }
func (ns NativeString) HasAttribute(key interface{}) bool               { return false }
func (ns NativeString) GetAttribute(key interface{}) interface{}        { return nil }
func (ns NativeString) SetAttribute(key interface{}, value interface{}) { panic("N/A") }

func NewNativeStringFrom(x Ast) NativeString {
	return NativeStringFrom(x).(NativeString)
}

// creates a NativeString from Ast, only when that makes sense.
// It panics if necessary.
func NativeStringFrom(x Ast) Ast {
	switch v := x.(type) {
	case NativeString:
		return v
	case NativeInt:
		return NewNativeString(strconv.Itoa(v.Int()))
	case *NativeArray:
		return v.Concat()
	default:
		panic("Unable to make NativeString from " + x.String())
	}
}
