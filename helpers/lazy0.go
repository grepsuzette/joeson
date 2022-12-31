package helpers

// niladic lazy varcache, i.e. the callback has no argument

// The difference with varcache is that
// when !IsCacheSet(), GetCache() calls the
// Lazy.f callback if not nil.

type Lazy0[T any] struct {
	val *T
	f   func() T // lazy evaluator
}

func NewLazy0[T any](af ...func() T) *Lazy0[T] {
	if len(af) == 0 {
		return &Lazy0[T]{}
	} else {
		return &Lazy0[T]{f: af[0]}
	}
}
func (k *Lazy0[T]) IsSet() bool { return k.val != nil }

func (k *Lazy0[T]) Set(t T) { k.val = &t }
func (k *Lazy0[T]) Clear()  { k.val = nil; k.f = nil }
func (k *Lazy0[T]) Get() T {
	// The difference with varcache is that
	// when !IsSet(), lazy's Get() calls the
	// Lazy0.f callback if not nil...
	if !k.IsSet() && k.f != nil {
		ptr := k.f()
		(*k).val = &ptr
		return *k.val
	} else {
		return *k.val
	}
}

func (k *Lazy0[T]) SetLazy(f func() T) { k.f = f }

// doubt it's useful
// func (k *Lazy0[T]) GetOrSet(f func() T) T {
// 	if k.val == nil {
// 		k.f = f
// 		v := f()
// 		k.val = &v
// 	}
// 	return *k.val
// }
