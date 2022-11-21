package main

// vim: fdm=marker fdl=0
// run with `go test .`

import (
	"fmt"

	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/ast/handcompiled"

	. "grepsuzette/joeson/core"

	line "grepsuzette/joeson/line"
	"testing"
)

// func o(a ...any) line.OLine               { return line.O(a...) }
// func i(a ...any) line.ILine               { return line.I(a...) }
// func Rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }
// func Named(name string, lineStringOrAstnode any) line.NamedRule {
// 	return line.Named(name, lineStringOrAstnode)
// }

var RAW_GRAMMAR = []line.Line{ // {{{1
	o(Named("EXPR", Rules(
		o("CHOICE _"),
		o(Named("CHOICE", Rules(
			o("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Astnode) Astnode { return NewChoice(it) }),
			o(Named("SEQUENCE", Rules(
				o("UNIT{2,}", func(it Astnode) Astnode { return NewSequence(it) }),
				o(Named("UNIT", Rules(
					o("_ LABELED"),
					o(Named("LABELED", Rules(
						o("(label:LABEL ':')? &:(DECORATED|PRIMARY)"),
						o(Named("DECORATED", Rules(
							o("PRIMARY '?'", func(it Astnode) Astnode { return NewExistential(it) }),
							o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", func(it Astnode) Astnode { return NewPattern(it) }),
							o("value:PRIMARY '+' join:(!__ PRIMARY)?", func(it Astnode) Astnode {
								h := it.(NativeMap)
								h.Set("Min", NewNativeInt(1))
								h.Set("Max", NewNativeInt(-1))
								return NewPattern(h)
							}),
							o("value:PRIMARY @:RANGE", func(it Astnode) Astnode { return NewPattern(it) }), // note: the @ label will "source" and "import" the labels from RANGE node into `it`
							o("'!' PRIMARY", func(it Astnode) Astnode { return NewNot(it) }),
							o("'(?' expr:EXPR ')' | '?' expr:EXPR", func(it Astnode) Astnode { return NewLookahead(it) }),
							i(Named("RANGE", "'{' _ min:INT? _ ',' _ max:INT? _ '}'")),
						))),
						o(Named("PRIMARY", Rules(
							o("WORD '(' EXPR ')'", func(it Astnode) Astnode {
								na := it.(*NativeArray)
								if na.Length() != 4 {
									panic("logic")
								}
								return NewRef(NewNativeArray([]Astnode{na.Get(1), na.Get(3)}))
							}),
							o("WORD", func(it Astnode) Astnode { return NewRef(it) }), // TODO really need callback here?
							// TODO inlineLabel, but code in joeson deprecated
							// TODO i "CODE "
							// Note:this v--- P() here *seems* useless. It's a Pattern(value, join=nil, min=-1, max=-1). However, that kind of pattern has a @capture = @value.capture, which allows it to capture the string.  at least that's my guess atm
							o("'\\'' (!'\\'') (ESC1 | .))* '\\''", func(it Astnode) Astnode {
								// this will require heavy testings, best case
								// scenario is captures got into a NativeArray
								// of NativeString, there are very few chances
								// we get it right at first though; better get
								// prepared.
								return NewStr(AttemptToJoinANativeArrayOrPanic(it))
							}),
							o("'/' (!'/') (ESC2 | .))* '/'", func(it Astnode) Astnode { return NewRegexFromString(AttemptToJoinANativeArrayOrPanic(it)) }),
							o("'[' (!']') (ESC2 | .))* ']'", func(it Astnode) Astnode { return NewRegexFromString("[" + AttemptToJoinANativeArrayOrPanic(it) + "]") }),
						))),
					))),
				))),
			))),
		))),
	))),
	i(Named("LABEL", "'&' | '@' | WORD")),
	i(Named("WORD", "/[a-zA-Z\\._][a-zA-Z\\._0-9]*/")),
	i(Named("INT", "/[0-9]+/"), func(it Astnode) Astnode { return NewNativeIntFromNativeString(it.(NativeString)) }),
	i(Named("_PIPE", "_ '|'")),
	i(Named("_", "(' ' | '\n')*")),
	i(Named("__", "(' ' | '\n')+")),
	i(Named(".", "/[\\s\\S]/")),
	i(Named("ESC1", "'\\\\' .")),
	i(Named("ESC2", "'\\\\' ."), func(chr Astnode) Astnode { return NewNativeString("\\" + chr.(NativeString).Str) }),
} // }}}1

// func testItCompilesToBeLikeHandcompiled() {
// 	var PARSED_GRAMMAR *Grammar = line.NewGrammarFromLines("joeson raw grammar", RAW_GRAMMAR)
// 	fmt.Println(PARSED_GRAMMAR.ContentString())
// }

func TestHandcompiledGrammar(t *testing.T) {
	gm := line.NewGrammarFromLines("joeson from handcompiled", JOESON_GRAMMAR_RULES, NewEmptyGrammarNamed("empty grammar"))
	if gm.GetGNode().Name != "joeson from handcompiled" {
		t.Fail()
	}
	// 33 right now, because 2/35 rules (about embedded code) were deprecated in joeson.coffee
	if gm.CountRules() != gm.NumRules || gm.CountRules() != JoesonNbRules {
		t.Errorf("Expected %d rules, got %d\n", JoesonNbRules, gm.CountRules())
	}
	if !gm.IsReady() {
		t.Fail()
	}
}

func TestAab(t *testing.T) {
	var JOESON *Grammar = line.NewGrammarFromLines("joeson from handcompiled", JOESON_GRAMMAR_RULES, NewEmptyGrammarNamed("empty grammar"))
	var AAB = []line.Line{
		o(Named("EXPRESSION", Rules(
			o("A EXPRESSION|B"),
			i(Named("A", "'A' | 'a'")),
			i(Named("B", "'B' | 'b'")),
		))),
	}
	var gm = line.NewGrammarFromLines("aab", AAB, JOESON)
	if gm.NumRules != 4 {
		t.Fail()
	}
	fmt.Println(gm.ContentString())
	t.Log("Done")
}
