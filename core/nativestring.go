package core

import . "grepsuzette/joeson/colors"

// NativeInt and NativeString denote our need to express a terminal element
// that at the same time satisfies the `Astnode` interface.
type NativeString struct {
	Str string
}

func NewNativeString(s string) NativeString {
	return NativeString{s}
}
func (ns NativeString) GetGNode() *GNode                { return nil }
func (ns NativeString) HandlesChildLabel() bool         { return false }
func (ns NativeString) Labels() []string                { return []string{} }
func (ns NativeString) Captures() []Astnode             { return []Astnode{} }
func (ns NativeString) Prepare()                        {}
func (ns NativeString) ContentString() string           { return Red(`"`) + ns.Str + Red(`"`) }
func (ns NativeString) Parse(ctx *ParseContext) Astnode { panic("uncallable") }

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (n NativeString) ForEachChild(f func(Astnode) Astnode) Astnode { return n }
