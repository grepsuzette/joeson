package core

import "grepsuzette/joeson/lambda"

// import . "grepsuzette/joeson/colors"
import "strings"

type NativeArray struct {
	Array []Astnode
}

func NewNativeArray(a []Astnode) *NativeArray {
	if a == nil {
		return &NativeArray{[]Astnode{}}
	} else {
		return &NativeArray{a}
	}
}

func (na *NativeArray) Get(i int) Astnode       { return na.Array[i] }
func (na *NativeArray) Length() int             { return len(na.Array) }
func (na *NativeArray) GetGNode() *GNode        { return nil }
func (na *NativeArray) HandlesChildLabel() bool { return false }
func (na *NativeArray) Labels() []string        { return []string{} }
func (na *NativeArray) Captures() []Astnode     { return []Astnode{} }
func (na *NativeArray) Prepare()                {}
func (na *NativeArray) ContentString() string {
	return "[" + strings.Join(lambda.Map(na.Array, func(x Astnode) string { return x.ContentString() }), ",") + "]"
}

func (na *NativeArray) Parse(ctx *ParseContext) Astnode { panic("uncallable") }

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (na *NativeArray) ForEachChild(f func(Astnode) Astnode) Astnode { return na }
