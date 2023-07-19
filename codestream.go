package joeson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// E.g. NewPeek().BeforeLines(2).AfterLines(4)
type peekOper struct {
	beforeChars int // they all use -1 for unspecified
	beforeLines int
	afterChars  int
	afterLines  int
}

func NewPeek() *peekOper {
	return &peekOper{beforeChars: -1, beforeLines: -1, afterChars: -1, afterLines: -1}
}
func (ps *peekOper) BeforeChars(n int) *peekOper { ps.beforeChars = n; return ps }
func (ps *peekOper) AfterChars(n int) *peekOper  { ps.afterChars = n; return ps }
func (ps *peekOper) BeforeLines(n int) *peekOper { ps.beforeLines = n; return ps }
func (ps *peekOper) AfterLines(n int) *peekOper  { ps.afterLines = n; return ps }

// Pos acts as a cursor
type CodeStream struct {
	text       string
	Pos        int // "Hello, 世界, X" <- Pos of o is 4, Pos of 界 is 10
	lineStarts []int
}

func NewCodeStream(text string) *CodeStream {
	lineStarts := []int{0}
	for pos, rune := range text {
		if rune == '\n' {
			lineStarts = append(lineStarts, pos+1)
		}
	}
	return &CodeStream{text, 0, lineStarts}
}

func (code *CodeStream) PosToLine(pos int) int { return helpers.BisectRight(code.lineStarts, pos) - 1 }
func (code *CodeStream) PosToCol(pos int) int  { return pos - code.lineStarts[code.PosToLine(pos)] }

func (code *CodeStream) Line() int   { return code.PosToLine(code.Pos) }
func (code *CodeStream) Col() int    { return code.PosToCol(code.Pos) }
func (code *CodeStream) Length() int { return len(code.text) }

// Get until the string `end` is encountered.
// Change code.pos accordingly, including the `end`.
func (code *CodeStream) GetUntil(end string) string { return code.GetUntilWithIgnoreEOF(end, true) }

func (code *CodeStream) GetUntilWithIgnoreEOF(end string, ignoreEOF bool) string {
	index := strings.Index(code.text[code.Pos:], end)
	if index == -1 {
		if ignoreEOF {
			index = len(code.text)
		} else {
			panic("EOFError")
		}
	} else {
		index += code.Pos // because we searched from this pos
		index += len(end) // what we're after is length in bytes
	}
	oldpos := code.Pos
	code.Pos = index
	s := helpers.SliceString(code.text, oldpos, code.Pos)
	fmt.Printf("index=%d return=%s∎\n", index, s)
	return s
}

func (code *CodeStream) Peek(oper *peekOper) string {
	if oper.beforeLines < 0 && oper.beforeChars < 0 {
		oper.beforeChars = 0
	}
	if oper.afterLines < 0 && oper.afterChars < 0 {
		oper.afterChars = 0
	}
	if oper.beforeChars == 0 && oper.afterChars == 0 {
		return ""
	}
	start := 0
	end := 0
	if oper.beforeLines > -1 {
		startLine := helpers.Max(0, code.Line()-oper.beforeLines)
		start = code.lineStarts[startLine]
	} else {
		start = code.Pos - oper.beforeChars
	}
	if oper.afterLines > -1 {
		endLine := helpers.Min(len(code.lineStarts)-1, code.Line()+oper.afterLines)
		if endLine < len(code.lineStarts)-1 {
			end = code.lineStarts[endLine+1] - 1
		} else {
			end = len(code.text)
		}
	} else {
		end = code.Pos + oper.afterChars
	}
	return helpers.SliceString(code.text, start, end)
}

// Match string `s` against current code.pos.
// didMatch indicates whether is succeeded
// in which case the match is in `m`
func (code *CodeStream) MatchString(s string) (didMatch bool, m string) {
	if s != helpers.SliceString(code.text, code.Pos, code.Pos+len(s)) {
		return false, ""
	}
	code.Pos += len(s)
	return true, s
}

// Match regex `re` against current code.Pos.
// didMatch indicates whether is succeeded
// in which case the match[0] is in `m`, which may be ”
func (code *CodeStream) MatchRegexp(re regexp.Regexp) (didMatch bool, m string) {
	if firstMatchLoc := re.FindStringIndex(code.text[code.Pos:]); firstMatchLoc == nil {
		return false, ""
	} else {
		if firstMatchLoc[0] != 0 {
			return false, ""
		} else {
			s := helpers.SliceString(code.text, code.Pos+firstMatchLoc[0], code.Pos+firstMatchLoc[1])
			code.Pos += firstMatchLoc[1]
			return true, s
		}
	}
}

func (code *CodeStream) Print() string {
	s := "Code at offset " + BoldYellow(strconv.Itoa(code.Pos)) + "/" + BoldYellow(strconv.Itoa(len(code.text))) + ": '"
	s += Cyan(helpers.SliceString(code.text, helpers.Max(0, code.Pos-20), code.Pos))
	s += BoldCyan("|")
	s += BoldWhite(helpers.SliceString(code.text, code.Pos, code.Pos+40)) + "'"
	return s
}
