package helpers

// variable caching (can be set, or not)
// E.g.:
// var precalculated Varcache[string]
// return precalculated.GetCacheOrSet(func (s string) string {
//	  return "concatenation " + "happens " + "only once"
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
