package ast

import (
	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/core"
	// "grepsuzette/joeson/helpers"
	. "grepsuzette/joeson/line"
	"strings"
)

// Use NewJoeson() to instantiate a new joeson grammar.
// This is the hand-compiled joeson grammar
// it comes from the original joeson.coffee

// {{{1 imports and funcs
func C(a ...Ast) Ast {
	r := NewChoice(NewNativeArray(a))
	s := ""
	for _, v := range a {
		s += v.ContentString() + ","
	}
	// fmt.Println("ohm Choice in=" + s + " out=" + r.ContentString())
	return r
}
func E(it Ast) Ast { return NewExistential(it) }
func L(label string, node Ast) Ast {
	// fmt.Println("ohm " + reflect.TypeOf(node).String() + " -  " + label)
	node.GetGNode().Label = label
	return node
}
func N(it Ast) Ast {
	// fmt.Println("ohm Not it=" + it.ContentString())
	return NewNot(it)
}
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
func R(s string) *Ref {
	r := NewRef(NewNativeString(s))
	return r
}
func Re(s string) *Regex { return NewRegexFromString(s) }
func S(a ...Ast) *Sequence {
	r := NewSequence(NewNativeArray(a))
	s := ""
	for _, v := range a {
		s += Prefix(v) + v.ContentString() + ","
	}
	return r
}
func St(s string) Str {
	r := NewStr(s)
	return r
}

func AttemptToJoinANativeArrayOrPanic(it Ast) string {
	var b strings.Builder
	na := it.(*NativeArray)
	for _, ns := range na.Array {
		b.WriteString(ns.(NativeString).Str)
	}
	return b.String()
}

func o(a ...any) OLine          { return O(a...) }
func i(a ...any) ILine          { return I(a...) }
func rules(lines ...Line) ALine { return NewALine(lines) }

// }}}1

const JoesonNbRules int = 35

const QUOTE string = "'\\''"
const JOESON_GRAMMAR_NAME = "__grammar__"

func NewJoeson() *Grammar {
	return NewGrammarFromLines(
		JOESON_GRAMMAR_NAME,
		JOESON_GRAMMAR_RULES(),
		NewEmptyGrammarNamed("empty grammar"),
	)
}

