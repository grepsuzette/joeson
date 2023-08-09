package joeson

import (
	"fmt"
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

func NewRuneStream(text string) CodeStream {
	lineStarts := []int{0}
	for pos, rune := range text {
		if rune == '\n' {
			lineStarts = append(lineStarts, pos+1)
		}
	}
	return &RuneStream{text, 0, lineStarts}
}

func (code *RuneStream) Pos() int { return code.pos }
func (code *RuneStream) SetPos(n int) {
	// n == len(code.text) is allowed for now, for EOF
	if n < 0 || n > len(code.text) {
		panic(fmt.Sprintf("%d is out of bound", n))
	}
	code.pos = n
}

func (code *RuneStream) PosToLine(pos int) int { return helpers.BisectRight(code.lineStarts, pos) - 1 }
func (code *RuneStream) PosToCol(pos int) int  { return pos - code.lineStarts[code.PosToLine(pos)] }

func (code *RuneStream) Line() int       { return code.PosToLine(code.pos) }
func (code *RuneStream) Col() int        { return code.PosToCol(code.pos) }
func (code *RuneStream) Code() string    { return code.text }
func (code *RuneStream) Length() int     { return len(code.text) }
func (code *RuneStream) workLength() int { return len(code.text) }

// Get until the string `end` is encountered.
// Change code.pos accordingly, including the `end`.
func (code *RuneStream) GetUntil(end string) string {
	offset := strings.Index(code.text[code.pos:], end)
	if offset == -1 {
		offset = len(code.text)
	} else {
		offset += code.pos // because we searched from this pos
		offset += len(end) // what we're after is length in bytes
	}
	oldpos := code.pos
	code.pos = offset
	return code.text[oldpos:offset]
}

// take a look n runes before or after, don't update position
func (code *RuneStream) PeekRunes(n int) string {
	if n <= 0 {
		return helpers.LastNRunes(code.text[:code.pos], -n)
	} else {
		var b strings.Builder
		i := 0
		for _, rune := range code.text[code.pos:] {
			b.WriteRune(rune)
			i++
			if i >= n {
				break
			}
		}
		return b.String()
	}
}

// Take a look n lines before or after, don't update position
// Negative means to look n lines before, positive means after.
// It is possible to provide 2 or more arguments.
// In that case, it will peek from the minimum to the maximum of the series.
// When only 1 value is given, a second value of 0 is implied to create a range.
func (code *RuneStream) PeekLines(n ...int) string {
	if len(n) <= 0 {
		return ""
	} else if len(n) == 1 {
		n = []int{n[0], 0} // implied 0
	}
	min := n[0]
	max := n[0]
	for _, n := range n {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	start := code.lineStarts[helpers.Max(0, code.Line()+min)]
	var end int
	endLine := helpers.Min(len(code.lineStarts)-1, code.Line()+max)
	if endLine < len(code.lineStarts)-1 {
		end = code.lineStarts[endLine+1] - 1
	} else {
		end = len(code.text)
	}
	return helpers.SliceString(code.text, start, end)
}

// Match string `s` against current code.pos.
// didMatch indicates whether is succeeded
// in which case the match is in `m`
func (code *RuneStream) MatchString(s string) (didMatch bool, m string) {
	if s != code.text[code.pos:helpers.Min(code.pos+len(s), len(code.text))] {
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
			from := code.pos + firstMatchLoc[0]
			to := helpers.Min(code.pos+firstMatchLoc[1], len(code.text))
			s := code.text[from:to]
			code.pos += firstMatchLoc[1]
			return true, s
		}
	}
}

// Match func(rune) bool against rune at current position.
// didMatch indicates whether is succeeded. If so the rune is m and position is
// updated. When at EOF it never match.
func (code *RuneStream) MatchRune(f func(rune rune) bool) (didMatch bool, m rune) {
	if code.pos >= len(code.text) {
		return false, '\x00' // never match at EOF
	}
	var ret rune
	newPos := code.pos
	iter := 0
	for offset, rune := range code.text[code.pos:] {
		if iter == 1 {
			newPos += offset // before leaving add offset of the next character
			break
		}
		if !f(rune) {
			return false, ' '
		} else {
			ret = rune
			iter++ // another round to take offset of the next rune and immediately break
		}
	}
	if newPos == code.pos {
		// when not updated, it means rune matched was the last in text
		code.SetPos(len(code.text))
	} else {
		code.SetPos(newPos)
	}
	return true, ret
}

// short, single line information to be integrated in parse errors
func (code *RuneStream) Print() string {
	s := "Code at offset " + BoldYellow(strconv.Itoa(code.pos)) + "/" + BoldYellow(strconv.Itoa(len(code.text))) + ": '"
	s += Cyan(helpers.SliceString(code.text, helpers.Max(0, code.pos-20), code.pos))
	s += BoldCyan("|")
	s += BoldWhite(helpers.SliceString(code.text, code.pos, code.pos+40)) + "'"
	return s
}

func (code *RuneStream) PrintDebug() string {
	return code.Print()
}
