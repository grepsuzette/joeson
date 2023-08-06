package joeson

import "testing"

func TestNativeMapGetIntExists(t *testing.T) {
	nm := NewEmptyNativeMap()
	nm.Set("int", NewNativeInt(99))
	nm.Set("int2", NewNativeString("99"))
	if n, ok := nm.GetIntExists("int"); ok {
		if n != 99 {
			t.Fail()
		}
	} else {
		t.Fail()
	}
	if n, ok := nm.GetIntExists("int2"); ok {
		if n != 99 {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
