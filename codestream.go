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
	// important: for Pos, SetPos, PosToLine, PosToCol, if source is tokenized,
	// pos is relative to the transformed text (that is the working text).
	Pos() int
	SetPos(int)
	PosToLine(pos int) int
	PosToCol(pos int) int

	Line() int // line in the original text
	Col() int  // col in the original text

	// all relating to the working text
	Length() int
	GetUntil(end string) string
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
