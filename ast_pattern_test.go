package joeson

import (
	"testing"
)

func TestPattern(t *testing.T) {
	var f = func(patt *pattern, tcase string, expectedMin int, expectedMax int, expectedContent string) {
		if patt.Value == nil || patt.Value.(str).Str != "foo" {
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
	patt := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt, tcase, 2, -1, green("'foo'")+cyan("*")+cyan("{2,}"))
	tcase = "TestPattern case#2"
	patt2 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeInt(4),
	}))
	f(patt2, tcase, 2, 4, green("'foo'")+cyan("*")+cyan("{2,4}"))
	tcase = "TestPattern case#3"
	patt3 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(-1),
		"max":   NewNativeInt(-1),
	}))
	f(patt3, tcase, -1, -1, green("'foo'")+cyan("*")+cyan(""))
	tcase = "TestPattern case#4"
	patt4 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt4, tcase, 2, -1, green("'foo'")+cyan("*")+cyan("{2,}"))
	tcase = "TestPattern case#5(non-nil join)"
	patt5 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"join":  newStr("bar"),
		"min":   NewNativeInt(2),
		"max":   NewNativeUndefined(),
	}))
	f(patt5, tcase, 2, -1, green("'foo'")+cyan("*")+green("'bar'")+cyan("{2,}"))
	tcase = "TestPattern case#5(non-nil join)"
	patt6 := newPattern(NewNativeMap(map[string]Ast{
		"value": newStr("foo"),
		"join":  newStr("bar"),
		"min":   NewNativeInt(-1),
		"max":   NewNativeUndefined(),
	}))
	f(patt6, tcase, -1, -1, green("'foo'")+cyan("*")+green("'bar'")+cyan(""))
}
