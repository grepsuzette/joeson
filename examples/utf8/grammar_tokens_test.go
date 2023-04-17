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
		duo("break", "keyword"),
		duo("default", "keyword"),
		duo("func", "keyword"),
		duo("interface", "keyword"),
		duo("select", "keyword"),
		duo("case", "keyword"),
		duo("defer", "keyword"),
		duo("go", "keyword"),
		duo("map", "keyword"),
		duo("struct", "keyword"),
		duo("chan", "keyword"),
		duo("else", "keyword"),
		duo("goto", "keyword"),
		duo("package", "keyword"),
		duo("switch", "keyword"),
		duo("const", "keyword"),
		duo("fallthrough", "keyword"),
		duo("if", "keyword"),
		duo("range", "keyword"),
		duo("type", "keyword"),
		duo("continue", "keyword"),
		duo("for", "keyword"),
		duo("import", "keyword"),
		duo("return", "keyword"),
		duo("var", "keyword"),
		duo("_", "identifier"),
		duo("+", "operator"),
		duo("&", "operator"),
		duo("+=", "operator"),
		duo("&=", "operator"),
		duo("&&", "operator"),
		duo("==", "operator"),
		duo("!=", "operator"),
		duo("(", "operator"),
		duo(")", "operator"),
		duo("-", "operator"),
		duo("|", "operator"),
		duo("-=", "operator"),
		duo("|=", "operator"),
		duo("||", "operator"),
		duo("<", "operator"),
		duo("<=", "operator"),
		duo("[", "operator"),
		duo("]", "operator"),
		duo("*", "operator"),
		duo("^", "operator"),
		duo("*=", "operator"),
		duo("^=", "operator"),
		duo("<-", "operator"),
		duo(">", "operator"),
		duo(">=", "operator"),
		duo("{", "operator"),
		duo("}", "operator"),
		duo("/", "operator"),
		duo("<<", "operator"),
		duo("/=", "operator"),
		duo("<<=", "operator"),
		duo("++", "operator"),
		duo("=", "operator"),
		duo(":=", "operator"),
		duo("%", "operator"),
		duo(">>", "operator"),
		duo("%=", "operator"),
		duo(">>=", "operator"),
		duo("--", "operator"),
		duo("!", "operator"),
		duo("...", "operator"),
		duo("&^", "operator"),
		duo("&^=", "operator"),
		duo("~", "operator"),
		duo(".", "punctuation"),
		duo(";", "punctuation"),
		duo(":", "punctuation"),
		duo(",", "punctuation"),
		duo("345678", "literal"),
		duo("aaegeagr", "identifier"),
		duo("345678", "int_lit"),
		duo("42", "int_lit"),
		duo("4_2", "int_lit"),
		duo("0600", "ERROR:i think"),
		duo("0_600", "ERROR:i think"),
		duo("0o600", "octal_lit"),
		duo("0O600", "octal_lit"),
		duo("0xBadFace", "hex_lit"),
		duo("0xBad_Face", "hex_lit"),
		duo("0x_67_7a_2f_cc_40_c6", "hex_lit"),
		duo("170141183460469231731687303715884105727", "int_lit"),
		duo("170_141183_460469_231731_687303_715884_105727", "int_lit"),
		duo("_42", "identifier"),
		duo("42_", "ERROR invalid: _ must separate successive digits"),
		duo("4__2", "ERROR invalid: only one _ at a time"),
		duo("0_xBadFace", "ERROR invalid: _ must separate successive digits"),
		duo("0b0101001010011001", "binary_lit"),
		duo("0B01_01001010_01_1001", "binary_lit"),
		duo("0o1013375762602732", "octal_lit"),
		duo("0O1013375762602732726", "octal_lit"),
		duo("0o19", "ERROR invalid octal digit"),
		duo("0B01_01001010_01_1001", "binary_lit"),
		duo("'\\125'", "rune_lit"),
		duo("'\\xF2'", "rune_lit"),
		duo("'\a'", "escaped_char"),
		// duo("'\b'", "escaped_char"), // skipped BELL RING as it doesn't work
		duo("'\f'", "escaped_char"),
		duo("'\n'", "escaped_char"),
		duo("'\r'", "escaped_char"),
		duo("'\t'", "escaped_char"),
		duo("'\v'", "escaped_char"),
		duo("'\\u13F8'", "little_u_value"),
		duo("'\\u13a'", "ERROR little_u_value requires 4 hex"),
		duo("'\\u1a248'", "ERROR little_u_value requires 4 hex"),
		duo("'\\UFFeeFFee'", "big_u_value"),
		duo("'\\UFFeeFFe'", "ERROR big_u_value requires 8 hex"),
		duo("'\\UFFeeFFeeA'", "ERROR big_u_value requires 8 hex"),
		duo("'a'", "rune_lit"),
		duo("'ä'", "rune_lit"),
		duo("'本'", "rune_lit"),
		duo("'a'", "unicode_char"),
		duo("'ä'", "unicode_char"),
		duo("'本'", "unicode_char"),
		duo("'\\000'", "octal_byte_value"),
		duo("'\\007'", "octal_byte_value"),
		duo("'\\x07'", "hex_byte_value"),
		duo("'\\xff'", "hex_byte_value"),
		duo("'\\u12e4'", "little_u_value"),
		duo("'\\U00101234'", "big_u_value"),
		duo(`'`, "rune_lit"), // rune literal containing single quote character
		duo("'aa'", "ERROR illegal: too many characters"),
		duo("'\\k'", "ERROR illegal: k is not recognized after a backslash"),
		duo("'\\xa'", "ERROR illegal: too few hexadecimal digits"),
		duo("'\\0'", "ERROR illegal: too few octal digits"),
		duo("'\\400'", "ERROR illegal: octal value over 255"),
		// duo("'\\uDFFF'", "ERROR illegal: surrogate half"), // TODO
		// duo("'\\U00110000'", "ERROR illegal: invalid Unicode code point"), // TODO
		// -- string_lit -- tests adapted from https://go.dev/ref/spec#String_literals
		duo("`abc`", "raw_string_lit"),
		duo("`\\n`", "raw_string_lit"), // original example is `\n<Actual CR>\n` // same as "\\n\n\\n". But's a bit hard to reproduce...
		duo(`"
"`, "interpreted_string_lit"),
		duo(`"`+"\""+`"`, "interpreted_string_lit"), // same as `"`
		duo(`"`+"Hello, world!\\n"+`"`, "interpreted_string_lit"),
	} {
		if ast, e := gm.ParseString(pair.a); e != nil {
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
