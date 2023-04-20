package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestTypes(t *testing.T) {
	gm := j.GrammarFromLines(rules(
		o(named("Type", partial_rules_types)),
		o(named("Expression", partial_rules_expressions)),
	), "go-types-and-exprs")
	for _, pair := range []Duo{
		duo("[32]byte", "ArrayType"),
		duo("[]32byte", "ERROR "),
		// duo("[2*N] struct { x, y int32 }", "?"),
		// duo("_", ""),
	} {
		test(t, gm, pair)
	}
}
