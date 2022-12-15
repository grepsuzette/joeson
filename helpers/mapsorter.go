package helpers

import "sort"

// this brings determinism, golang maps being unordered

func SortIntKeys[V any](h map[int]V) []int {
	a := []int{}
	for k, _ := range h {
		a = append(a, k)
	}
	sort.Ints(a)
	return a
}

func SortStringKeys[V any](h map[string]V) []string {
	a := []string{}
	for k, _ := range h {
		a = append(a, k)
	}
	sort.Strings(a)
	return a
}
