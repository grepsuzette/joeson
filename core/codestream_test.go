package core

import (
	"grepsuzette/joeson/helpers"
	"strconv"
	"testing"
)

func eq_int(t *testing.T, a int, b int) {
	if a != b {
		// t.Errorf("assert %d == %d", a, b)
		panic("assert " + strconv.Itoa(a) + " == " + strconv.Itoa(b))
	}
}
func eq_str(t *testing.T, a string, b string) {
	if a != b {
		t.Errorf("assert %s == %s", a, b)
	}
}

func TestState(t *testing.T) {
	cs := newCodeStream(`line 0
line 1
line 2 -- rest of line 2
line 3
line 4`)
	eq_int(t, cs.line(), 0)
	eq_int(t, cs.col(), 0)
	line := cs.getUntil("\n")
	eq_str(t, line, "line 0\n")

	eq_int(t, cs.line(), 1) // FAIL  is 0
	eq_int(t, cs.col(), 0)
	line = cs.getUntil("\n")
	eq_str(t, line, "line 1\n")
	eq_int(t, cs.line(), 2)
	eq_int(t, cs.col(), 0)

	part := cs.getUntil(" -- ")
	eq_str(t, part, "line 2 -- ")
	eq_int(t, cs.line(), 2)
	eq_int(t, cs.col(), 10)
	eq_str(t, cs.peek(Peek{afterChars: helpers.NewNullInt(4)}), "rest")
	eq_str(t, cs.peek(Peek{beforeLines: helpers.NewNullInt(0), afterLines: helpers.NewNullInt(0)}), "line 2 -- rest of line 2")
	eq_str(t, cs.peek(Peek{afterLines: helpers.NewNullInt(0)}), "rest of line 2")
	eq_str(t, cs.peek(Peek{afterLines: helpers.NewNullInt(1)}), "rest of line 2\nline 3")
}
