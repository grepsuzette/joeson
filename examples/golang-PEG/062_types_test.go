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

	test(t, gm, "[32]byte", "ArrayType")
	test(t, gm, "[]32byte", "ERROR ")

	// TODO uncomment when Type finished
	// test(t, gm, "[2*N] struct { x, y int32 }", ""),
	// test(t, gm, "[1000]*float64", ""),
	// test(t, gm, "[3][5]int", ""),
	// test(t, gm, "[2][2][2]float64", ""),  // same as [2]([2]([2]float64))
	// TODO
	// // invalid array types
	// type (
	//         T1 [10]T1                 // element type of T1 is T1
	//         T2 [10]struct{ f T2 }     // T2 contains T2 as component of a struct
	//         T3 [10]T4                 // T3 contains T3 as component of a struct in T4
	//         T4 struct{ f T3 }         // T4 contains T4 as component of array T3 in a struct
	// )
	// // valid array types
	// type (
	//         T5 [10]*T5                // T5 contains T5 as component of a pointer
	//         T6 [10]func() T6          // T6 contains T6 as component of a function type
	//         T7 [10]struct{ f []T7 }   // T7 contains T7 as component of a slice in a struct
	// )
}
