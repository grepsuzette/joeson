package main

// vim: fdm=marker fdl=0
// run with `go test .` or `make test`
// or individually like `go test . --run "^TestCalculator\$"`

import (
	"fmt"
	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/ast/handcompiled"
	. "grepsuzette/joeson/colors"
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
								o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", func(it Astnode) Astnode {
									return NewPattern(it)
								}),
								o("value:PRIMARY '+' join:(!__ PRIMARY)?", func(it Astnode) Astnode {
									h := it.(NativeMap)
									// maybe TODO check value and join is range imported?
									h.Set("Min", NewNativeInt(1))
									h.Set("Max", NewNativeInt(-1))
									return NewPattern(h)
								}),
								o("value:PRIMARY @:RANGE", func(it Astnode) Astnode {
									// maybe TODO
									return NewPattern(it)
								}),
								o("'!' PRIMARY", func(it Astnode) Astnode { return NewNot(it) }),
								o("'(?' expr:EXPR ')' | '?' expr:EXPR", func(it Astnode) Astnode { return NewLookahead(it) }),
								i(Named("RANGE", "'{' _ min:INT? _ ',' _ max:INT? _ '}'")), // TODO do this work
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
									// TODO
									// assert.ok not code?, "code in joeson deprecated"
									// return expr
									panic("TODO")
									// var caps = it.GetGNode().Captures_.Get()
									// if len(caps) == 0 {
									// 	return NewNativeUndefined()
									// }
									// switch v := caps[0].(type) {
									// case NativeString:
									// 	if len(v.Str) == 0 {
									// 		return NewNativeUndefined()
									// 	} else {
									// 		panic("code in joeson is obsolete")
									// 	}
									// default:
									// 	panic("assert")
									// }
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
		i(Named("INT", "/[0-9]+/"), func(it Astnode) Astnode { return NewNativeIntFrom(it) }),
		i(Named("_PIPE", "_ '|'")),
		i(Named("_", "(' ' | '\n')*")),
		i(Named("__", "(' ' | '\n')+")),
		i(Named(".", "/[\\s\\S]/")),
		i(Named("ESC1", "'\\\\' .")),
		i(Named("ESC2", "'\\\\' ."), func(chr Astnode) Astnode { return NewNativeString("\\" + chr.(NativeString).Str) }),
	}
} // }}}1

// this is the bootstrapped grammar, using joeson to define itself,
// similar to joeson_test.coffee
func TestRaw(t *testing.T) {
	raw := line.NewGrammarFromLines(
		"bootstrapped grammar",
		RAW_GRAMMAR(),
		NewJoeson(),
	)
	if !raw.IsReady() {
		t.Fail()
	}
}

func TestPattern(t *testing.T) {
	var f = func(patt *Pattern, tcase string, expectedMin int, expectedMax int, expectedContent string) {
		if patt.Value == nil || patt.Value.(Str).Str != "foo" {
			t.Error(tcase + " patt.Value expected foo")
		}
		if int(patt.Min) != expectedMin {
			t.Errorf(tcase+" patt.Min expected %d, got %d", expectedMin, patt.Min)
		}
		if int(patt.Max) != expectedMax {
			t.Errorf(tcase+" patt.Max expected %d, got %d", expectedMax, patt.Max)
		}
		if patt.ContentString() != expectedContent {
			t.Errorf(tcase+" patt.ContentString() expected %s, got %s", expectedContent, patt.ContentString())
		}
	}
	tcase := "TestPattern case#1"
	patt := NewPattern(NewNativeMap(map[string]Astnode{
		"value": NewStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt, tcase, 2, -1, Green("'foo'")+Cyan("*")+Cyan("{2,}"))
	tcase = "TestPattern case#2"
	patt2 := NewPattern(NewNativeMap(map[string]Astnode{
		"value": NewStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeInt(4),
	}))
	f(patt2, tcase, 2, 4, Green("'foo'")+Cyan("*")+Cyan("{2,4}"))
	tcase = "TestPattern case#3"
	patt3 := NewPattern(NewNativeMap(map[string]Astnode{
		"value": NewStr("foo"),
		"min":   NewNativeInt(-1),
		"max":   NewNativeInt(-1),
	}))
	f(patt3, tcase, -1, -1, Green("'foo'")+Cyan("*")+Cyan(""))
	tcase = "TestPattern case#4"
	patt4 := NewPattern(NewNativeMap(map[string]Astnode{
		"value": NewStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt4, tcase, 2, -1, Green("'foo'")+Cyan("*")+Cyan("{2,}"))
	tcase = "TestPattern case#5(non-nil join)"
	patt5 := NewPattern(NewNativeMap(map[string]Astnode{
		"value": NewStr("foo"),
		"join":  NewStr("bar"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt5, tcase, 2, -1, Green("'foo'")+Cyan("*")+Green("'bar'")+Cyan("{2,}"))
	tcase = "TestPattern case#5(non-nil join)"
	patt6 := NewPattern(NewNativeMap(map[string]Astnode{
		"value": NewStr("foo"),
		"join":  NewStr("bar"),
		"min":   NewNativeInt(-1),
		"max":   NewNativeUndefined(),
	}))
	f(patt6, tcase, -1, -1, Green("'foo'")+Cyan("*")+Green("'bar'")+Cyan(""))

	// in handcompiled, test the function P(..) *Pattern
	tcase = "TestPattern:TestHandcompiledFuncP case#1"
	patt = P(NewStr("foo"), NewNativeUndefined(), 2, -1)
	f(patt, tcase, 2, -1, Green("'foo'")+Cyan("*")+Cyan("{2,}"))
}

func TestHandcompiled(t *testing.T) {
	gm := NewJoeson()
	if gm.GetGNode().Name != JOESON_GRAMMAR_NAME {
		t.Fail()
	}
	if gm.CountRules() != gm.NumRules || gm.CountRules() != JoesonNbRules {
		t.Errorf("Expected %d rules, got %d\n", JoesonNbRules, gm.CountRules())
	}
	if !gm.IsReady() {
		t.Fail()
	}
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

// this short grammar was useful for debugging
func TestDebugLabel(t *testing.T) {
	debuglabel := line.NewGrammarFromLines(
		"gm_DebugLabel",
		[]line.Line{
			o(Named("In", "l:Br")),
			i(Named("Br", "'Toy' | 'BZ'")),
		},
		NewJoeson(),
	)
	debuglabel.PrintRules()
	if x, error := debuglabel.ParseString("Toy"); error == nil {
		if nm, isNativeMap := x.(NativeMap); !isNativeMap {
			t.Errorf("expected NativeMap, got %T. ContentString: %s\n", x, x.ContentString())
		} else {
			// in two operations...
			if label, exists := nm.GetExists("l"); !exists {
				t.Fail()
			} else if label.(NativeString).Str != "Toy" {
				t.Fail()
			}
			// ...or in 1 operation
			if label, exists := nm.GetStringExists("l"); !exists || label != "Toy" {
				t.Fail()
			}
		}
	} else {
		t.Error(error)
	}
}
