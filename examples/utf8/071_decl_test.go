package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestDecls(t *testing.T) {
	gm := j.GrammarFromLines(rules_decl, "go-declaraions")
	for _, pair := range []Duo{
		duo("const a int = 1", ""),
		// duo("const Pi float64 = ", "3.14159265358979323846"), // TODO enable when float supported
		// duo("const zero = 0.0", ""), // TODO enable when float supported
		// duo("const (
		// 	size int64 = 1024
		// 	eof        = -1
		// )", ""),
		// const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", untyped integer and string constants
		// const u, v float32 = 0, 3    // u = 0.0, v = 3.0
	} {
		test(t, gm, pair)
	}
}