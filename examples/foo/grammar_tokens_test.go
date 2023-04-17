package main

import (
	"fmt"
	"strings"
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestTokens(t *testing.T) {
	gm := j.GrammarFromLines(rules_tokens, "go-tokens")
	for _, pair := range []Duo{
		duo("QqQ", "quoted"),
		// duo("QQQ", "quoted"),
		// duo("QQ", "quoted"),
	} {
		if ast, e := gm.ParseString(pair.a, j.ParseOptions{Debug: true}); e != nil {
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
}
