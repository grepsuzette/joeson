package core

import (
	"grepsuzette/joeson/helpers"
	"strings"
)

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

func (na *NativeArray) Get(i int) Ast           { return na.Array[i] }
func (na *NativeArray) Length() int             { return len(na.Array) }
func (na *NativeArray) GetGNode() *GNode        { return nil }
func (na *NativeArray) HandlesChildLabel() bool { return false }
func (na *NativeArray) Prepare()                {}
func (na *NativeArray) ContentString() string {
	return "[" + strings.Join(helpers.AMap(na.Array, func(x Ast) string { return x.ContentString() }), ",") + "]"
}

func (na *NativeArray) Parse(ctx *ParseContext) Ast      { panic("uncallable") }
func (na *NativeArray) ForEachChild(f func(Ast) Ast) Ast { return na }
