package helpers

// a simple cache in a variable (can be set, or not set),
// helping for laziness

// A bad example:

// var precalculated Varcache[string]
// return precalculated.GetCacheOrSet(func (s string) string {
//	  return "concat " + "virtually happening " + "only once"
// })

type Varcache[T any] struct{ val *T }

func (k Varcache[T]) IsCacheSet() bool { return k.val != nil }
func (k Varcache[T]) SetCache(t T)     { k.val = &t }
func (k Varcache[T]) GetCache() T      { return *k.val }
func (k Varcache[T]) GetCacheOrSet(f func() T) T {
	if k.val == nil {
		v := f()
		k.val = &v
	}
	return *k.val
}
