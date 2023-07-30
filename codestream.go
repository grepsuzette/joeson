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
	Peek(*PeekOper) string
	MatchString(string) (didMatch bool, m string)
	MatchRegexp(regexp.Regexp) (didMatch bool, m string)
	Print() string
}

var _ CodeStream = &RuneStream{} // _ CodeStream = &TokenStream{}

// E.g. NewPeek().BeforeLines(2).AfterLines(4)
type PeekOper struct {
	beforeChars int // they all use -1 for unspecified
	beforeLines int
	afterChars  int
	afterLines  int
}

func NewPeek() *PeekOper {
	return &PeekOper{beforeChars: -1, beforeLines: -1, afterChars: -1, afterLines: -1}
}
func (ps *PeekOper) BeforeChars(n int) *PeekOper { ps.beforeChars = n; return ps }
func (ps *PeekOper) AfterChars(n int) *PeekOper  { ps.afterChars = n; return ps }
func (ps *PeekOper) BeforeLines(n int) *PeekOper { ps.beforeLines = n; return ps }
func (ps *PeekOper) AfterLines(n int) *PeekOper  { ps.afterLines = n; return ps }
