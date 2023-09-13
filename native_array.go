package joeson

/*
"Native*" types
===============

These types wrap `int, string, map, array` and implement `Ast`.

There is also NativeUndefined (to represent javascript `undefined`, as nil
can't be used).

These types are absent from the original coffeescript implementation, as js is
a much more dynamic language.
*/

import (
	"strings"
)

// Native* types implement Ast.
// They are used to represent array[Ast], int, string, map[string]Ast and
// the undefined value (a parsing function returns nil to indicate failure,
// undefined is something else).

type NativeArray []Ast

func NewNativeArray(a []Ast) *NativeArray {
	na := NativeArray(a)
	return &na
}

func NewEmptyNativeArray() *NativeArray {
	na := NativeArray([]Ast{})
	return &na
}

func (na *NativeArray) Get(i int) Ast { return (*na)[i] }
func (na *NativeArray) Length() int   { return len(*na) }
func (na *NativeArray) String() string {
	var b strings.Builder
	b.WriteString("[")
	first := true
	for _, it := range *na {
		if !first {
			b.WriteString(",")
		}
		b.WriteString(it.String())
		first = false
	}
	b.WriteString("]")
	return b.String()
}

func (na *NativeArray) Array() []Ast { return *na }
func (na *NativeArray) Append(it Ast) {
	*na = append(*na, it)
}

// `["a","","bc"]` -> `"abc"`
// (while we use strings in that example,
// a NativeArray will not contain strings directly
// but rather NativeString or children NativeArray,
// it however works the same)
func (na *NativeArray) Concat() NativeString {
	var b strings.Builder
	for _, element := range *na {
		switch v := element.(type) {
		case NativeString:
			b.WriteString(string(v))
		case *NativeArray:
			b.WriteString(string(v.Concat()))
		case NativeUndefined:
		default:
			b.WriteString(v.String())
		}
	}
	return NewNativeString(b.String())
}

func (na *NativeArray) SetLine(n int)                                   {}
func (na *NativeArray) GetLine() int                                    { return 1 }
func (na *NativeArray) SetOrigin(o Origin)                              {}
func (na *NativeArray) GetOrigin() Origin                               { return Origin{} }
func (na *NativeArray) HasAttribute(key interface{}) bool               { return false }
func (na *NativeArray) GetAttribute(key interface{}) interface{}        { return nil }
func (na *NativeArray) SetAttribute(key interface{}, value interface{}) { panic("N/A") }
