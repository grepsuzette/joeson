package main

import (
	// "fmt"
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
