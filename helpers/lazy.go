package helpers

// niladic lazy varcache, i.e. the callback has no argument

type Lazy[T any] struct {
	val *T
	f   func() T // lazy evaluator
}

func NewLazy[T any]() *Lazy[T]                   { return &Lazy[T]{} }
func NewLazyFromValue[T any](t T) *Lazy[T]       { r := &Lazy[T]{}; r.Set(t); return r }
func NewLazyFromFunc[T any](f func() T) *Lazy[T] { return &Lazy[T]{f: f} }

func (k *Lazy[T]) IsSet() bool { return k.val != nil }
func (k *Lazy[T]) Set(t T)     { k.val = &t }
func (k *Lazy[T]) Clear()      { k.val = nil; k.f = nil }
func (k *Lazy[T]) Get() T {
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

func (k *Lazy[T]) SetLazy(f func() T) { k.f = f }
