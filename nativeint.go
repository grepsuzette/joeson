package main

import (
	"strconv"
)

// NativeInt and NativeString denote our need to express a terminal element
// that at the same time satisfies the `astnode` interface. Joeson.coffee used
// Number or string directly, but javascript is different language
type NativeInt int

func NewNativeIntFromString(s string) NativeInt {
	if n, e := strconv.Atoi(s); e != nil {
		return NativeInt{n}
	} else {
		panic(e)
	}
}

func (n NativeInt) ContentString() string { return Yellow(strconv.Itoa(n.int)) }
func (n NativeInt) HandlesChildLabel()    { return false }
func (n NativeInt) Labels() []string      { return n.GNode.Labels() }
func (n NativeInt) Captures() []astnode   { return n.GNode.Captures() }
func (n NativeInt) Prepare()              {}
func (n NativeInt) Parse(ctx *ParseContext) astnode {
	// there is nothing to parse... so this should not even ever be called?
	// how to distinguish whether it's a terminal or not?
	// TODO
	// maybe the grammar needs to check the type,
	//  NativeString or NativeInt are terminal
	panic("uncallable") // or nil? will see
}
