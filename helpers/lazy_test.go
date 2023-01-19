package helpers

import (
	"testing"
)

func TestLazy0a(t *testing.T) {
	lazy := NewLazy[int]()
	lazy.SetLazy(func() int { return 5 })
	if lazy.IsSet() {
		t.FailNow()
	}
	if lazy.Get() != 5 {
		t.FailNow()
	}
	if !lazy.IsSet() {
		t.FailNow()
	}
	lazy.Set(3)
	if lazy.Get() != 3 {
		t.FailNow()
	}
}

func TestLazy0b(t *testing.T) {
	lazy := NewLazyFromFunc[int](func() int { return 5 })
	if lazy.IsSet() {
		t.FailNow()
	}
	if lazy.Get() != 5 {
		t.FailNow()
	}
	if !lazy.IsSet() {
		t.FailNow()
	}
	lazy.Set(3)
	if lazy.Get() != 3 {
		t.FailNow()
	}
}

func TestLazy0c(t *testing.T) {
	lazy := NewLazyFromValue[int](5)
	if !lazy.IsSet() {
		t.FailNow()
	}
	if lazy.Get() != 5 {
		t.FailNow()
	}
	lazy.Set(3)
	if lazy.Get() != 3 {
		t.FailNow()
	}
}
