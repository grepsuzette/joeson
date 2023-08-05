package joeson

import (
	"testing"
)

func assertOk(t *testing.T, ast Ast) {
	t.Helper()
	if IsParseError(ast) {
		t.Error()
	}
}

func assertError(t *testing.T, ast Ast) {
	t.Helper()
	if !IsParseError(ast) {
		t.Error()
	}
}

func TestLookahead(t *testing.T) {
	gm := GrammarFromLines([]Line{
		o(named("Input", "foo_followed_by_b 'bar'")),
		i(named("foo_followed_by_b", "foo ?('b')")), // the 'b' is not consumed (not "captured")
		i(named("foo", "'foo'")),
	}, "lookahead")
	assertError(t, gm.ParseString("fool"))
	assertOk(t, gm.ParseString("foobar"))
	assertError(t, gm.ParseString("foo  bar"))
	assertError(t, gm.ParseString("foo - bar"))
}
