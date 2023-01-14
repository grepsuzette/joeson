package core

import (
	"regexp"
	"testing"
)

func eq_int(t *testing.T, a int, b int) {
	t.Helper()
	if a != b {
		t.Errorf("assert %d == %d", a, b)
	}
}
func eq_str(t *testing.T, a string, b string) {
	t.Helper()
	if a != b {
		t.Errorf("assertion failed: \"%s\" == \"%s\"", a, b)
	}
}

func TestState(t *testing.T) {
	cs := NewCodeStream(`line 0
line 1
line 2 -- rest of line 2
line 3
line 4`)
	eq_int(t, cs.Line(), 0)
	eq_int(t, cs.Col(), 0)
	line := cs.GetUntil("\n")
	eq_str(t, line, "line 0\n")

	eq_int(t, cs.Line(), 1)
	eq_int(t, cs.Col(), 0)
	line = cs.GetUntil("\n")
	eq_str(t, line, "line 1\n")
	eq_int(t, cs.Line(), 2)
	eq_int(t, cs.Col(), 0)

	part := cs.GetUntil(" -- ")
	eq_str(t, part, "line 2 -- ")
	eq_int(t, cs.Line(), 2)
	eq_int(t, cs.Col(), 10)
	eq_str(t, cs.Peek(NewPeek().AfterChars(4)), "rest")
	eq_str(t, cs.Peek(NewPeek().BeforeLines(0).AfterLines(0)), "line 2 -- rest of line 2")
	eq_str(t, cs.Peek(NewPeek().AfterLines(0)), "rest of line 2")
	eq_str(t, cs.Peek(NewPeek().AfterLines(1)), "rest of line 2\nline 3")
}

// func TestDumb(t *testing.T) {
// cs := NewCodeStream("FOO")
// eq_str(t, cs.Peek(NewPeek().AfterChars(20)), "")
// }

func TestRegexp(t *testing.T) {
	cs := NewCodeStream("A EXPRESSION|B")
	cs.Pos = 2
	if re, err := regexp.CompilePOSIX("([a-zA-Z\\._][a-zA-Z\\._0-9]*)"); err != nil {
		t.Fatalf("regexp failed to compile")
	} else {
		didMatch, m := cs.MatchRegexp(*re)
		if !didMatch {
			t.Fatalf("expected regexp to match")
		} else if m != "EXPRESSION" {
			t.Fail()
		}
	}
}
