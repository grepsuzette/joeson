package helpers

import "testing"

func asseq(t *testing.T, a string, b string) {
	if a != b {
		t.Errorf("assert failed %s == %s", a, b)
		panic("a")
	}
}

func TestSliceString(t *testing.T) {
	for i := -5; i < 10; i++ {
		asseq(t, SliceString("june", i, 0), "")
	}
	asseq(t, SliceString("june", -2, 1), "")
	asseq(t, SliceString("june", -1, 1), "")
	asseq(t, SliceString("june", 0, 1), "j")
	asseq(t, SliceString("june", 1, 1), "")
	asseq(t, SliceString("june", 2, 1), "")
	asseq(t, SliceString("june", 3, 1), "")
	asseq(t, SliceString("june", -2, 2), "")
	asseq(t, SliceString("june", -1, 2), "")
	asseq(t, SliceString("june", 0, 2), "ju")
	asseq(t, SliceString("june", 1, 2), "u")
	asseq(t, SliceString("june", 2, 2), "")
	asseq(t, SliceString("june", 3, 2), "")
	asseq(t, SliceString("june", -6, 3), "jun")
	asseq(t, SliceString("june", -5, 3), "jun")
	asseq(t, SliceString("june", -4, 3), "jun")
	asseq(t, SliceString("june", -3, 3), "un")
	asseq(t, SliceString("june", -2, 3), "n")
	asseq(t, SliceString("june", -1, 3), "")
	asseq(t, SliceString("june", 0, 3), "jun")
	asseq(t, SliceString("june", 1, 3), "un")
	asseq(t, SliceString("june", 2, 3), "n")
	asseq(t, SliceString("june", 3, 3), "")
	asseq(t, SliceString("june", -5, 4), "june")
	asseq(t, SliceString("june", -5, 3), "jun")
	asseq(t, SliceString("june", -5, 2), "ju")
	asseq(t, SliceString("june", -5, 1), "j")
	asseq(t, SliceString("june", -5, 0), "")
	asseq(t, SliceString("june", -5, -1), "jun")
	asseq(t, SliceString("june", -5, -2), "ju")
	asseq(t, SliceString("june", -5, -3), "j")
	asseq(t, SliceString("june", -5, -4), "")
	asseq(t, SliceString("june", 0, 4), "june")
	asseq(t, SliceString("june", 0, 2), "ju")
	asseq(t, SliceString("june", 3, 5), "e")
	asseq(t, SliceString("june", -1, 5), "e")
	asseq(t, SliceString("june", 2, 5), "ne")
	asseq(t, SliceString("june", -2, 5), "ne")
}
