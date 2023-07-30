package joeson

import (
	"fmt"
	"regexp"
	"strings"
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
	cs := NewRuneStream(`line 0
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

func TestRegexp(t *testing.T) {
	cs := NewRuneStream("A EXPRESSION|B")
	cs.pos = 2
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

func TestPosToLine(t *testing.T) {
	source := `
	// RuneStream is a very simple code holder, cursor, matcher.
	type RuneStream struct {
		text       string
		pos        int // "Hello, 世界, X" <- Pos of o is 4, Pos of 界 is 10
		lineStarts []int
	}
	`
	code := NewRuneStream(source)
	index := strings.Index(code.text, "text")
	line := code.PosToLine(index)
	if line != 3 {
		t.Errorf(fmt.Sprintf("Invalid line found %d, expected 3", line))
	}
	col := code.PosToCol(index)
	if col != 2 {
		t.Errorf(fmt.Sprintf("Invalid col found %d, expected 2", col))
	}
}
