package helpers

// lazy, when lazy.Get() is called:
// - returns lazy.val if set,
// - returns lazy.val = lazy.f() otherwise

type Lazy[T any] struct {
	val *T
	f   func() T // lazy evaluator
}

func NewLazy[T any]() *Lazy[T]                { return &Lazy[T]{} }
func LazyFromValue[T any](t T) *Lazy[T]       { r := &Lazy[T]{}; r.Set(t); return r }
func LazyFromFunc[T any](f func() T) *Lazy[T] { return &Lazy[T]{f: f} }

func (k *Lazy[T]) IsSet() bool { return k.val != nil }
func (k *Lazy[T]) Set(t T)     { k.val = &t }
func (k *Lazy[T]) Clear()      { k.val = nil; k.f = nil }
func (k *Lazy[T]) Get() T {
	if !k.IsSet() && k.f != nil {
		ptr := k.f()
		(*k).val = &ptr
		return *k.val
	} else {
		return *k.val
	}
}

func (k *Lazy[T]) SetLazy(f func() T) { k.f = f }
