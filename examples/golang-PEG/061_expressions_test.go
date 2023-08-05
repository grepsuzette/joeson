package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

// Expressions - https://go.dev/ref/spec#Expressions
func TestExpressions(t *testing.T) {
	gm := j.GrammarFromLines(
		"go-types-and-exprs",
		rules(
			o(named("__expression", partial_rules_expressions)),
			o(named("__type", partial_rules_types)),
		))

	test(t, gm, "math.Sin", "QualifiedIdent") // denotes the Sin function in package math
}
