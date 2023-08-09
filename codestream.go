package joeson

import (
	"regexp"
)

// A code holder, cursor, matcher.
type CodeStream interface {
	Pos() int   // The current offset
	SetPos(int) // Panics when n is out of bounds
	PosToLine(pos int) int
	PosToCol(pos int) int
	Line() int    // Current line. First line is 0.
	Col() int     // Current column. First column is 1.
	Code() string // User-provided code
	Length() int  // Like len(Code())

	GetUntil(end string) string                          // Move position until the string `end` and return the string from Pos() to there. Return "" if unfound, not changing position in that case.
	PeekRunes(n int) string                              // Does not change position. TODO multiargs like PeekLines  // e.g. -3 to peek 3 runes back. 2 to peek 2 runes forward.
	PeekLines(n ...int) string                           // Does not change position. Peeks a range of lines relative to current one. A range is built from ...int. Empty peeks current line. When 1 parameter is given a second 0 is implied. E.g. PeekLines(-1, 2), 2 to peek 2 lines forward. -1 to peek 1 line backwards. You may pass more args, the min and max of the series will be used.
	MatchString(string) (didMatch bool, m string)        // Advance position if matched
	MatchRegexp(regexp.Regexp) (didMatch bool, m string) // Advance position if matched
	MatchRune(func(rune) bool) (didMatch bool, m rune)   // Advance position if matched

	Print() string      // short line info, can be integrated to longer parse error messages
	PrintDebug() string // free-form multiline detailed debug information

	// usable only by grammar, parsectx and packrat
	workLength() int
}

var (
	_ CodeStream = &RuneStream{}
	_ CodeStream = &TokenStream{}
)
