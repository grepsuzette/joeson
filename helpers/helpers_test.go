package helpers

import "testing"

func eq_str(t *testing.T, a string, b string) {
	t.Helper()
	if a != b {
		t.Errorf("assert failed %s == %s", a, b)
	}
}

func eq_int(t *testing.T, a int, b int) {
	t.Helper()
	if a != b {
		t.Errorf("assert %d == %d", a, b)
	}
}

func TestSliceString(t *testing.T) {
	for i := -5; i < 10; i++ {
		eq_str(t, SliceString("june", i, 0), "")
	}
	eq_str(t, SliceString("june", -2, 1), "")
	eq_str(t, SliceString("june", -1, 1), "")
	eq_str(t, SliceString("june", 0, 1), "j")
	eq_str(t, SliceString("june", 1, 1), "")
	eq_str(t, SliceString("june", 2, 1), "")
	eq_str(t, SliceString("june", 3, 1), "")
	eq_str(t, SliceString("june", -2, 2), "")
	eq_str(t, SliceString("june", -1, 2), "")
	eq_str(t, SliceString("june", 0, 2), "ju")
	eq_str(t, SliceString("june", 1, 2), "u")
	eq_str(t, SliceString("june", 2, 2), "")
	eq_str(t, SliceString("june", 3, 2), "")
	eq_str(t, SliceString("june", -6, 3), "jun")
	eq_str(t, SliceString("june", -5, 3), "jun")
	eq_str(t, SliceString("june", -4, 3), "jun")
	eq_str(t, SliceString("june", -3, 3), "un")
	eq_str(t, SliceString("june", -2, 3), "n")
	eq_str(t, SliceString("june", -1, 3), "")
	eq_str(t, SliceString("june", 0, 3), "jun")
	eq_str(t, SliceString("june", 1, 3), "un")
	eq_str(t, SliceString("june", 2, 3), "n")
	eq_str(t, SliceString("june", 3, 3), "")
	eq_str(t, SliceString("june", -5, 4), "june")
	eq_str(t, SliceString("june", -5, 3), "jun")
	eq_str(t, SliceString("june", -5, 2), "ju")
	eq_str(t, SliceString("june", -5, 1), "j")
	eq_str(t, SliceString("june", -5, 0), "")
	eq_str(t, SliceString("june", -5, -1), "jun")
	eq_str(t, SliceString("june", -5, -2), "ju")
	eq_str(t, SliceString("june", -5, -3), "j")
	eq_str(t, SliceString("june", -5, -4), "")
	eq_str(t, SliceString("june", 0, 4), "june")
	eq_str(t, SliceString("june", 0, 2), "ju")
	eq_str(t, SliceString("june", 3, 5), "e")
	eq_str(t, SliceString("june", -1, 5), "e")
	eq_str(t, SliceString("june", 2, 5), "ne")
	eq_str(t, SliceString("june", -2, 5), "ne")
}

// ToAscii() is still to do
// func TestToAscii(t *testing.T) {
// 	s := "abcd♥"
// 	sExpected := "abcd\\u0098\\u2665"
// 	if ToAscii(s) != sExpected {
// 		t.Errorf("ToAscii(%s) expected %s, got %s", s, sExpected, ToAscii(s))
// 		t.Fail()
// 	}
// }

func TestMinMax(t *testing.T) {
	eq_int(t, Min(-1999, 0), -1999)
	eq_int(t, Min(0, -10), -10)
	eq_int(t, Min(0, 10), 0)
	eq_int(t, Max(-1999, 0), 0)
	eq_int(t, Max(0, -10), 0)
	eq_int(t, Max(0, 10), 10)
}

func TestPad(t *testing.T) {
	eq_str(t, PadLeft("foo", 5), "foo  ")
	eq_str(t, PadLeft("foo", 4), "foo ")
	eq_str(t, PadLeft("foo", 3), "foo")
	eq_str(t, PadLeft("foo", 2), "foo")
	eq_str(t, PadLeft("foo", 1), "foo")
	eq_str(t, PadLeft("foo", 0), "foo")
	eq_str(t, PadLeft("foo", -1), "foo")
	eq_str(t, PadLeft("foo", -10), "foo")
	eq_str(t, PadRight("foo", 5), "  foo")
	eq_str(t, PadRight("foo", 4), " foo")
	eq_str(t, PadRight("foo", 3), "foo")
	eq_str(t, PadRight("foo", 2), "foo")
	eq_str(t, PadRight("foo", 1), "foo")
	eq_str(t, PadRight("foo", 0), "foo")
	eq_str(t, PadRight("foo", -1), "foo")
	eq_str(t, PadRight("foo", -10), "foo")
}

func TestIndent(t *testing.T) {
	eq_str(t, Indent(-1), "")
	eq_str(t, Indent(0), "")
	eq_str(t, Indent(1), "  ")
	eq_str(t, Indent(2), "    ")
	eq_str(t, Indent(3), "      ")
}
