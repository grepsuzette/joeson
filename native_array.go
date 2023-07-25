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
	// "reflect"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// Native* types implement Ast.
// They are used to represent array[Ast], int, string, map[string]Ast and
// the undefined value (a parsing function returns nil to indicate failure,
// undefined is something else).

type NativeArray struct {
	Attr
	Array []Ast
}

func NewNativeArray(a []Ast) *NativeArray {
	if a == nil {
		return &NativeArray{newAttr(), []Ast{}}
	} else {
		return &NativeArray{newAttr(), a}
	}
}

func (na *NativeArray) Get(i int) Ast { return na.Array[i] }
func (na *NativeArray) Length() int   { return len(na.Array) }
func (na *NativeArray) String() string {
	return "[" + strings.Join(helpers.AMap(na.Array, func(x Ast) string { return x.String() }), ",") + "]"
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
