package main

import (
	"fmt"
	"strings"
	"testing"

	j "github.com/grepsuzette/joeson"
)

const E = "should not happen" // panic(E)

// -- The parsing grammar

func i(a ...any) j.ILine                                 { return j.I(a...) }
func o(a ...any) j.OLine                                 { return j.O(a...) }
func rules(a ...j.Line) []j.Line                         { return a }
func named(name string, lineStringOrAst any) j.NamedRule { return j.Named(name, lineStringOrAst) }

func main() {
	gm_chars := j.GrammarFromLines("go-characters", rules_chars)
	for _, s := range []string{"345678", "aaegeagr", "_"} {
		ast := gm_chars.ParseString(s)
		if j.IsParseError(ast) {
			panic(ast.String())
		}
	}
	gm_tokens := j.GrammarFromLines("go-tokens", rules_tokens)
	gm_tokens.ParseString("aawfe8f2")
}

// Parse all keys from map `h` sequentially using grammar `gm`.
// `test()` will be called for each key-value pair.
func testParse(t *testing.T, gm *j.Grammar, h map[string]string) {
	t.Helper()
	for k, v := range h {
		test(t, gm, k, v)
	}
}

// Parse `s` using grammar `gm`, asserts `s` is parsed
// with an ast type beginning with string `expect`.
//
// That last one is given using `ast.String()`, where `ast` is the
// result of `gm.ParseString(s)`.
//
// Two special cases:
//
//  1. When `expect` starts with "ERROR", it means a parse error is expected.
//     You can specify the exact error like so: "ERROR illegal: octal value over 255".
//
// 2. When `expect` is "", the test passes as long as parsing did not fail.
func test(t *testing.T, gm *j.Grammar, s string, expect string) {
	t.Helper()
	ast := gm.ParseString(s)
	if j.IsParseError(ast) {
		if strings.HasPrefix(expect, "ERROR") {
			fmt.Printf("[32m%s[0m gave an error as expected [32m✓[0m\n", s)
		} else {
			t.Fatalf("Error parsing %s. Expected ast.String() to contain '%s', got '%s'", s, expect, ast.String())
		}
	} else {
		if strings.Contains(ast.String(), expect) {
			fmt.Printf("[32m%s[0m parsed as [33m%s[0m [32m✓[0m %s\n", s, ast.String(), expect)
		} else {
			t.Fatalf(
				"Error, \"[1m%s[0m\" [1;31mparsed[0m as %s [1;31mbut expected [0;31m%s[0m",
				s,
				ast.String(),
				expect,
			)
		}
	}
}
