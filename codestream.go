package joeson

import (
	"regexp"
)

// A code holder, cursor, matcher.
type CodeStream interface {
	Pos() int   // the current offset
	SetPos(int) // panics when n is out of bounds
	PosToLine(pos int) int
	PosToCol(pos int) int
	Line() int    // Current line. First line is 0.
	Col() int     // Current column. First column is 1.
	Code() string // User-provided code
	Length() int  // len(Code())

	GetUntil(end string) string // Get until the string `end` is encountered.  Change current position accordingly, including the string
	GetUntilWithIgnoreEOF(end string, ignoreEOF bool) string
	PeekRunes(n int) string    // TODO multiargs like PeekLines  // e.g. -3 to peek 3 runes back. 2 to peek 2 runes forward. Does not change position.
	PeekLines(n ...int) string // Peeks a range of lines relative to current one. A range is built from ...int. Empty peeks current line. When 1 parameter is given a second 0 is implied. E.g. PeekLines(-1, 2), 2 to peek 2 lines forward. -1 to peek 1 line backwards. You may pass more args, the min and max of the series will be used.
	MatchString(string) (didMatch bool, m string)
	MatchRegexp(regexp.Regexp) (didMatch bool, m string)
	MatchRune(func(rune) bool) (didMatch bool, m rune)

	Print() string      // short line info, can be integrated to longer parse error messages
	PrintDebug() string // free-form multiline detailed debug information

	// usable only by grammar, parsectx and packrat
	workLength() int
}

var (
	_ CodeStream = &RuneStream{}
	_ CodeStream = &TokenStream{}
)
