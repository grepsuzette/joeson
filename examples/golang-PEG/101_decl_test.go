package main

import (
	"testing"
	// j "github.com/grepsuzette/joeson"
)

func TestDecls(t *testing.T) {
	// gm := j.GrammarFromLines(rules_decl, "go-declarations")
	// test(t, gm, "const a int = 1", "")
	// test(t, gm, "const Pi float64 = ", "3.14159265358979323846"), // TODO enable when float supported
	// test(t, gm, "const zero = 0.0", ""), // TODO enable when float supported
	// test(t, gm, "const (
	// 	size int64 = 1024
	// 	eof        = -1
	// )", ""),
	// const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", untyped integer and string constants
	// const u, v float32 = 0, 3    // u = 0.0, v = 3.0
	// } {
	// sourceCode > tokenizedString (i.e. []ParseContext) >
	// gm.ParseString(pair.a, WithTokenScanner(sCaN))
	// ast := gm.Parse(ctx)
	// ast := gm.ParseTokens(pair.a, scanner)
	// 	if strings.HasPrefix(pair.b, "ERROR") {
	// 		fmt.Printf("[32m%s[0m gave an error as expected [32m✓[0m\n", pair.a)
	// 	} else {
	// 		t.Fatalf("Error parsing %s. Expected ast.ContentString() to contain '%s', got '%s'", pair.a, pair.b, e.Error())
	// 	}
	// } else {
	// 	if strings.Contains(ast.ContentString(), pair.b) {
	// 		fmt.Printf("[32m%s[0m parsed as [33m%s[0m [32m✓[0m %s\n", pair.a, ast.ContentString(), pair.b)
	// 	} else {
	// 		t.Fatalf(
	// 			"Error, \"[1m%s[0m\" [1;31mparsed[0m as %s [1;31mbut expected [0;31m%s[0m",
	// 			pair.a,
	// 			ast.ContentString(),
	// 			pair.b,
	// 		)
	// 	}
	// }
	// }
}
