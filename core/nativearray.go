package core

import "grepsuzette/joeson/lambda"
import . "grepsuzette/joeson/colors"
import "strings"

type NativeArray struct {
	Array []Astnode
}

func NewNativeArray(a []Astnode) *NativeArray { return &NativeArray{a} }

func (na *NativeArray) Get(i int) Astnode       { return na.Array[i] }
func (na *NativeArray) Length() int             { return len(na.Array) }
func (na *NativeArray) GetGNode() *GNode        { return nil }
func (na *NativeArray) HandlesChildLabel() bool { return false }
func (na *NativeArray) Labels() []string        { return []string{} }
func (na *NativeArray) Captures() []Astnode     { return []Astnode{} }
func (na *NativeArray) Prepare()                {}
func (na *NativeArray) ContentString() string {
	return Blue("[") + strings.Join(lambda.Map(na.Array, func(x Astnode) string { return x.ContentString() }), Blue(", ")) + Blue("]")
}

func (na *NativeArray) Parse(ctx *ParseContext) Astnode {
	panic("uncallable") // or nil? will see
}

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (na *NativeArray) ForEachChild(f func(Astnode) Astnode) Astnode { return na }

// func (na *NativeArray) ForEachChild(f func(Astnode) Astnode) Astnode {
// 	a := []Astnode{}
// 	for _, child := range na.Array {
// 		if r := f(child); r != nil {
// 			a = append(a, r)
// 		} // remove from array when r != nil
// 	}
// 	na.Array = a
// 	return na
// }