func JOESON_GRAMMAR_RULES() Lines {
	return []Line{
		o(Named("EXPR", rules(
			o(S(R("CHOICE"), R("_"))),
			o(Named("CHOICE", rules(
				o(S(P(R("_PIPE"), nil), P(R("SEQUENCE"), R("_PIPE"), 2), P(R("_PIPE"), nil)), func(it Ast) Ast {
					return NewChoice(it)
				}),
				o(Named("SEQUENCE", rules(
					o(P(R("UNIT"), nil, 2), func(it Ast) Ast { return NewSequence(it) }),
					o(Named("UNIT", rules(
						o(S(R("_"), R("LABELED"))),
						o(Named("LABELED", rules(
							o(S(E(S(L("label", R("LABEL")), St(":"))), L("&", C(R("DECORATED"), R("PRIMARY"))))),
							o(Named("DECORATED", rules(
								o(S(R("PRIMARY"), St("?")), func(it Ast) Ast { return NewExistential(it) }),
								o(S(L("value", R("PRIMARY")), St("*"), L("join", E(S(N(R("__")), R("PRIMARY")))), L("@", E(R("RANGE")))), func(it Ast) Ast {
									// fmt.Println("ohm DECORATED[1] value:PRIMARY* join:(!__ PRIMARY)? @:RANGE? NewPattern from it=" + it.ContentString())
									return NewPattern(it)
								}),
								o(S(L("value", R("PRIMARY")), St("+"), L("join", E(S(N(R("__")), R("PRIMARY"))))), func(it Ast) Ast {
									// fmt.Println("ohm DECORATED[2] value:PRIMARY+ join:(!__ PRIMARY)? @:RANGE? NewPattern from it=" + it.ContentString())
									h := it.(NativeMap)
									h.Set("Min", NewNativeInt(1))
									h.Set("Max", NewNativeInt(-1))
									return NewPattern(h)
								}),
								o(S(L("value", R("PRIMARY")), L("@", R("RANGE"))), func(it Ast) Ast {
									// fmt.Println("ohm DECORATED[3] value:PRIMARY @:RANGE NewPattern from it=" + it.ContentString())
									return NewPattern(it)
								}),
								o(S(St("!"), R("PRIMARY")), func(it Ast) Ast { return NewNot(it) }),
								o(C(S(St("(?"), L("expr", R("EXPR")), St(")")), S(St("?"), L("expr", R("EXPR")))), func(it Ast) Ast { return NewLookahead(it) }),
								i(Named("RANGE", o(S(St("{"), R("_"), L("min", E(R("INT"))), R("_"), St(","), R("_"), L("max", E(R("INT"))), R("_"), St("}"))))),
							))),
							o(Named("PRIMARY", rules(
								o(S(R("WORD"), St("("), R("EXPR"), St(")")), func(it Ast) Ast {
									na := it.(*NativeArray)
									if na.Length() != 4 {
										panic("logic")
									}
									return NewRef(NewNativeArray([]Ast{na.Get(1), na.Get(3)}))
								}),
								o(R("WORD"), func(it Ast) Ast { return NewRef(it) }),
								o(S(St("("), L("inlineLabel", E(S(R("WORD"), St(": ")))), L("expr", R("EXPR")), St(")"), E(S(R("_"), St("->"), R("_"), L("code", R("CODE"))))), func(it Ast) Ast {
									h := it.(NativeMap)
									if h.Get("code") != nil {
										panic("code in joeson deprecated")
									}
									return h.Get("expr")
								}),
								// i "CODE", o S(St("{"), P(S(N(St("}")), C(R("ESC1"), R(".")))), St("}")), (it) -> require('./joescript').parse(it.join '')),
								i(Named("CODE", o(S(St("{"), P(S(N(St("}")), C(R("ESC1"), R("."))), nil, -1, -1), St("}")))), func(it Ast) Ast {
									// TODO condense it or write a function
									// deprecated code in joeson
									switch v := it.(type) {
									case NativeMap:
										if v.Get("code") != nil {
											panic("code in joeson is obsolete")
										}
										expr := v.Get("expr")
										return expr
									default:
										panic("assert")
									}
								}),
								// Note:this v--- P() here *seems* useless. It's a Pattern(value, join=nil, min=-1, max=-1). However, that kind of pattern has a @capture = @value.capture, which allows it to capture the string.  at least that's my guess atm
								o(S(St("'"), P(S(N(St("'")), C(R("ESC1"), R("."))), nil), St("'")), func(it Ast) Ast {
									// this will require heavy testings, best case
									// scenario is captures got into a NativeArray
									// of NativeString, there are very few chances
									// we get it right at first though; better get
									// prepared.
									return NewStr(AttemptToJoinANativeArrayOrPanic(it))
								}),
								o(S(St("/"), P(S(N(St("/")), C(R("ESC2"), R("."))), nil), St("/")), func(it Ast) Ast { return NewRegexFromString(AttemptToJoinANativeArrayOrPanic(it)) }),
								o(S(St("["), P(S(N(St("]")), C(R("ESC2"), R("."))), nil), St("]")), func(it Ast) Ast { return NewRegexFromString("[" + AttemptToJoinANativeArrayOrPanic(it) + "]") }),
							))),
						))),
					))),
				))),
			))),
		))),
		i(Named("LABEL", C(St("&"), St("@"), R("WORD")))),
		i(Named("WORD", Re("[a-zA-Z\\._][a-zA-Z\\._0-9]*"))),
		i(Named("INT", Re("[0-9]+")), func(it Ast) Ast {
			return NewNativeIntFrom(it)
		}),
		i(Named("_PIPE", S(R("_"), St("|")))),
		i(Named("_", P(C(St(" "), St("\n")), nil))),
		i(Named("__", P(C(St(" "), St("\n")), nil, 1))),
		i(Named(".", Re("[\\s\\S]"))),
		i(Named("ESC1", S(St("\\"), R(".")))),
		i(Named("ESC2", S(St("\\"), R("."))), func(chr Ast) Ast { return NewNativeString("\\" + chr.(NativeString).Str) }),
		// i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true},
		// 	func(it Astnode, ctx *ParseContext) Astnode {
		// 		// ctx.SkipLog is false
		// 		// ctx.SkipCache is true
		// 		// ctx.Debug is false
		// 		return nil
		// })),
	}
}

// vim: fdm=marker fdl=0
