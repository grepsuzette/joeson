package helpers

import (
	"strconv"
	"testing"
)

func asseq_arr(t *testing.T, a []int, b []int) {
	if a == nil || b == nil || len(a) != len(b) {
		t.Fatal("different arrays")
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Fatal("arrays are different at index " + strconv.Itoa(i))
		}
	}
}

func TestBisectRight(t *testing.T) {
	a := []int{1, 2, 3, 5}
	a = InsortRight(a, 4)
	asseq_arr(t, a, []int{1, 2, 3, 4, 5})
	a = InsortLeft(a, 4)
	asseq_arr(t, a, []int{1, 2, 3, 4, 4, 5})
}
