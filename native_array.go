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

type NativeArray struct {
	Array []Ast
}

func NewNativeArray(a []Ast) *NativeArray {
	if a == nil {
		// TODO try to remove this nil test
		return NewEmptyNativeArray()
	} else {
		return &NativeArray{a}
	}
}

func NewEmptyNativeArray() *NativeArray {
	return &NativeArray{[]Ast{}}
}

func (na *NativeArray) Get(i int) Ast { return na.Array[i] }
func (na *NativeArray) Length() int   { return len(na.Array) }
func (na *NativeArray) String() string {
	var b strings.Builder
	b.WriteString("[")
	first := true
	for _, it := range na.Array {
		if !first {
			b.WriteString(",")
		}
		b.WriteString(it.String())
		first = false
	}
	b.WriteString("]")
	return b.String()
}

func (na *NativeArray) Append(it Ast) {
	na.Array = append(na.Array, it)
}

// `["a","","bc"]` -> `"abc"`
// (with respect to the fact elements of the example are not strings
// but either NativeString or embedded *NativeArray)
func (na *NativeArray) Concat() string {
	var b strings.Builder
	for _, element := range na.Array {
		switch v := element.(type) {
		case NativeString:
			b.WriteString(v.Str)
		case *NativeArray:
			b.WriteString(v.Concat())
		case NativeUndefined:
		default:
			b.WriteString(v.String())
		}
	}
	return b.String()
}

func (na *NativeArray) SetLine(n int)                                   {}
func (na *NativeArray) GetLine() int                                    { return 1 }
func (na *NativeArray) SetOrigin(o Origin)                              {}
func (na *NativeArray) GetOrigin() Origin                               { return Origin{} }
func (na *NativeArray) HasAttribute(key interface{}) bool               { return false }
func (na *NativeArray) GetAttribute(key interface{}) interface{}        { return nil }
func (na *NativeArray) SetAttribute(key interface{}, value interface{}) { panic("N/A") }
