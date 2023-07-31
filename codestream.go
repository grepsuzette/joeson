package joeson

import (
	"regexp"
)

// CodeStream is a code holder, cursor, matcher.
//
// Implementations:
// - RuneStream is a simple implementation.
// - TokenStream allows to work with pre-tokenized source code.
type CodeStream interface {
	// `Pos` means the offset in the tokenized file
	// for non-tokenized stream, there is of course no such ambiguity.
	Pos() int
	SetPos(int)
	PosToLine(pos int) int
	PosToCol(pos int) int
	Line() int // first line is 1
	Col() int  // first column is 1
	Length() int

	// all relating to the working text
	GetUntil(end string) string // Get until the string `end` is encountered.  Change code.pos accordingly, including the string
	GetUntilWithIgnoreEOF(end string, ignoreEOF bool) string
	PeekRunes(n int) string // e.g. -3 to go back 3 runes. 2 to advance 2 runes
	PeekLines(n int) string // e.g. 2 to advance 2 lines. 1 to advance 1 line (this one is not necessarily precise, meant for printing purposes)
	MatchString(string) (didMatch bool, m string)
	MatchRegexp(regexp.Regexp) (didMatch bool, m string)
	Print() string
}

var _ CodeStream = &RuneStream{} // _ CodeStream = &TokenStream{}
