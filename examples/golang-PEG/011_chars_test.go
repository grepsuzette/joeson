package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestCharacters(t *testing.T) {
	testParse(t, j.GrammarFromLines(rules_chars, "go-chars"), map[string]string{
		//<string to parse>: <string beginning by "ERROR" if it must fail>
		// note:            we leave these type empty "" here,
		//                  as the rules for 010_chars are still a bit naive.
		//                  Refer to 021_tokens_test.go for better examples
		"\n":               "", // newline
		"\\x{20}":          "", // unicode_char
		"a":                "", // unicode_char
		"7":                "", // rule unicode_digit
		"abcdefghIJKLMNOP": "", // "
		"_":                "", // letter
		"1234567":          "", // digits
	})
}
