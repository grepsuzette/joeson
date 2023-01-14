package helpers

import (
	"testing"
)

func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Errorf("should have panicked")
}

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

// test niladic lazy.go, ctor gets the callback this time
func TestLazy0b(t *testing.T) {
	lazy := NewLazy[int](func() int { return 5 })
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
