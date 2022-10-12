package helpers

import (
	"sort"
)

func BisectLeft(a []int, n int) int {
	return _bisectLeft(a, n, 0, len(a))
}

func _bisectLeft(a []int, n int, low int, high int) int {
	s := a[low:high]
	return sort.Search(len(s), func(i int) bool {
		return s[i] >= n
	})
}

/*
 Return the index where to insert item n in list a, assuming a is sorted.

 The return value i is such that all e in a[:i] have e <= x, and all e in
 a[i:] have e > x.  So if x already appears in the list, a.insert(x) will
 insert just after the rightmost n already there.

 Optional to _bisectRight args are lo (default 0) and hi (default len(a)) bound
 the slice of a to be searched.
*/
func BisectRight(a []int, n int) int {
	return _bisectRight(a, n, 0, len(a))
}

func _bisectRight(a []int, n int, low int, high int) int {
	s := a[low:high]
	return sort.Search(len(s), func(i int) bool {
		return s[i] > n
	})
}

/*
   Insert item x in list a, and keep it sorted assuming a is sorted.

   If x is already in a, insert it to the right of the rightmost x.

   Optional args lo (default 0) and hi (default len(a)) bound the
   slice of a to be searched.
*/
func InsortRight(a []int, x int) []int {
	return insert(a, BisectRight(a, x), x)
}

func InsortLeft(a []int, x int) []int {
	return insert(a, BisectLeft(a, x), x)
}

type intSlice []int

func insert(a []int, at int, val int) []int {
	a = append(a, 0)
	copy(a[at+1:], a[at:])
	a[at] = val
	return a
}
