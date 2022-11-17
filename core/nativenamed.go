package core

import (
//import . "grepsuzette/joeson/colors"
)

// NativeNamed is simply a key:value pair that
// also satisfies Astnode.
type NativeNamed struct {
	Name  string
	Value Astnode
}

func NewNativeNamed(name string, value Astnode) NativeNamed { return NativeNamed{name, value} }

func (nn NativeNamed) String() string {
	return "named:" + nn.Name + ";value:" + nn.Value.ContentString()
}
func (nn NativeNamed) ContentString() string   { return nn.String() }
func (nn NativeNamed) HandlesChildLabel() bool { return false }
func (nn NativeNamed) GetGNode() *GNode        { return nil }
func (nn NativeNamed) Labels() []string        { return []string{} }
func (nn NativeNamed) Captures() []Astnode     { return []Astnode{} }
func (nn NativeNamed) Prepare()                {}
func (nn NativeNamed) Parse(ctx *ParseContext) Astnode {
	panic("uncallable") // or nil? will see
}

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (nn NativeNamed) ForEachChild(f func(Astnode) Astnode) Astnode { return nn }
