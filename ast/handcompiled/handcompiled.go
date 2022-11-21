package ast

import (
	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/core"
	. "grepsuzette/joeson/line"
	"strings"
)

// {{{1 imports and funcs
// this is the hand-compiled joeson grammar
// it comes from the original joeson.coffee

// the ast files alone won't do anything.
// for joeson to parse joeson grammars, its
// own grammar needs to be written somewhere.

// this is where. We will use line package
// to follow original coffee impl. as much
// as possible

func C(a ...Astnode) Astnode { return NewChoice(NewNativeArray(a)) }
func E(it Astnode) Astnode   { return NewExistential(it) }
func L(label string, node Astnode) Astnode {
	// fmt.Println(reflect.TypeOf(node).String() + " -  " + label)
	node.GetGNode().Label = label
	return node
}
func N(it Astnode) Astnode { return NewNot(it) }
func P(value, join Astnode, minmax ...int) *Pattern {
	min := -1
	max := -1
	if len(minmax) > 0 {
		min = minmax[0]
		if len(minmax) > 1 {
			max = minmax[1]
		}
	}
	return NewPattern(NewNativeMap(map[string]Astnode{
		"value": value,
		"join":  join,
		"min":   NewNativeInt(min),
		"max":   NewNativeInt(max),
	}))
}
func R(s string) *Ref          { return NewRef(NewNativeString(s)) }
func Re(s string) *Regex       { return NewRegexFromString(s) }
func S(a ...Astnode) *Sequence { return NewSequence(NewNativeArray(a)) }
func St(s string) Str          { return NewStr(s) }

func AttemptToJoinANativeArrayOrPanic(it Astnode) string {
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

const JoesonNbRules int = 33

var QUOTE string = "'\\''"
var JOESON_GRAMMAR_RULES Lines = []Line{
	o(Named("EXPR", rules(
		o(S(R("CHOICE"), R("_"))),
		o(Named("CHOICE", rules(
			o(S(P(R("_PIPE"), nil), P(R("SEQUENCE"), R("_PIPE"), 2), P(R("_PIPE"), nil)), func(it Astnode) Astnode { return NewChoice(it) }),
			o(Named("SEQUENCE", rules(
				o(P(R("UNIT"), nil, 2), func(it Astnode) Astnode { return NewSequence(it) }),
				o(Named("UNIT", rules(
					o(S(R("_"), R("LABELED"))),
					o(Named("LABELED", rules(
						o(S(E(S(L("label", R("LABEL")), St(":"))), L("&", C(R("DECORATED"), R("PRIMARY"))))),
						o(Named("DECORATED", rules(
							o(S(R("PRIMARY"), St("?")), func(it Astnode) Astnode { return NewExistential(it) }),
							o(S(L("value", R("PRIMARY")), St("*"), L("join", E(S(N(R("__")), R("PRIMARY")))), L("@", E(R("RANGE")))), func(it Astnode) Astnode { return NewPattern(it) }),
							o(S(L("value", R("PRIMARY")), St("+"), L("join", E(S(N(R("__")), R("PRIMARY"))))), func(it Astnode) Astnode {
								h := it.(NativeMap)
								h.Set("Min", NewNativeInt(1))
								h.Set("Max", NewNativeInt(-1))
								return NewPattern(h)
							}),
							o(S(L("value", R("PRIMARY")), L("@", R("RANGE"))), func(it Astnode) Astnode { return NewPattern(it) }), // note: the @ label will "source" and "import" the labels from RANGE node into `it`
							o(S(St("!"), R("PRIMARY")), func(it Astnode) Astnode { return NewNot(it) }),
							o(C(S(St("(?"), L("expr", R("EXPR")), St(")")), S(St("?"), L("expr", R("EXPR")))), func(it Astnode) Astnode { return NewLookahead(it) }),
							i(Named("RANGE", o(S(St("{"), R("_"), L("min", E(R("INT"))), R("_"), St(","), R("_"), L("max", E(R("INT"))), R("_"), St("}"))))),
						))),
						o(Named("PRIMARY", rules(
							o(S(R("WORD"), St("("), R("EXPR"), St(")")), func(it Astnode) Astnode {
								na := it.(*NativeArray)
								if na.Length() != 4 {
									panic("logic")
								}
								return NewRef(NewNativeArray([]Astnode{na.Get(1), na.Get(3)}))
							}),
							o(R("WORD"), func(it Astnode) Astnode { return NewRef(it) }), // TODO really need callback here?
							// o S(St('('), L("inlineLabel",E(S(R('WORD'), St(': ')))), L("expr",R("EXPR")), St(')'), E(S(R('_'), St('->'), R('_'), L("code",R("CODE"))))), ({expr, code}) ->
							//   assert.ok not code?, "code in joeson deprecated"
							//   return expr
							// i "CODE", o S(St("{"), P(S(N(St("}")), C(R("ESC1"), R(".")))), St("}")), (it) -> require('./joescript').parse(it.join '')),
							// Note:this v--- P() here *seems* useless. It's a Pattern(value, join=nil, min=-1, max=-1). However, that kind of pattern has a @capture = @value.capture, which allows it to capture the string.  at least that's my guess atm
							o(S(St("'"), P(S(N(St("'")), C(R("ESC1"), R("."))), nil), St("'")), func(it Astnode) Astnode {
								// this will require heavy testings, best case
								// scenario is captures got into a NativeArray
								// of NativeString, there are very few chances
								// we get it right at first though; better get
								// prepared.
								return NewStr(AttemptToJoinANativeArrayOrPanic(it))
							}),
							o(S(St("/"), P(S(N(St("/")), C(R("ESC2"), R("."))), nil), St("/")), func(it Astnode) Astnode { return NewRegexFromString(AttemptToJoinANativeArrayOrPanic(it)) }),
							o(S(St("["), P(S(N(St("]")), C(R("ESC2"), R("."))), nil), St("]")), func(it Astnode) Astnode { return NewRegexFromString("[" + AttemptToJoinANativeArrayOrPanic(it) + "]") }),
						))),
					))),
				))),
			))),
		))),
	))),
	i(Named("LABEL", C(St("&"), St("@"), R("WORD")))),
	i(Named("WORD", Re("[a-zA-Z\\._][a-zA-Z\\._0-9]*"))),
	i(Named("INT", Re("[0-9]+"))), //, func(it Astnode) Astnode { return NewNativeIntFromNativeString(it.(NativeString)) }),
	i(Named("_PIPE", S(R("_"), St("|")))),
	i(Named("_", P(C(St(" "), St("\n")), nil))),
	i(Named("__", P(C(St(" "), St("\n")), nil, 1))),
	i(Named(".", Re("[\\s\\S]"))),
	i(Named("ESC1", S(St("\\"), R(".")))),
	i(Named("ESC2", S(St("\\"), R("."))), func(chr Astnode) Astnode { return NewNativeString("\\" + chr.(NativeString).Str) }),
	// i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true},
	// 	func(it Astnode, ctx *ParseContext) Astnode {
	// 		// ctx.SkipLog is false
	// 		// ctx.SkipCache is true
	// 		// ctx.Debug is false
	// 		return nil
	// })),
}

// vim: fdm=marker fdl=0
