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

func TestRuneStreamState(t *testing.T) {
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
	eq_str(t, cs.PeekRunes(4), "rest")
	// TODO uncomment and fix these 2, behavior changed as of August 07, 2023
	// eq_str(t, cs.PeekLines(0), "rest of line 2")
	// eq_str(t, cs.PeekLines(1), "rest of line 2\nline 3")
}

func TestRuneStreamRegexp(t *testing.T) {
	cs := NewRuneStream("A EXPRESSION|B")
	cs.SetPos(2)
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

func TestRuneStreamPosToLine(t *testing.T) {
	source := `
	// RuneStream is a very simple code holder, cursor, matcher.
	type RuneStream struct {
		text       string
		pos        int // "Hello, 世界, X" <- Pos of o is 4, Pos of 界 is 10
		lineStarts []int
	}
	`
	code := NewRuneStream(source)
	index := strings.Index(code.Code(), "text")
	line := code.PosToLine(index)
	if line != 3 {
		t.Errorf(fmt.Sprintf("Invalid line found %d, expected 3", line))
	}
	col := code.PosToCol(index)
	if col != 2 {
		t.Errorf(fmt.Sprintf("Invalid col found %d, expected 2", col))
	}
}

func TestRuneStreamPeekLines(t *testing.T) {
	s := "rose are blue\nblue are violet\nviolet are pi/2"
	code := NewRuneStream(s)
	index := strings.Index(code.Code(), "blue are violet")
	code.SetPos(index)
	{
		assertPanics(t, func() { code.SetPos(-1) })
		assertPanics(t, func() { code.SetPos(9999999) })
	}
	{
		peeked := code.PeekLines(-1, 1)
		if peeked != s {
			t.Errorf("expected %q, got %q\n", s, peeked)
		}
	}
	{
		peeked := code.PeekLines(0, 1)
		expected := "blue are violet\nviolet are pi/2"
		if peeked != expected {
			t.Errorf("expected %q, got %q\n", expected, peeked)
		}
	}
	{
		peeked := code.PeekLines(-1, 0)
		expected := "rose are blue\nblue are violet"
		if peeked != expected {
			t.Errorf("expected %q, got %q\n", expected, peeked)
		}
	}
	{
		peeked := code.PeekLines(-99, 0)
		expected := "rose are blue\nblue are violet"
		if peeked != expected {
			t.Errorf("expected %q, got %q\n", expected, peeked)
		}
	}
	{
		peeked := code.PeekLines(99, 9, 1, 0)
		expected := "blue are violet\nviolet are pi/2"
		if peeked != expected {
			t.Errorf("expected %q, got %q\n", expected, peeked)
		}
	}
}

func TestRuneStreamUnicode(t *testing.T) {
	{
		abc := NewRuneStream("abc")
		s := abc.GetUntil("a")
		if s != "a" {
			t.Errorf("should have obtained \"a\", not %q", s)
		}
	}
	{
		abc := NewRuneStream("αβγ")
		s := abc.GetUntil("α")
		if s != "α" {
			t.Errorf("should have obtained \"α\", not %q", s)
		}
	}
	{
		abc := NewRuneStream("αβγ")
		abc.SetPos(0)
		ok, m := abc.MatchString("α")
		if !ok {
			t.Errorf("should have matched")
		} else if m != "α" {
			t.Errorf("should have matched α")
		}
	}
	{
		abc := NewRuneStream("abc")
		abc.SetPos(1)
		ok, m := abc.MatchString("b")
		if !ok {
			t.Errorf("should have matched")
		} else if m != "b" {
			t.Errorf("should have obtained \"b\", not %q", m)
		}
	}
	{
		abc := NewRuneStream("αβγ")
		abc.SetPos(2) // note: not meant to be used like that!
		ok, m := abc.MatchString("β")
		if !ok {
			t.Errorf("should have matched")
		} else if m != "β" {
			t.Errorf("should have matched β")
		}
	}
	{
		abc := NewRuneStream("αβγ")
		abc.SetPos(4) // note: not meant to be used like that!
		fmt.Printf("% x\n", abc.Code())
		ok, m := abc.MatchString("γ")
		if !ok {
			t.Errorf("should have matched")
		} else if m != "γ" {
			t.Errorf("should have matched γ")
		}
	}
	{
		abc := NewRuneStream("αβγ")
		abc.SetPos(0) // note: not meant to be used like that!
		fmt.Printf("% x\n", abc.Code())
		ok, c := abc.MatchRune(func(rune rune) bool { return 'γ' == rune })
		if ok {
			t.Errorf("should NOT have matched")
		}
		ok, c = abc.MatchRune(func(rune rune) bool { return 'α' == rune })
		if !ok {
			t.Errorf("should have matched")
		} else if c != 'α' {
			t.Errorf("should have matched α, got %q", c)
		}
		if abc.Pos() != 2 {
			t.Errorf("Pos should have been updated")
		}
		ok, c = abc.MatchRune(func(rune rune) bool { return 'γ' == rune })
		if ok {
			t.Errorf("should NOT have matched")
		}
		ok, c = abc.MatchRune(func(rune rune) bool { return 'β' == rune })
		if !ok {
			t.Errorf("should have matched")
		} else if c != 'β' {
			t.Errorf("should have matched β")
		}
		if abc.Pos() != 4 {
			t.Errorf("Pos should have been updated")
		}
		ok, c = abc.MatchRune(func(rune rune) bool { return 'γ' == rune })
		if !ok {
			t.Errorf("should have matched")
		} else if c != 'γ' {
			t.Errorf("should have matched γ")
		}
		if abc.Pos() != 6 {
			t.Errorf("Pos should have been updated")
		}
		// We are now at the end of the text...
		// it MUST not match
		if ok, _ = abc.MatchRune(func(rune rune) bool { return true }); ok {
			t.Errorf("MatchRune must never match at EOF")
		}
	}
}
