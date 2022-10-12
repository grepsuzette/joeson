package main

import "grepsuzette/joeson/lambda"
import "strings"

// this was unnecessary in joeson.coffee because javascript is not strongly typed
type NativeArray struct {
	GNode
	Array []astnode
}

func NewNativeArray(a []astnode) NativeArray { return NativeArray{NewGNode(), a} }

func (na NativeArray) GetGNode() GNode     { return na.GNode }
func (na NativeArray) HandlesChildLabel()  { return false }
func (na NativeArray) Labels() []string    { return na.GNode.Labels() }
func (na NativeArray) Captures() []astnode { return na.GNode.Captures() }
func (na NativeArray) Prepare()            {}
func (na NativeArray) ContentString() string {
	return Blue("[") + strings.Join(Map(na, func(x) { return x.ContentString() }), Blue(", ")) + Blue("]")
}

func (na NativeArray) Parse(ctx *ParseContext) astnode {
	a := []astnode{}
	for _, x := range Array {
		r = x.Parse(ctx)
		if r != nil {
			a = append(a, r)
		}
	}
	return NativeArray{a}
}
