package joeson

import "strings"

// NewJoeson() creates a new instance of the handcompiled grammar.
// You may want to use this with NewGrammarFromLines(), though
// with default options NewJoeson() is used anyway by default
// (best example probably in TestBootstrap())
func NewJoeson() *Grammar {
	return NewJoesonWithOptions(DefaultTraceOptions())
}

func NewJoesonWithOptions(opts TraceOptions) *Grammar {
	gm := GrammarFromLines(
		JoesonRules(),
		JoesonGrammarName,
		GrammarOptions{
			TraceOptions: opts,
			LazyGrammar:  nil,
		},
	)
	return gm
}

func c(a ...Ast) *choice                 { return newChoice(NewNativeArray(a)) }
func e(it Ast) *existential              { return newExistential(it) }
func l(label string, node Parser) Parser { node.SetLabel(label); return node }
func n(it Parser) *not                   { return newNot(it) }
func r(s string) *ref                    { return newRef(NewNativeString(s)) }
func re(s string) *regex                 { return newRegexFromString(s) }
func s(a ...Ast) *sequence               { return newSequence(NewNativeArray(a)) }
func st(s string) str                    { return newStr(s) }
func p(value, join Parser, minmax ...int) *pattern {
	min := -1
	max := -1
	if len(minmax) > 0 {
		min = minmax[0]
		if len(minmax) > 1 {
			max = minmax[1]
		}
	}
	p := newPattern(NewNativeMap(map[string]Ast{
		"value": value,
		"join":  join,
		"min":   NewNativeInt(min),
		"max":   NewNativeInt(max),
	}))
	return p
}

const JoesonNbRules int = 35
const JoesonGrammarName = "__joeson__"

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

// provide the Lines of the joeson grammar
func JoesonRules() []Line {
	return []Line{
		o(Named("EXPR", rules(
			o(s(r("CHOICE"), r("_"))),
			o(Named("CHOICE", rules(
				o(s(p(r("_PIPE"), nil), p(r("SEQUENCE"), r("_PIPE"), 2), p(r("_PIPE"), nil)), func(it Ast) Ast { return newChoice(it) }),
				o(Named("SEQUENCE", rules(
					o(p(r("UNIT"), nil, 2), func(it Ast) Ast { return newSequence(it) }),
					o(Named("UNIT", rules(
						o(s(r("_"), r("LABELED"))),
						o(Named("LABELED", rules(
							o(s(e(s(l("label", r("LABEL")), st(":"))), l("&", c(r("DECORATED"), r("PRIMARY"))))),
							o(Named("DECORATED", rules(
								o(s(r("PRIMARY"), st("?")), func(it Ast) Ast { return newExistential(it) }),
								o(s(l("value", r("PRIMARY")), st("*"), l("join", e(s(n(r("__")), r("PRIMARY")))), l("@", e(r("RANGE")))), func(it Ast) Ast { return newPattern(it) }),
								o(s(l("value", r("PRIMARY")), st("+"), l("join", e(s(n(r("__")), r("PRIMARY"))))), func(it Ast) Ast {
									h := it.(NativeMap)
									h.Set("min", NewNativeInt(1))
									h.Set("max", NewNativeInt(-1))
									return newPattern(h)
								}),
								o(s(l("value", r("PRIMARY")), l("@", r("RANGE"))), func(it Ast) Ast { return newPattern(it) }),
								o(s(st("!"), r("PRIMARY")), func(it Ast) Ast { return newNot(it) }),
								o(c(s(st("(?"), l("expr", r("EXPR")), st(")")), s(st("?"), l("expr", r("EXPR")))), func(it Ast) Ast { return newLookahead(it) }),
								i(Named("RANGE", o(s(st("{"), r("_"), l("min", e(r("INT"))), r("_"), st(","), r("_"), l("max", e(r("INT"))), r("_"), st("}"))))),
							))),
							o(Named("PRIMARY", rules(
								o(s(r("WORD"), st("("), r("EXPR"), st(")")), func(it Ast) Ast { return newRef(it) }), // it is a NativeArray of [r("WORD"), r("EXPR")]
								o(r("WORD"), func(it Ast) Ast { return newRef(it) }),
								o(s(st("("), l("inlineLabel", e(s(r("WORD"), st(": ")))), l("expr", r("EXPR")), st(")"), e(s(r("_"), st("->"), r("_"), l("code", r("CODE"))))), fCode),
								i(Named("CODE", o(s(st("{"), p(s(n(st("}")), c(r("ESC1"), r("."))), nil, -1, -1), st("}")))), fCode),
								o(s(st("'"), p(s(n(st("'")), c(r("ESC1"), r("."))), nil), st("'")), func(it Ast) Ast { return newStr(stringFromNativeArray(it)) }),
								o(s(st("/"), p(s(n(st("/")), c(r("ESC2"), r("."))), nil), st("/")), func(it Ast) Ast { return newRegexFromString(stringFromNativeArray(it)) }),
								o(s(st("["), p(s(n(st("]")), c(r("ESC2"), r("."))), nil), st("]")), func(it Ast) Ast { return newRegexFromString("[" + stringFromNativeArray(it) + "]") }),
							))),
						))),
					))),
				))),
			))),
		))),
		i(Named("LABEL", c(st("&"), st("@"), r("WORD")))),
		i(Named("WORD", re("[a-zA-Z\\._][a-zA-Z\\._0-9]*"))),
		i(Named("INT", re("[0-9]+")), func(it Ast) Ast { return NewNativeIntFrom(it) }),
		i(Named("_PIPE", s(r("_"), st("|")))),
		i(Named("_", p(c(st(" "), st("\n")), nil))),
		i(Named("__", p(c(st(" "), st("\n")), nil, 1))),
		i(Named(".", re("[\\s\\S]"))),
		i(Named("ESC1", s(st("\\"), r(".")))),
		i(Named("ESC2", s(st("\\"), r("."))), func(chr Ast) Ast { return NewNativeString("\\" + chr.(NativeString).Str) }),
		// i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true}, func(it Ast, ctx *ParseContext) Ast { return it })),
	}
}
