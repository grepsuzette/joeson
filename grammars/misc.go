package grammars

import (
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/line"
)

func o(a ...any) line.OLine               { return line.O(a...) }
func i(a ...any) line.ILine               { return line.I(a...) }
func rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }
func named(name string, lineStringOrAstnode any) line.NamedRule {
	return line.Named(name, lineStringOrAstnode)
}

func fCode(it Ast) Ast {
	h := it.(NativeMap)
	if !h.IsUndefined("code") {
		panic("code in joeson is obsolete")
	}
	return h.GetOrPanic("expr")
}
