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
	lazy := NewLazy0[int]()
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
	lazy := NewLazy0[int](func() int { return 5 })
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

// test monadic lazy1.go, ctor with no argument
func TestLazy1a(t *testing.T) {
	lazy1 := NewLazy1[string, int]()
	lazy1.SetLazy(func(n int) string {
		s := ""
		for i := 0; i < n; i++ {
			s += "a"
		}
		return s
	})
	if lazy1.IsSet() {
		t.FailNow()
	}
	if lazy1.Get(3) != "aaa" {
		t.FailNow()
	}
	if lazy1.Get(3) != "aaa" {
		t.FailNow()
	}
	shouldPanic(t, func() { lazy1.Get(5) })
}

// test monadic lazy1.go, ctor with the callback
func TestLazy1b(t *testing.T) {
	lazy1 := NewLazy1[string, int](func(n int) string {
		s := ""
		for i := 0; i < n; i++ {
			s += "a"
		}
		return s
	})
	if lazy1.IsSet() {
		t.FailNow()
	}
	if lazy1.Get(5) != "aaaaa" {
		t.FailNow()
	}
	if lazy1.Get(5) != "aaaaa" {
		t.FailNow()
	}
	shouldPanic(t, func() { lazy1.Get(4) })
}
