package joeson

import (
	"testing"
)

func TestPattern(t *testing.T) {
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
