package core

import (
	. "grepsuzette/joeson/colors"
	"strconv"
)

// NativeInt and NativeString denote our need to express a terminal element
// that at the same time satisfies the `Astnode` interface. Joeson.coffee used
// Number or string directly, but javascript is different language
type NativeInt int

func NewNativeInt(n int) NativeInt { return NativeInt(n) }
func NewNativeIntFromString(s string) NativeInt {
	if n, e := strconv.Atoi(s); e != nil {
		return NativeInt(n)
	} else {
		panic(e)
	}
}
func NewNativeIntFromNativeString(ns NativeString) NativeInt {
	if n, e := strconv.Atoi(ns.Str); e != nil {
		return NativeInt(n)
	} else {
		panic(e)
	}
}

func (n NativeInt) String() string          { return strconv.Itoa(int(n)) }
func (n NativeInt) ContentString() string   { return Yellow(strconv.Itoa(int(n))) }
func (n NativeInt) HandlesChildLabel() bool { return false }
func (n NativeInt) GetGNode() *GNode        { return nil }
func (n NativeInt) Labels() []string        { return []string{} }
func (n NativeInt) Captures() []Astnode     { return []Astnode{} }
func (n NativeInt) Prepare()                {}
func (n NativeInt) Parse(ctx *ParseContext) Astnode {
	// there is nothing to parse... so this should not even ever be called?
	// how to distinguish whether it's a terminal or not?
	// TODO
	// maybe the grammar needs to check the type,
	//  NativeString or NativeInt are terminal
	panic("uncallable") // or nil? will see
}

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (n NativeInt) ForEachChild(f func(Astnode) Astnode) Astnode { return n }
