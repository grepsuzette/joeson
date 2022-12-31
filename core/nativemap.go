package core

import (
	"grepsuzette/joeson/helpers"
	"strings"
)

type NativeMap map[string]Astnode

func NewEmptyNativeMap() NativeMap { return NewNativeMap(map[string]Astnode{}) }
func NewNativeMap(h map[string]Astnode) NativeMap {
	// return NativeMap{map[string]Astnode}
	return h
}

func (nm NativeMap) HandlesChildLabel() bool         { return false }
func (nm NativeMap) Labels() []string                { return []string{} }
func (nm NativeMap) Captures() []Astnode             { return []Astnode{} }
func (nm NativeMap) Prepare()                        {}
func (nm NativeMap) GetGNode() *GNode                { return nil }
func (nm NativeMap) Parse(ctx *ParseContext) Astnode { panic("uncallable") }

func (nm NativeMap) ContentString() string {
	var b strings.Builder
	b.WriteString("NativeMap{")
	first := true
	// for k, v := range nm {
	for _, k := range helpers.SortStringKeys(nm) {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(k + ":" + nm.Get(k).ContentString())
		first = false
	}
	b.WriteString("}")
	return b.String()
}

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (nm NativeMap) ForEachChild(f func(Astnode) Astnode) Astnode { return nm }

func (nm NativeMap) Keys() []string {
	a := []string{}
	for k := range nm {
		a = append(a, k)
	}
	return a
}

func (nm NativeMap) IsEmpty() bool {
	for range nm {
		return true
	}
	return false
}

func (nm NativeMap) GetExists(k string) (Astnode, bool) {
	v, exists := nm[k]
	return v, exists
}

// specialized getter when value is known to be a NativeString
// it panics otherwise
func (nm NativeMap) GetStringExists(k string) (string, bool) {
	if v, exists := nm[k]; exists {
		return v.(NativeString).Str, true
	} else {
		return "", false
	}
}

// specialized getter when value is known to be a NativeInt
// it panics otherwise
func (nm NativeMap) GetIntExists(k string) (int, bool) {
	if v, exists := nm[k]; exists {
		return int(v.(NativeInt)), true
	} else {
		return 0, false
	}
}
func (nm NativeMap) Get(k string) Astnode {
	return nm[k]
}

func (nm NativeMap) Set(k string, v Astnode) {
	nm[k] = v
}

func (nm NativeMap) Has(k string) bool {
	_, ok := nm[k]
	return ok
}
