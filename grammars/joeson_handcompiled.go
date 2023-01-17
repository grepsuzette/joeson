package grammars

import (
	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/core"
	. "grepsuzette/joeson/line"
	"strings"
)

// Instantiate a new joeson grammar.
// It can be used to parse all joeson grammars.
func NewJoeson() *Grammar {
	return NewJoesonWithOptions(DefaultTraceOptions())
}

func NewJoesonWithOptions(opts TraceOptions) *Grammar {
	opts = CheckEnvironmentForTraceOptions(opts)
	return GrammarFromLinesWithOptions(
		JOESON_GRAMMAR_NAME,
		JOESON_GRAMMAR_RULES(),
		NewEmptyGrammarWithOptions(opts),
		opts,
	)
}

func C(a ...Ast) Ast               { return NewChoice(NewNativeArray(a)) }
func E(it Ast) Ast                 { return NewExistential(it) }
func L(label string, node Ast) Ast { node.GetGNode().Label = label; return node }
func N(it Ast) Ast                 { return NewNot(it) }
func R(s string) *Ref              { return NewRef(NewNativeString(s)) }
func Re(s string) *Regex           { return NewRegexFromString(s) }
func S(a ...Ast) *Sequence         { return NewSequence(NewNativeArray(a)) }
func St(s string) Str              { return NewStr(s) }
func P(value, join Ast, minmax ...int) *Pattern {
	min := -1
	max := -1
	if len(minmax) > 0 {
		min = minmax[0]
		if len(minmax) > 1 {
			max = minmax[1]
		}
	}
	p := NewPattern(NewNativeMap(map[string]Ast{
		"value": value,
		"join":  join,
		"min":   NewNativeInt(min),
		"max":   NewNativeInt(max),
	}))
	return p
}

func StringFromNativeArray(it Ast) string {
	var b strings.Builder
	na := it.(*NativeArray)
	for _, ns := range na.Array {
		b.WriteString(ns.(NativeString).Str)
	}
	return b.String()
}

const JoesonNbRules int = 35
const JOESON_GRAMMAR_NAME = "__grammar__"

func JOESON_GRAMMAR_RULES() Lines {
	return []Line{
		o(Named("EXPR", rules(
			o(S(R("CHOICE"), R("_"))),
			o(Named("CHOICE", rules(
				o(S(P(R("_PIPE"), nil), P(R("SEQUENCE"), R("_PIPE"), 2), P(R("_PIPE"), nil)), func(it Ast) Ast { return NewChoice(it) }),
				o(Named("SEQUENCE", rules(
					o(P(R("UNIT"), nil, 2), func(it Ast) Ast { return NewSequence(it) }),
					o(Named("UNIT", rules(
						o(S(R("_"), R("LABELED"))),
						o(Named("LABELED", rules(
							o(S(E(S(L("label", R("LABEL")), St(":"))), L("&", C(R("DECORATED"), R("PRIMARY"))))),
							o(Named("DECORATED", rules(
								o(S(R("PRIMARY"), St("?")), func(it Ast) Ast { return NewExistential(it) }),
								o(S(L("value", R("PRIMARY")), St("*"), L("join", E(S(N(R("__")), R("PRIMARY")))), L("@", E(R("RANGE")))), func(it Ast) Ast { return NewPattern(it) }),
								o(S(L("value", R("PRIMARY")), St("+"), L("join", E(S(N(R("__")), R("PRIMARY"))))), func(it Ast) Ast {
									h := it.(NativeMap)
									h.Set("Min", NewNativeInt(1))
									h.Set("Max", NewNativeInt(-1))
									return NewPattern(h)
								}),
								o(S(L("value", R("PRIMARY")), L("@", R("RANGE"))), func(it Ast) Ast { return NewPattern(it) }),
								o(S(St("!"), R("PRIMARY")), func(it Ast) Ast { return NewNot(it) }),
								o(C(S(St("(?"), L("expr", R("EXPR")), St(")")), S(St("?"), L("expr", R("EXPR")))), func(it Ast) Ast { return NewLookahead(it) }),
								i(Named("RANGE", o(S(St("{"), R("_"), L("min", E(R("INT"))), R("_"), St(","), R("_"), L("max", E(R("INT"))), R("_"), St("}"))))),
							))),
							o(Named("PRIMARY", rules(
								o(S(R("WORD"), St("("), R("EXPR"), St(")")), func(it Ast) Ast {
									na := it.(*NativeArray)
									if na.Length() != 4 {
										panic("assert")
									}
									return NewRef(NewNativeArray([]Ast{na.Get(1), na.Get(3)}))
								}),
								o(R("WORD"), func(it Ast) Ast { return NewRef(it) }),
								o(S(St("("), L("inlineLabel", E(S(R("WORD"), St(": ")))), L("expr", R("EXPR")), St(")"), E(S(R("_"), St("->"), R("_"), L("code", R("CODE"))))), fCode),
								i(Named("CODE", o(S(St("{"), P(S(N(St("}")), C(R("ESC1"), R("."))), nil, -1, -1), St("}")))), fCode),
								o(S(St("'"), P(S(N(St("'")), C(R("ESC1"), R("."))), nil), St("'")), func(it Ast) Ast { return NewStr(StringFromNativeArray(it)) }),
								o(S(St("/"), P(S(N(St("/")), C(R("ESC2"), R("."))), nil), St("/")), func(it Ast) Ast { return NewRegexFromString(StringFromNativeArray(it)) }),
								o(S(St("["), P(S(N(St("]")), C(R("ESC2"), R("."))), nil), St("]")), func(it Ast) Ast { return NewRegexFromString("[" + StringFromNativeArray(it) + "]") }),
							))),
						))),
					))),
				))),
			))),
		))),
		i(Named("LABEL", C(St("&"), St("@"), R("WORD")))),
		i(Named("WORD", Re("[a-zA-Z\\._][a-zA-Z\\._0-9]*"))),
		i(Named("INT", Re("[0-9]+")), func(it Ast) Ast { return NewNativeIntFrom(it) }),
		i(Named("_PIPE", S(R("_"), St("|")))),
		i(Named("_", P(C(St(" "), St("\n")), nil))),
		i(Named("__", P(C(St(" "), St("\n")), nil, 1))),
		i(Named(".", Re("[\\s\\S]"))),
		i(Named("ESC1", S(St("\\"), R(".")))),
		i(Named("ESC2", S(St("\\"), R("."))), func(chr Ast) Ast { return NewNativeString("\\" + chr.(NativeString).Str) }),
		// i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true}, func(it Astnode, ctx *ParseContext) Astnode { return it })),
	}
}
