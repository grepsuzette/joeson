package core

import (
	"grepsuzette/joeson/lambda"
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
func (na *NativeArray) Labels() []string        { return []string{} }
func (na *NativeArray) Captures() []Ast         { return []Ast{} }
func (na *NativeArray) Prepare()                {}
func (na *NativeArray) ContentString() string {
	return "[" + strings.Join(lambda.Map(na.Array, func(x Ast) string { return x.ContentString() }), ",") + "]"
}

func (na *NativeArray) Parse(ctx *ParseContext) Ast { panic("uncallable") }

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (na *NativeArray) ForEachChild(f func(Ast) Ast) Ast { return na }
