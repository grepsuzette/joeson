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
	testParse(t, j.GrammarFromLines(rules_tokens, "go-tokens"), map[string]string{
		"break":                "keyword",
		"default":              "keyword",
		"func":                 "keyword",
		"interface":            "keyword",
		"select":               "keyword",
		"case":                 "keyword",
		"defer":                "keyword",
		"go":                   "keyword",
		"map":                  "keyword",
		"struct":               "keyword",
		"chan":                 "keyword",
		"else":                 "keyword",
		"goto":                 "keyword",
		"package":              "keyword",
		"switch":               "keyword",
		"const":                "keyword",
		"fallthrough":          "keyword",
		"if":                   "keyword",
		"range":                "keyword",
		"type":                 "keyword",
		"continue":             "keyword",
		"for":                  "keyword",
		"import":               "keyword",
		"return":               "keyword",
		"var":                  "keyword",
		"_":                    "identifier",
		"+":                    "operator",
		"&":                    "operator",
		"+=":                   "operator",
		"&=":                   "operator",
		"&&":                   "operator",
		"==":                   "operator",
		"!=":                   "operator",
		"(":                    "operator",
		")":                    "operator",
		"-":                    "operator",
		"|":                    "operator",
		"-=":                   "operator",
		"|=":                   "operator",
		"||":                   "operator",
		"<":                    "operator",
		"<=":                   "operator",
		"[":                    "operator",
		"]":                    "operator",
		"*":                    "operator",
		"^":                    "operator",
		"*=":                   "operator",
		"^=":                   "operator",
		"<-":                   "operator",
		">":                    "operator",
		">=":                   "operator",
		"{":                    "operator",
		"}":                    "operator",
		"/":                    "operator",
		"<<":                   "operator",
		"/=":                   "operator",
		"<<=":                  "operator",
		"++":                   "operator",
		"=":                    "operator",
		":=":                   "operator",
		"%":                    "operator",
		">>":                   "operator",
		"%=":                   "operator",
		">>=":                  "operator",
		"--":                   "operator",
		"!":                    "operator",
		"...":                  "operator",
		"&^":                   "operator",
		"&^=":                  "operator",
		"~":                    "operator",
		".":                    "punctuation",
		";":                    "punctuation",
		":":                    "punctuation",
		"aaegeagr":             "identifier",
		"345678":               "int_lit",
		"42":                   "int_lit",
		"4_2":                  "int_lit",
		"0600":                 "ERROR:i think",
		"0_600":                "ERROR:i think",
		"0o600":                "octal_lit",
		"0O600":                "octal_lit",
		"0xBadFace":            "hex_lit",
		"0xBad_Face":           "hex_lit",
		"0x_67_7a_2f_cc_40_c6": "hex_lit",
		"170141183460469231731687303715884105727":       "int_lit",
		"170_141183_460469_231731_687303_715884_105727": "int_lit",
		"_42":                   "identifier",
		"42_":                   "ERROR invalid: _ must separate successive digits",
		"4__2":                  "ERROR invalid: only one _ at a time",
		"0_xBadFace":            "ERROR invalid: _ must separate successive digits",
		"0b0101001010011001":    "binary_lit",
		"0B01_01001010_01_1001": "binary_lit",
		"0o1013375762602732":    "octal_lit",
		"0O1013375762602732726": "octal_lit",
		"0o19":                  "ERROR invalid octal digit",
		"'\\125'":               "rune_lit",
		"'\\xF2'":               "rune_lit",
		"'\a'":                  "escaped_char",
		// "'\b'": "escaped_char", // skipped BELL RING as it doesn't work
		"'\f'":           "escaped_char",
		"'\n'":           "escaped_char",
		"'\r'":           "escaped_char",
		"'\t'":           "escaped_char",
		"'\v'":           "escaped_char",
		"'\\u13F8'":      "little_u_value",
		"'\\u13a'":       "ERROR little_u_value requires 4 hex",
		"'\\u1a248'":     "ERROR little_u_value requires 4 hex",
		"'\\UFFeeFFee'":  "big_u_value",
		"'\\UFFeeFFe'":   "ERROR big_u_value requires 8 hex",
		"'\\UFFeeFFeeA'": "ERROR big_u_value requires 8 hex",
		"'a'":            "rune_lit",
		"'ä'":            "rune_lit",
		"'本'":            "rune_lit",
		"'\\000'":        "octal_byte_value",
		"'\\007'":        "octal_byte_value",
		"'\\x07'":        "hex_byte_value",
		"'\\xff'":        "hex_byte_value",
		"'\\u12e4'":      "little_u_value",
		"'\\U00101234'":  "big_u_value",
		`'`:              "rune_lit", // rune literal containing single quote character
		"'aa'":           "ERROR illegal: too many characters",
		"'\\k'":          "ERROR illegal: k is not recognized after a backslash",
		"'\\xa'":         "ERROR illegal: too few hexadecimal digits",
		"'\\0'":          "ERROR illegal: too few octal digits",
		"'\\400'":        "ERROR illegal: octal value over 255",
		// "'\\uDFFF'": "ERROR illegal: surrogate half", // TODO
		// "'\\U00110000'": "ERROR illegal: invalid Unicode code point", // TODO
		// -- string_lit -- tests adapted from https://go.dev/ref/spec#String_literals
		"`abc`":                          "raw_string_lit",
		"`\\n`":                          "raw_string_lit",         // original example is `\n<Actual CR>\n` // same as "\\n\n\\n". But's a bit hard to reproduce...
		"\"i like guitar\"":              "interpreted_string_lit", // this is an added example
		"\"i like \\\"bass\\\" guitar\"": "interpreted_string_lit", // this is an added example
		`"
"`: "interpreted_string_lit",
		`"\""`:                "interpreted_string_lit", // same as `"`
		"\"Hello: world!\n\"": "interpreted_string_lit",
		`"日本語"`:               "interpreted_string_lit",
		`"\u65e5本\U00008a9e"`: "interpreted_string_lit",
		`"\xff\u00FF"`:        "interpreted_string_lit",
		// TODO `"\uD800"`: "ERROR illegal: surrogate half",
		// TODO `"\U00110000"`: "ERROR illegal: invalid Unicode code point",
		// these 5 following all represent the same thing:
		"`日本語`":                                  "raw_string_lit",         // UTF-8 input text as a raw literal
		`"\u65e5\u672c\u8a9e"`:                   "interpreted_string_lit", // the explicit Unicode code points
		`"\U000065e5\U0000672c\U00008a9e"`:       "interpreted_string_lit", // the explicit Unicode code points
		`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`: "interpreted_string_lit", // the explicit UTF-8 bytes
	})
}
