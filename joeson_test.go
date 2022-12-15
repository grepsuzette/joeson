package main

// vim: fdm=marker fdl=0
// run with `go test .` or `make test`
// or individually like `go test . --run=TestAab`

import (
	"fmt"
	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/ast/handcompiled"
	// . "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	// "grepsuzette/joeson/helpers"
	line "grepsuzette/joeson/line"
	"testing"
	"time"
)

func RAW_GRAMMAR() line.Lines {
	return []line.Line{ // {{{1
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
								o("'(' inlineLabel:(WORD ': ')? expr:EXPR ')' ( _ '->' _ code:CODE )?", func(it Astnode) Astnode {
									h := it.(NativeMap)
									if _, isNativeUndefined := h.Get("code").(NativeUndefined); !isNativeUndefined {
										panic("code in joeson deprecated")
									}
									return h.Get("expr")
								}),
								i(Named("CODE", "'{' (!'}' (ESC1 | .))* '}'"), func(it Astnode) Astnode {
									// deprecate code in joeson
									var caps = it.Captures()
									if len(caps) == 0 {
										return NewNativeUndefined()
									}
									switch v := caps[0].(type) {
									case NativeString:
										if len(v.Str) == 0 {
											return NewNativeUndefined()
										} else {
											panic("code in joeson is obsolete")
										}
									default:
										panic("assert")
									}
								}),
								o("'\\'' (!'\\'' (ESC1 | .))* '\\''", func(it Astnode) Astnode {
									// this will require heavy testings, best case
									// scenario is captures got into a NativeArray
									// of NativeString, there are very few chances
									// we get it right at first though; better get
									// prepared.
									return NewStr(AttemptToJoinANativeArrayOrPanic(it))
								}),
								o("'/' (!'/' (ESC2 | .))* '/'", func(it Astnode) Astnode { return NewRegexFromString(AttemptToJoinANativeArrayOrPanic(it)) }),
								o("'[' (!']' (ESC2 | .))* ']'", func(it Astnode) Astnode { return NewRegexFromString("[" + AttemptToJoinANativeArrayOrPanic(it) + "]") }),
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
	}
} // }}}1

func TestRaw(t *testing.T) {
	raw := line.NewGrammarFromLines(
		"bootstrapped grammar",
		RAW_GRAMMAR(),
		NewJoeson(),
	)
	fmt.Println(raw.ContentString())
}

func TestHandcompiled(t *testing.T) {
	gm := NewJoeson()
	if gm.GetGNode().Name != "joeson from handcompiled" {
		t.Fail()
	}
	if gm.CountRules() != gm.NumRules || gm.CountRules() != JoesonNbRules {
		t.Errorf("Expected %d rules, got %d\n", JoesonNbRules, gm.CountRules())
	}
	if !gm.IsReady() {
		t.Fail()
	}
}

func TestAab(t *testing.T) {
	joeson := NewJoeson()
	AAB := []line.Line{
		o(Named("EXPRESSION", Rules(
			o("A EXPRESSION|B"),
			i(Named("A", "'A' | 'a'")),
			i(Named("B", "'B'|'b'")),
		))),
	}
	aab := line.NewGrammarFromLines("aab", AAB, joeson)
	if aab.NumRules != 4 {
		t.Fail()
	}
	fmt.Println(aab.ContentString())
}

func Test100Times(t *testing.T) {
	// this test comes directly from joeson_test.coffee
	start := time.Now()
	iter := 100
	for i := 0; i < iter; i++ {
		// testGrammar(line.NewALine(RAW_GRAMMAR()), 0, "")
		fmt.Println(line.NewALine(RAW_GRAMMAR()).StringIndent(0))
		fmt.Println("-------------")
	}
	fmt.Printf("Duration for %d iterations: %d ms\n", iter, time.Now().Sub(start).Milliseconds())
}

func TestCalculator(t *testing.T) {
	joeson := NewJoeson()
	CALC := []line.Line{
		o(Named("Input", "expr:Expression")),
		i(Named("Expression", "first:Term rest:( _ AddOp _ Term )*")),
		i(Named("Term", "first:Factor rest:( _ MulOp _ Factor  )*")),
		i(Named("Factor", "'(' expr:Expression _ ')' | integer:Integer")),
		i(Named("AddOp", "'+' | '-'")),
		i(Named("MulOp", "'*' | '/'")),
		// i(Named("Integer", "'-'? [0-9]{1,}")),
		i(Named("Integer", "[0-9]{1,}")),
		i(Named("_", "[ \n\t\r]*")),
		// i(Named("EOF", "!.")),
	}
	calc := line.NewGrammarFromLines("calc", CALC, joeson)
	if !calc.IsReady() {
		t.Fail()
	} else {
		fmt.Println(calc.ContentString())
		// x := calc.ParseString("241 + (513 ) * -24 + ((1934 - 192 *2)/7) +1", ParseOptions{})
		x := calc.ParseString("241 + (513 ) -24", ParseOptions{})
		fmt.Println(x.ContentString())
		panic("ok")
	}
}
