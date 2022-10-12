package main

import "strconv"

// NativeInt and NativeString denote our need to express a terminal element
// that at the same time satisfies the `astnode` interface.
type NativeString struct {
	GNode
	str string
}

func NewNativeString(s string) NativeString {
	g := NewGNode()
	g.capture = false
	return NativeString{g, s}
}
func (ns NativeString) GetGNode() GNode       { return ns.GNode }
func (ns NativeString) HandlesChildLabel()    { return false }
func (ns NativeString) Labels() []string      { return ns.GNode.Labels() }
func (ns NativeString) Captures() []astnode   { return ns.GNode.Captures() }
func (ns NativeString) Prepare()              {}
func (ns NativeString) ContentString() string { return Red("\"") + ns.str + Red("\"") }
func (ns NativeString) Parse(ctx *ParseContext) astnode {
	// there is nothing to parse... so this should not even ever be called?
	// how to distinguish whether it's a terminal or not?
	// TODO
	// maybe the grammar needs to check the type,
	//  NativeString or NativeInt are terminal
	panic("uncallable") // or nil? will see
}
