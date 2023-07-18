package main

import (
	j "github.com/grepsuzette/joeson"
)

// https://go.dev/ref/spec#Characters

// TODO Floating-point literals
// TODO Imaginary literals (maybe)
// TODO String literals

// "Tokens form the vocabulary of the Go language. There are four classes:
// identifiers, keywords, operators and punctuation, and literals (SIC: that
// makes 5). White space,
// formed from spaces (U+0020), horizontal tabs (U+0009), carriage returns
// (U+000D), and newlines (U+000A), is ignored except as it separates tokens
// that would otherwise combine into a single token. Also, a newline or end of
// file may trigger the insertion of a semicolon. While breaking the input into
// tokens, the next token is the longest sequence of characters that form
// a valid token."

// Comments represent original spec, where useful (we often optimized it or reordered it for PEG)
var rules_tokens = rules(
	o(named("token", "keyword | identifier | operator | punctuation | literal"), x("token")),
	o(named("characters", rules_chars)), // import previous rules from grammar_chars.go
	i(named("keyword", "'break' | 'default' | 'func' | 'interface' | 'select' | 'case' | 'defer' | 'goto' | 'map' | 'struct' | 'chan' | 'else' | 'go' | 'package' | 'switch' | 'const' | 'fallthrough' | 'if' | 'range' | 'type' | 'continue' | 'for' | 'import' | 'return' | 'var'"), x("keyword")),
	i(named("identifier", "[a-zA-Z_][a-zA-Z0-9_]*"), x("identifier")), // letter { letter | unicode_digit } .   We rewrite it so to accelerate parsing
	i(named("operator", "'+=' | '&=' | '&&' | '==' | '!=' | '(' | ')' | '-=' | '|=' | '||' | '[' | ']' | '*=' | '^=' | '<-' | '>=' | '{' | '}' | '/=' | '<<=' | '<<' | '<=' | '++' | ':=' | '%=' | '>>=' | '>>' | '--'  | '...' | '&^=' | '&^' | '~' | '+' | '&' | '-' | '|' | '*' | '^' | '!' | '%' | '/' |  '=' | '>' | '<'"), x("operator")),
	i(named("punctuation", "',' | ';' | '.' | ':'"), x("punctuation")),
	o(named("literal", rules(
		o("int_lit | string_lit | rune_lit"),
		o(named("rune_lit", rules(
			o("'\\'' ( byte_value | unicode_value ) '\\''"),
			o(named("byte_value", rules(
				o("octal_byte_value | hex_byte_value"),
				i(named("octal_byte_value", "'\\\\' octal_digit{3,3}"), func(ast j.Ast, ctx *j.ParseContext) j.Ast {
					// check <= 255
					if j.NewNativeIntFrom(ast).Int() > 255 {
						return j.NewParseError(ctx, "ERROR illegal: octal value over 255")
					} else {
						return dumb{"octal_byte_value", ast, &j.Attributes{}}
					}
				}),
				i(named("hex_byte_value", "'\\\\x' hex_digit{2,2}"), x("hex_byte_value")),
			))),
			o(named("unicode_value", rules(
				o("escaped_char | little_u_value | big_u_value | unicode_char"),
				i(named("escaped_char", `[\a\f\n\r\t\v]`), x("escaped_char")), // TODO NOTE: we skip \b (BELL RING) as for some reason it doesn't work in the regex
				i(named("little_u_value", "'\\\\u' hex_digit{4,4}"), x("little_u_value")),
				i(named("big_u_value", "'\\\\U' hex_digit{8,8}"), x("big_u_value")),
			))),
			i(named("foo", "[0-9a-zA-Z]")),
		)), x("rune_lit")),
		o(named("string_lit", rules(
			o("raw_string_lit | interpreted_string_lit"),
			// o(named("raw_string_lit", "'`' ( !('`') ( unicode_char | newline ) )* '`'"), x("raw_string_lit")),                      // "`" { unicode_char | newline } "`" .   <- since unicode_char is everything but \n, it means any char.
			o(named("raw_string_lit", "/`[^`]*`/"), x("raw_string_lit")),                                                                           // "`" { unicode_char | newline } "`" .   <- since unicode_char is everything but \n, it means any char.
			o(named("interpreted_string_lit", "'\"' (!'\\\"' ('\\\\' [\\s\\S] | unicode_value | byte_value))* '\"'"), x("interpreted_string_lit")), // interpreted_string_lit = `"` { unicode_value | byte_value } `"` .
		)), x("string_lit")),
	)), x("literal")),
	i(named("int_lit", "hex_lit | octal_lit | binary_lit | decimal_lit"), x("int_lit")),
	i(named("decimal_lit", "/^0|[1-9](_?[0-9])*/"), x("decimal_lit")),
	i(named("binary_lit", "/^0[bB](_?[01])*/"), x("binary_lit")),
	i(named("octal_lit", "/^0[oO](_?[0-7])*/"), x("octal_lit")),
	i(named("hex_lit", "/^0[xX](_?[0-9a-fA-F])*/"), x("hex_lit")),
	i(named("decimal_digits", "[0-9][_0-9]*")), // decimal_digits = decimal_digit { [ "_" ] decimal_digit } .
	i(named("binary_digits", "binary_digit ('_'? binary_digit)+")),
	i(named("octal_digits", "octal_digit ('_'? octal_digit)+")),
	i(named("hex_digits", "hex_digit ('_'? hex_digit)+")),

	// "White space, formed from spaces (U+0020), horizontal tabs (U+0009),
	// carriage returns (U+000D), and newlines (U+000A), is ignored except as
	// it separates tokens that would otherwise combine into a single token."
	i(named("_", "/[ \t\n\r]*/")),
	i(named("__", "/[ \t\n\r]+/")),
)
