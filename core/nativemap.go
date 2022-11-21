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
func (nm NativeMap) HandlesChildLabel() bool { return false }
func (nm NativeMap) Labels() []string        { return []string{} }
func (nm NativeMap) Captures() []Astnode     { return []Astnode{} }
func (nm NativeMap) Prepare()                {}
func (nm NativeMap) GetGNode() *GNode        { return nil }
func (nm NativeMap) Parse(ctx *ParseContext) Astnode {
	// there is nothing to parse... so this should not even ever be called?
	panic("uncallable") // or nil? will see
}

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (nm NativeMap) ForEachChild(f func(Astnode) Astnode) Astnode { return nm }

// func (nm NativeMap) ForEachChild(f func(Astnode) Astnode) Astnode {
// 	h := map[string]Astnode{}
// 	for k, child := range nm {
// 		if r := f(child); r != nil {
// 			h[k] = r
// 		} // removed from map when r is nil
// 	}
// 	nm = h
// 	return nm
// }
