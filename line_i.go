package joeson

import (
	"grepsuzette/joeson/helpers"
)

type ILine struct {
	name    string // ILine, as terminal elements, are always named
	content Line
	attrs   ParseOptions
}

/*
I() is a variadic function which allows a variety of declarations, for example:
- I("INT", "/[0-9]+/")
- I("INT", "/[0-9]+/", func(it Ast) Ast { return new NativeInt(it) })
- I("INT", "/[0-9]+/", func(it Ast, ctx *ParseContext) Ast { return <...> })
- I("INT", "/[0-9]+/", func(it Ast) Ast { return <...> }, ParseOptions{SkipLog: false, SkipCache: true})
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
func (il ILine) Name() string     { return il.name }
func (il ILine) Content() Line    { return il.content }
func (il ILine) StringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += il.LineType()
	s += " "
	s += il.content.StringIndent(nIndent)
	if il.attrs.CbBuilder != nil {
		s += green(", ") + yellow("𝘧")
	}
	return s
}

func (il ILine) toRule(rank *Rank, parentRule Ast, opts TraceOptions, lazyGrammar *helpers.Lazy[*Grammar]) (name string, rule Ast) {
	return il.name, getRule(rank, il.name, il.content, parentRule, il.attrs, opts, lazyGrammar)
}
