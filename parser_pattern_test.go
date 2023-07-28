package joeson

import (
	"strings"
	"testing"
)

// this tests the parsing of the pattern parser
func TestPattern(t *testing.T) {
	{
		rule := "a*b"
		gm := GrammarFromLines([]Line{
			o(named("Input", rule)),
			i(named("a", "'a'")),
			i(named("b", "'b'")),
		}, rule)
		mustParse(t, gm, "", "[]")
		mustParse(t, gm, "a", "[a]")
		mustParse(t, gm, "abababababa", "[a,a,a,a,a,a]")
		mustFail(t, gm, "ab")
		mustFail(t, gm, "aaa")
	}
	{
		rule := "a*b{1,2}"
		// the {} repetition applies to the whole sequence
		gm := GrammarFromLines([]Line{
			o(named("Input", rule)),
			i(named("a", "'a'")),
			i(named("b", "'b'")),
		}, rule)
		mustParse(t, gm, "a", "[a]")
		mustFail(t, gm, "ab")
		mustParse(t, gm, "aba", "[a,a]")
		mustFail(t, gm, "abab")
		mustFail(t, gm, "ababa")
		mustFail(t, gm, "aabaa")
	}
	{
		rule := "a*(b{,2})"
		// the join value is b{,2} here,
		// this is why we use brackets
		gm := GrammarFromLines([]Line{
			o(named("Input", rule)),
			i(named("a", "'a'")),
			i(named("b", "'b'")),
		}, rule)
		mustParse(t, gm, "abbabababbabba")
	}
}

// this does not test the behavior but only the parameters
func TestNewPattern(t *testing.T) {
	f := func(patt *pattern, tcase string, expectedMin int, expectedMax int, expectedContent string) {
		if patt.value == nil || patt.value.(str).Str != "foo" {
			t.Error(tcase + " patt.Value expected foo")
		}
		if int(patt.min) != expectedMin {
			t.Errorf(tcase+" patt.Min expected %d, got %d", expectedMin, patt.min)
		}
		if int(patt.max) != expectedMax {
			t.Errorf(tcase+" patt.Max expected %d, got %d", expectedMax, patt.max)
		}
		if patt.String() != expectedContent {
			t.Errorf(tcase+" patt.String() expected %s, got %s", expectedContent, patt.String())
		}
	}
	tcase := "TestPattern case#1"
	patt := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt, tcase, 2, -1, Green("'foo'")+Cyan("*")+Cyan("{2,}"))
	tcase = "TestPattern case#2"
	patt2 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeInt(4),
	}))
	f(patt2, tcase, 2, 4, Green("'foo'")+Cyan("*")+Cyan("{2,4}"))
	tcase = "TestPattern case#3"
	patt3 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(-1),
		"max":   NewNativeInt(-1),
	}))
	f(patt3, tcase, -1, -1, Green("'foo'")+Cyan("*")+Cyan(""))
	tcase = "TestPattern case#4"
	patt4 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt4, tcase, 2, -1, Green("'foo'")+Cyan("*")+Cyan("{2,}"))
	tcase = "TestPattern case#5(non-nil join)"
	patt5 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"join":  newStr("bar"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt5, tcase, 2, -1, Green("'foo'")+Cyan("*")+Green("'bar'")+Cyan("{2,}"))
	tcase = "TestPattern case#5(non-nil join)"
	patt6 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"join":  newStr("bar"),
		"min":   NewNativeInt(-1),
		"max":   NewNativeUndefined(),
	}))
	f(patt6, tcase, -1, -1, Green("'foo'")+Cyan("*")+Green("'bar'")+Cyan(""))
}

func mustParse(t *testing.T, gm *Grammar, s string, sMustContain ...string) {
	t.Helper()
	it := gm.ParseString(s)
	if IsParseError(it) {
		t.Errorf("grammar \"%s\": \"%s\" failed to parse: \"%s\"",
			gm.name, s, it.String(),
		)
	} else {
		if len(sMustContain) > 0 && !strings.Contains(it.String(), sMustContain[0]) {
			t.Errorf("grammar \"%s\": \"%s\" parsed as \"%s\" but fails to contain \"%s\"",
				gm.name, s, it.String(), sMustContain[0],
			)
		}
	}
}

func mustFail(t *testing.T, gm *Grammar, s string, sMustContain ...string) {
	t.Helper()
	it := gm.ParseString(s)
	if IsParseError(it) {
		if len(sMustContain) > 0 && !strings.Contains(it.String(), sMustContain[0]) {
			t.Errorf("grammar \"%s\": \"%s\" failed to parse but it does not contain \"%s\", instead it gave \"%s\"",
				gm.name, s, sMustContain[0], it.String(),
			)
		}
	} else {
		t.Errorf("grammar \"%s\": \"%s\" should have fail to parse. But it gave \"%s\"",
			gm.name, s, it.String(),
		)
	}
}
