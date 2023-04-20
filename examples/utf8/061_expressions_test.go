package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

// Expressions - https://go.dev/ref/spec#Expressions
func TestExpressions(t *testing.T) {
	gm := j.GrammarFromLines(rules(
		o(named("__expression", partial_rules_expressions)),
		o(named("__type", partial_rules_types)),
	), "go-types-and-exprs")
	for _, pair := range []Duo{
		duo("math.Sin", "QualifiedIdent"), // denotes the Sin function in package math
	} {
		test(t, gm, pair)
	}
}
