package main

import j "github.com/grepsuzette/joeson"

// https://go.dev/ref/spec#Characters
// https://go.dev/ref/spec#Letters_and_digits
//
// The following terms are used to denote specific Unicode character categories:
//
// newline        = /* the Unicode code point U+000A */ .
// unicode_char   = /* an arbitrary Unicode code point except newline */ .
// unicode_letter = /* a Unicode code point categorized as "Letter" */ .
// unicode_digit  = /* a Unicode code point categorized as "Number, decimal digit" */ .
//
// In The Unicode Standard 8.0, Section 4.5 "General Category" defines a set of character categories. Go treats all characters in any of the Letter categories Lu, Ll, Lt, Lm, or Lo as Unicode letters, and those in the Number category Nd as Unicode digits.
// Letters and digits
//
// The underscore character _ (U+005F) is considered a lowercase letter.
//
// letter        = unicode_letter | "_" .
// decimal_digit = "0" … "9" .
// binary_digit  = "0" | "1" .
// octal_digit   = "0" … "7" .
// hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .

var rules_chars = rules(
	o(named("expr", "(characters | letter | digits)+")),

	/* Characters. https://go.dev/ref/spec#string_lit:
	"In The Unicode Standard 8.0, Section 4.5 "General Category" defines a set of
	character categories. Go treats all characters in any of the Letter categories
	Lu, Ll, Lt, Lm, or Lo as Unicode letters, and those in the Number category Nd
	as Unicode digits.""
	*/
	i(named("characters", "(newline | unicode_char | unicode_letter | unicode_digit)")),
	i(named("newline", "'\n'")),                                 // "the Unicode code point U+000A"
	i(named("unicode_char", "[^\\x{0a}]"), x("unicode_char")),   // "an arbitrary Unicode code point except newline"
	i(named("unicode_letter", "[a-zA-Z]"), x("unicode_letter")), // "a Unicode code point categorized as "Letter""
	i(named("unicode_digit", "[0-9]"), x("unicode_digit")),      // "a Unicode code point categorized as "Number, decimal digit""
	//                         ^^^
	// For now we'll stick to ANSI for letters and digits. It can later be improved
	// That's because https://www.unicode.org/versions/Unicode8.0.0/ch04.pdf <- Section 4.5
	//   does not define them.

	// Letters and digits
	i(named("letter", "unicode_letter | '_'")), // "The underscore character _ (U+005F) is considered a lowercase letter."
	i(named("digits", "decimal_digit | binary_digit | octal_digit | hex_digit")),
	i(named("decimal_digit", "[0-9]"), j.ParseOptions{Debug:true}),
	i(named("binary_digit", "[01]")),
	i(named("octal_digit", "[0-7]")),
	i(named("hex_digit", "[0-9A-Fa-f]")),
)
