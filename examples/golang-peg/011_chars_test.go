package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestCharacters(t *testing.T) {
	gm := j.GrammarFromLines(rules_chars, "go-chars")
	for _, pair := range []Duo{
		duo("345678", ""),
		duo("aaergega", ""),
		duo("_", ""),
	} {
		test(t, gm, pair)
	}
}
