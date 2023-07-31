package joeson

import (
	"regexp"
)

// A code holder, cursor, matcher.
type CodeStream interface {
	Pos() int
	SetPos(int)
	PosToLine(pos int) int
	PosToCol(pos int) int
	Line() int    // Current line. First line is 1.
	Col() int     // Current column. First column is 1.
	Code() string // User-provided code
	Length() int  // len(Code())

	GetUntil(end string) string // Get until the string `end` is encountered.  Change current position accordingly, including the string
	GetUntilWithIgnoreEOF(end string, ignoreEOF bool) string
	PeekRunes(n int) string // e.g. -3 to peek 3 runes back. 2 to peek 2 runes forward. Does not change position.
	PeekLines(n int) string // e.g. 2 to peek 2 lines forward. -1 to peek 1 line backwards. (this one is not necessarily precise, it is meant for printing purposes and not for parsing)
	MatchString(string) (didMatch bool, m string)
	MatchRegexp(regexp.Regexp) (didMatch bool, m string)
	Print() string

	// usable only by grammar, parsectx and packrat
	workLength() int
}

var (
	_ CodeStream = &RuneStream{}
	_ CodeStream = &TokenStream{}
)
