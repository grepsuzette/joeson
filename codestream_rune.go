package joeson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// RuneStream is a very simple code holder, cursor, matcher.
type RuneStream struct {
	text       string
	pos        int // "Hello, 世界, X" <- Pos of o is 4, Pos of 界 is 10
	lineStarts []int
}

func NewRuneStream(text string) *RuneStream {
	lineStarts := []int{0}
	for pos, rune := range text {
		if rune == '\n' {
			lineStarts = append(lineStarts, pos+1)
		}
	}
	return &RuneStream{text, 0, lineStarts}
}

func (code *RuneStream) Pos() int     { return code.pos }
func (code *RuneStream) SetPos(n int) { code.pos = n }

func (code *RuneStream) PosToLine(pos int) int { return helpers.BisectRight(code.lineStarts, pos) - 1 }
func (code *RuneStream) PosToCol(pos int) int  { return pos - code.lineStarts[code.PosToLine(pos)] }

func (code *RuneStream) Line() int   { return code.PosToLine(code.pos) }
func (code *RuneStream) Col() int    { return code.PosToCol(code.pos) }
func (code *RuneStream) Length() int { return len(code.text) }

// Get until the string `end` is encountered.
// Change code.pos accordingly, including the `end`.
func (code *RuneStream) GetUntil(end string) string { return code.GetUntilWithIgnoreEOF(end, true) }

func (code *RuneStream) GetUntilWithIgnoreEOF(end string, ignoreEOF bool) string {
	index := strings.Index(code.text[code.pos:], end)
	if index == -1 {
		if ignoreEOF {
			index = len(code.text)
		} else {
			panic("EOFError")
		}
	} else {
		index += code.pos // because we searched from this pos
		index += len(end) // what we're after is length in bytes
	}
	oldpos := code.pos
	code.pos = index
	s := helpers.SliceString(code.text, oldpos, code.pos)
	fmt.Printf("index=%d return=%s∎\n", index, s)
	return s
}

func (code *RuneStream) Peek(oper *PeekOper) string {
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
		start = code.pos - oper.beforeChars
	}
	if oper.afterLines > -1 {
		endLine := helpers.Min(len(code.lineStarts)-1, code.Line()+oper.afterLines)
		if endLine < len(code.lineStarts)-1 {
			end = code.lineStarts[endLine+1] - 1
		} else {
			end = len(code.text)
		}
	} else {
		end = code.pos + oper.afterChars
	}
	return helpers.SliceString(code.text, start, end)
}

// Match string `s` against current code.pos.
// didMatch indicates whether is succeeded
// in which case the match is in `m`
func (code *RuneStream) MatchString(s string) (didMatch bool, m string) {
	if s != helpers.SliceString(code.text, code.pos, code.pos+len(s)) {
		return false, ""
	}
	code.pos += len(s)
	return true, s
}

// Match regex `re` against current code.Pos.
// didMatch indicates whether is succeeded
// in which case the match[0] is in `m`, which may be ”
func (code *RuneStream) MatchRegexp(re regexp.Regexp) (didMatch bool, m string) {
	if firstMatchLoc := re.FindStringIndex(code.text[code.pos:]); firstMatchLoc == nil {
		return false, ""
	} else {
		if firstMatchLoc[0] != 0 {
			return false, ""
		} else {
			s := helpers.SliceString(code.text, code.pos+firstMatchLoc[0], code.pos+firstMatchLoc[1])
			code.pos += firstMatchLoc[1]
			return true, s
		}
	}
}

func (code *RuneStream) Print() string {
	s := "Code at offset " + BoldYellow(strconv.Itoa(code.pos)) + "/" + BoldYellow(strconv.Itoa(len(code.text))) + ": '"
	s += Cyan(helpers.SliceString(code.text, helpers.Max(0, code.pos-20), code.pos))
	s += BoldCyan("|")
	s += BoldWhite(helpers.SliceString(code.text, code.pos, code.pos+40)) + "'"
	return s
}
