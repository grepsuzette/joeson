package line

import (
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/core"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
)

type ILine struct {
	name    string // ILine, as terminal elements, are always named
	content Line
	attrs   core.ParseOptions
}

/*
I() is a variadic function which allows a variety of declarations, for example:
- I("INT", "/[0-9]+/")
- I("INT", "/[0-9]+/", func(it Astnode) Astnode { return new NativeInt(it) })
- I("INT", "/[0-9]+/", func(it Astnode, ctx *ParseContext) Astnode { return <...> })
- I("INT", "/[0-9]+/", func(it Astnode) Astnode { return <...> }, core.ParseOptions{SkipLog: false, SkipCache: true})
- I("RANGE", O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}"))))
  This one is a handcompiled rule with an O which the joeson grammar is initially defined as in ast/handcompiled
- I("LABEL", C(St('&'), St('@'), R("WORD"))),
  That one is a handcompiled rule that doesn't use an O rule.
*/
func I(a ...any) ILine {
	name, content, attrs := lineInit(a)
	if name == "" {
		panic("ILine must always be named")
	}
	return ILine{name, content, attrs}
}

func (il ILine) LineType() string { return "i" }
func (il ILine) Content() Line    { return il.content }
func (il ILine) StringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += il.LineType()
	s += " "
	s += il.content.StringIndent(nIndent)
	if il.attrs.CbBuilder != nil {
		s += Green(", ") + Yellow("𝘧")
	}
	return s
}

func (il ILine) ToRule(grammar *ast.Grammar, parentRule Ast) (name string, rule Ast) {
	return il.name, getRule(grammar, il.name, il.content, parentRule, il.attrs)
}
