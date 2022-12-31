package helpers

// specialized lazy varcache
// - it requires to be instantiated with NewLazy1[Type1,Type1]()
// - uses a monadic callback (with exactly one argument), it is set with SetLazy(cb)
// - when called with Get(p),
//     the first time p is memorized,
//     the result is then given in function of p.
//     When p is different from the first time, it panics!

type Lazy1[T any, P comparable] struct {
	val      *T
	f        func(P) T
	initialP P
}

// the callback can optionally be passed, or SetLazy()
// can be called later instead
func NewLazy1[T any, P comparable](af ...func(P) T) *Lazy1[T, P] {
	if len(af) == 0 {
		return &Lazy1[T, P]{}
	} else {
		return &Lazy1[T, P]{f: af[0]}
	}
}
func (k *Lazy1[T, P]) IsSet() bool { return k.val != nil }
func (k *Lazy1[T, P]) Set(t T)     { k.val = &t }
func (k *Lazy1[T, P]) Clear()      { k.val = nil; k.f = nil }
func (k *Lazy1[T, P]) Get(p P) T {
	// when !IsSet(), Get() calls the
	// Lazy.f callback if not nil
	if !k.IsSet() && k.f != nil {
		k.initialP = p
		ptr := k.f(p)
		(*k).val = &ptr
		return *k.val
	} else {
		if k.initialP == p {
			return *k.val
		} else {
			panic("lazy1:Get(): different argument from the 1st time is forbidden")
		}
	}
}

func (k *Lazy1[T, P]) SetLazy(f func(P) T) {
	k.f = f
}
