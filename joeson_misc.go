package joeson

import (
	"strings"
)

func o(a ...any) OLine { return O(a...) }
func i(a ...any) ILine { return I(a...) }

func rules(lines ...Line) ALine { return NewALine(lines) }
func named(name string, lineStringOrAst any) NamedRule {
	return Named(name, lineStringOrAst)
}

func fCode(it Ast) Ast {
	h := it.(NativeMap)
	if !h.IsUndefined("code") {
		panic("code in joeson is obsolete")
	}
	return h.GetOrPanic("expr")
}

func stringFromNativeArray(it Ast) string {
	var b strings.Builder
	na := it.(*NativeArray)
	for _, ns := range na.Array {
		b.WriteString(ns.(NativeString).Str)
	}
	return b.String()
}
