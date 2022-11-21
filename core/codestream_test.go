package core

import (
	"fmt"
	"grepsuzette/joeson/helpers"
	"regexp"
	"strconv"
	"testing"
)

func eq_int(t *testing.T, a int, b int) {
	if a != b {
		// t.Errorf("assert %d == %d", a, b)
		panic("assertion failed: " + strconv.Itoa(a) + " == " + strconv.Itoa(b))
	}
}
func eq_str(t *testing.T, a string, b string) {
	if a != b {
		// t.Errorf("assertion failed: \"%s\" == \"%s\"", a, b)
		panic(fmt.Sprintf("assertion failed: \"%s\" == \"%s\"", a, b))
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

func TestMinMax(t *testing.T) {
	eq_int(t, helpers.Min(-1999, 0), -1999)
	eq_int(t, helpers.Min(0, -10), -10)
	eq_int(t, helpers.Min(0, 10), 0)
	eq_int(t, helpers.Max(-1999, 0), 0)
	eq_int(t, helpers.Max(0, -10), 0)
	eq_int(t, helpers.Max(0, 10), 10)
}

func TestSliceString(t *testing.T) {
	eq_str(t, helpers.SliceString("FOo", 0, 0), "")
	eq_str(t, helpers.SliceString("FOo", 0, 1), "F")
	eq_str(t, helpers.SliceString("FOo", 1, 2), "O")
	eq_str(t, helpers.SliceString("FOo", 2, 3), "o")
	eq_str(t, helpers.SliceString("FOo", 2, 4), "o")
	eq_str(t, helpers.SliceString("FOo", 3, 4), "")
	eq_str(t, helpers.SliceString("FOo", 6, 10), "")
	eq_str(t, helpers.SliceString("FOo", 6, -2), "")
	eq_str(t, helpers.SliceString("FOo", -1, -2), "")
	eq_str(t, helpers.SliceString("FOo", -1, 2), "")
	eq_str(t, helpers.SliceString("FOo", -1, 4), "")
	eq_str(t, helpers.SliceString("FOo", 0, 3), "FOo")
	eq_str(t, helpers.SliceString("FOo", 1, 3), "Oo")
	eq_str(t, helpers.SliceString("FOo", 2, 3), "o")
}

func TestPad(t *testing.T) {
	eq_str(t, helpers.PadLeft("foo", 5), "foo  ")
	eq_str(t, helpers.PadLeft("foo", 4), "foo ")
	eq_str(t, helpers.PadLeft("foo", 3), "foo")
	eq_str(t, helpers.PadLeft("foo", 2), "foo")
	eq_str(t, helpers.PadLeft("foo", 1), "foo")
	eq_str(t, helpers.PadLeft("foo", 0), "foo")
	eq_str(t, helpers.PadLeft("foo", -1), "foo")
	eq_str(t, helpers.PadLeft("foo", -10), "foo")
	eq_str(t, helpers.PadRight("foo", 5), "  foo")
	eq_str(t, helpers.PadRight("foo", 4), " foo")
	eq_str(t, helpers.PadRight("foo", 3), "foo")
	eq_str(t, helpers.PadRight("foo", 2), "foo")
	eq_str(t, helpers.PadRight("foo", 1), "foo")
	eq_str(t, helpers.PadRight("foo", 0), "foo")
	eq_str(t, helpers.PadRight("foo", -1), "foo")
	eq_str(t, helpers.PadRight("foo", -10), "foo")
}

func TestIndent(t *testing.T) {
	eq_str(t, helpers.Indent(-1), "")
	eq_str(t, helpers.Indent(0), "")
	eq_str(t, helpers.Indent(1), "  ")
	eq_str(t, helpers.Indent(2), "    ")
	eq_str(t, helpers.Indent(3), "      ")
}

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
