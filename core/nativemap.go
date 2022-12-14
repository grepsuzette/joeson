package core

import (
	"strings"
)

type NativeMap map[string]Astnode

func NewEmptyNativeMap() NativeMap { return NewNativeMap(map[string]Astnode{}) }
func NewNativeMap(h map[string]Astnode) NativeMap {
	// return NativeMap{map[string]Astnode}
	return h
}

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

func (nm NativeMap) GetExist(k string) (Astnode, bool) {
	v, exist := nm[k]
	return v, exist
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

func (nm NativeMap) ContentString() string {
	var b strings.Builder
	b.WriteString("NativeMap{")
	for k, v := range nm {
		b.WriteString("  " + k + ": " + v.ContentString())
	}
	b.WriteString("}")
	return b.String()
}
func (nm NativeMap) HandlesChildLabel() bool         { return false }
func (nm NativeMap) Labels() []string                { return []string{} }
func (nm NativeMap) Captures() []Astnode             { return []Astnode{} }
func (nm NativeMap) Prepare()                        {}
func (nm NativeMap) GetGNode() *GNode                { return nil }
func (nm NativeMap) Parse(ctx *ParseContext) Astnode { panic("uncallable") }

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (nm NativeMap) ForEachChild(f func(Astnode) Astnode) Astnode { return nm }
