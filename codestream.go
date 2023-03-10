package joeson

import (
	"fmt"
	"github.com/grepsuzette/joeson/helpers"
	"regexp"
	"strconv"
	"strings"
)

type Cursor struct {
	line int
	col  int
	pos  int
}

type Origin struct {
	code  string
	start int
	end   int
}

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
	Pos        int // "Hello, δΈη, X" <- pos of o is 4, pos of η is 10
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
func (code *CodeStream) PosToCursor(pos int) Cursor {
	line := helpers.BisectRight(code.lineStarts, pos) - 1
	return Cursor{line: line, pos: pos}
}

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
	fmt.Printf("index=%d return=%sβ\n", index, s)
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

// Get next byte(s). Default value for len is 1,
// this is why its put as a pointer here (nil will give a value of 1)
func (code *CodeStream) Next(pLen *int) string {
	n := 1
	if pLen != nil {
		n = *pLen
	}
	if n <= 0 {
		panic("<CodeStream>.next wants len > 0")
	}
	oldpos := code.Pos
	code.Pos += n
	return helpers.SliceString(code.text, oldpos, code.Pos)
}

// Get next hex byte(s) as number.
// Default value for len is 1,
// this is why its put as a pointer here (nil will give a value of 1)
// If you use more than 8, note it will outflow the capacity of uint64.
func (code *CodeStream) Hex(pLen *int) uint64 {
	// Do we want int, or int64 etc?
	// We read 4 bits at a time, hence the <<4
	// Ultimately we can handle 64bits.
	// Let's return uint64
	var num uint64 = 0
	nextBytes := code.Next(pLen)
	len1 := len(nextBytes)
	for i := 0; i < len1; i++ {
		theByte := nextBytes[i]
		// "The bitSize argument [3rd one of ParseUint] specifies the integer
		// type that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
		// correspond to int, int8, int16, int32, and int64. If bitSize is
		// below 0 or above 64, an error is returned."
		//  -> As we read 4bit at a time, doesn't really matter
		//     and we can use uint8 here -----------------v
		if theUint, err := strconv.ParseUint(string(theByte), 16, 8); err != nil {
			panic("Invalid hex-character pattern in string")
		} else {
			num = (num << 4) | theUint
		}
	}
	return num
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
// in which case the match[0] is in `m`, which may be ''
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
	s := "Code at offset " + boldYellow(strconv.Itoa(code.Pos)) + "/" + boldYellow(strconv.Itoa(len(code.text))) + ": '"
	s += cyan(helpers.SliceString(code.text, helpers.Max(0, code.Pos-20), code.Pos))
	s += boldCyan("|")
	s += boldWhite(helpers.SliceString(code.text, code.Pos, code.Pos+40)) + "'"
	// s += "' (note: caret | and colors extraneously inserted)"
	return s
}
