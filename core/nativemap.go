package core

import (
	"grepsuzette/joeson/helpers"
	"strings"
)

type NativeMap map[string]Ast

func NewEmptyNativeMap() NativeMap            { return NewNativeMap(map[string]Ast{}) }
func NewNativeMap(h map[string]Ast) NativeMap { return h }

func (nm NativeMap) HandlesChildLabel() bool     { return false }
func (nm NativeMap) Prepare()                    {}
func (nm NativeMap) GetGNode() *GNode            { return nil }
func (nm NativeMap) Parse(ctx *ParseContext) Ast { panic("uncallable") }

func (nm NativeMap) ContentString() string {
	var b strings.Builder
	b.WriteString("NativeMap{")
	first := true
	for _, k := range helpers.SortStringKeys(nm) {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(k + ":" + nm.GetOrPanic(k).ContentString())
		first = false
	}
	b.WriteString("}")
	return b.String()
}

func (nm NativeMap) ForEachChild(f func(Ast) Ast) Ast { return nm }

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

func (nm NativeMap) GetExists(k string) (Ast, bool) {
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

// true when key is not defined or when its value is NativeUndefined.
// Note: successfully parsed Ast can't possibly return nil.
func (nm NativeMap) IsUndefined(k string) bool {
	if v, exists := nm[k]; exists {
		if _, ok := v.(NativeUndefined); ok {
			return true
		} else {
			return false
		}
	}
	return true
}

func (nm NativeMap) GetOrPanic(k string) Ast {
	if r, ok := nm[k]; ok {
		return r
	} else {
		panic("assert")
	}
}

func (nm NativeMap) Get(k string) Ast {
	return nm[k]
}

func (nm NativeMap) Set(k string, v Ast) {
	nm[k] = v
}

func (nm NativeMap) Exists(k string) bool {
	_, ok := nm[k]
	return ok
}
