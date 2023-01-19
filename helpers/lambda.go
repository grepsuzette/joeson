package helpers

// functional map over an array
func AMap[T any, U any](a []T, f func(T) U) []U {
	r := make([]U, len(a))
	for i, x := range a {
		r[i] = f(x)
	}
	return r
}
