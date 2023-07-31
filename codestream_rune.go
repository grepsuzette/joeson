package joeson

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// A simple code holder, cursor, matcher implementing CodeStream.
// The name "stream" is a little bit illusory. It's a string okay.
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
	// fmt.Printf("index=%d return=%s∎\n", index, s)
	return s
}

// take a look n runes before or after, don't update position
func (code *RuneStream) PeekRunes(n int) string {
	start := code.pos
	end := code.pos
	if n < 0 {
		start += n
	} else {
		end += n
	}
	return helpers.SliceString(code.text, start, end)
}

// take a look n lines before or after, don't update position
func (code *RuneStream) PeekLines(n int) string {
	start := code.pos
	end := code.pos
	if n < 0 {
		start = code.lineStarts[helpers.Max(0, code.Line()+n)]
	} else {
		endLine := helpers.Min(len(code.lineStarts)-1, code.Line()+n)
		if endLine < len(code.lineStarts)-1 {
			end = code.lineStarts[endLine+1] - 1
		} else {
			end = len(code.text)
		}
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

// Match regex `re` against current position.
// didMatch indicates whether is succeeded.
// If so the full text for the match (usually called match[0]) is in m.
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

// debugging purposes only
func (code *RuneStream) Print() string {
	s := "Code at offset " + BoldYellow(strconv.Itoa(code.pos)) + "/" + BoldYellow(strconv.Itoa(len(code.text))) + ": '"
	s += Cyan(helpers.SliceString(code.text, helpers.Max(0, code.pos-20), code.pos))
	s += BoldCyan("|")
	s += BoldWhite(helpers.SliceString(code.text, code.pos, code.pos+40)) + "'"
	return s
}
