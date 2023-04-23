package main

import (
	"fmt"
	"strings"
	"testing"

	j "github.com/grepsuzette/joeson"
)

const E = "should not happen" // panic(E)

// -- The parsing grammar

func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }
func named(name string, lineStringOrAst any) j.NamedRule {
	return j.Named(name, lineStringOrAst)
}
func rules(a ...j.Line) []j.Line { return a }

func main() {
	gm_chars := j.GrammarFromLines(rules_chars, "go-characters")
	for _, s := range []string{"345678", "aaegeagr", "_"} {
		if _, e := gm_chars.ParseString(s); e != nil {
			panic("parse error on " + s)
		}
	}
	gm_tokens := j.GrammarFromLines(rules_tokens, "go-tokens")
	gm_tokens.ParseString("aawfe8f2")
}

// test() is untokenized raw parsing, it is used to parse small units.
//
// It means "use `gm` grammar to `gm.ParseString(pair.a)`".
// The result if successful is expected to contain `pair.b`.
// If `pair.b` starts with "ERROR ", we instead expect an error.
func test(t *testing.T, gm *j.Grammar, pair Duo) {
	t.Helper()
	if ast, e := gm.ParseString(pair.a); e != nil {
		if strings.HasPrefix(pair.b, "ERROR") {
			fmt.Printf("[32m%s[0m gave an error as expected [32m✓[0m\n", pair.a)
		} else {
			t.Fatalf("Error parsing %s. Expected ast.ContentString() to contain '%s', got '%s'", pair.a, pair.b, e.Error())
		}
	} else {
		if strings.Contains(ast.ContentString(), pair.b) {
			fmt.Printf("[32m%s[0m parsed as [33m%s[0m [32m✓[0m %s\n", pair.a, ast.ContentString(), pair.b)
		} else {
			t.Fatalf(
				"Error, \"[1m%s[0m\" [1;31mparsed[0m as %s [1;31mbut expected [0;31m%s[0m",
				pair.a,
				ast.ContentString(),
				pair.b,
			)
		}
	}
}
