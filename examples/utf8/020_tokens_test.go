package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

// Lexical elements - https://go.dev/ref/spec#Lexical_elements
// - Comments TODO
// - Tokens
// - Semicolons TODO
// - Identifiers
// - Keywords
// - Operators and punctuation
// - Integer literals
// - Floating-point literals TODO
// - Imaginary literals TODO
// - Rune literals
// - String literals
func TestLexicalElements(t *testing.T) {
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
		duo("`\\n`", "raw_string_lit"),                                  // original example is `\n<Actual CR>\n` // same as "\\n\n\\n". But's a bit hard to reproduce...
		duo("\"i like guitar\"", "interpreted_string_lit"),              // this is an added example
		duo("\"i like \\\"bass\\\" guitar\"", "interpreted_string_lit"), // this is an added example
		duo(`"
"`, "interpreted_string_lit"),
		duo(`"\""`, "interpreted_string_lit"), // same as `"`
		duo("\"Hello, world!\n\"", "interpreted_string_lit"),
		duo(`"日本語"`, "interpreted_string_lit"),
		duo(`"\u65e5本\U00008a9e"`, "interpreted_string_lit"),
		duo(`"\xff\u00FF"`, "interpreted_string_lit"),
		// TODO duo(`"\uD800"`, "ERROR illegal: surrogate half"),
		// TODO duo(`"\U00110000"`, "ERROR illegal: invalid Unicode code point"),
		// these 5 following all represent the same thing:
		duo(`"日本語"`, "interpreted_string_lit"),                                  // UTF-8 input text
		duo("`日本語`", "raw_string_lit"),                                          // UTF-8 input text as a raw literal
		duo(`"\u65e5\u672c\u8a9e"`, "interpreted_string_lit"),                   // the explicit Unicode code points
		duo(`"\U000065e5\U0000672c\U00008a9e"`, "interpreted_string_lit"),       // the explicit Unicode code points
		duo(`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`, "interpreted_string_lit"), // the explicit UTF-8 bytes
	} {
		test(t, gm, pair)
	}
}
