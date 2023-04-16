package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestCharacters(t *testing.T) {
	gm := j.GrammarFromLines(rules_chars, "go-characters")
	for _, s := range []string{"345678", "aaegeagr", "_"} {
		if _, e := gm.ParseString(s); e != nil {
			t.Fatalf("Error parsing %s", s)
		}
	}
}
