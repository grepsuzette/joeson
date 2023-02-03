package joeson

/*
"Native*" types
===============

Short of finding a better name, these types wrap int, string, map, array to
satisfy interface core.Ast.

There is also NativeUndefined (to represent javascript `undefined`, as nil
can't be used).

These types are absent from the original coffeescript implementation, as js is
a much more dynamic language.

They are put here in a core/native package for 2 reasons:

1. core.ParseContext depends upon it,
2. it better separates the ast types from the original implementation in ast/
*/

import (
	"github.com/grepsuzette/joeson/helpers"
	"strings"
)

// Native* types satisfy Ast, but their GetGNode() returns nil.
// They are used to represent array[Ast], int, string, map[string]Ast and
// the undefined value (a parsing function returns nil to indicate failure,
// undefined is something else).

type NativeArray struct {
	Array []Ast
}

func NewNativeArray(a []Ast) *NativeArray {
	if a == nil {
		return &NativeArray{[]Ast{}}
	} else {
		return &NativeArray{a}
	}
}

func (na *NativeArray) Get(i int) Ast { return na.Array[i] }
func (na *NativeArray) Length() int   { return len(na.Array) }
func (na *NativeArray) ContentString() string {
	return "[" + strings.Join(helpers.AMap(na.Array, func(x Ast) string { return x.ContentString() }), ",") + "]"
}
