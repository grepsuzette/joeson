package core

import (
	"fmt"
	"grepsuzette/joeson/helpers"
	"regexp"
	"strconv"
	"strings"
)

// surely NullInt is more usable than *int...
// https://stackoverflow.com/questions/68800319/how-to-differentiate-int-null-and-defaulted-to-zero-from-int-actually-equal-to-z
// let's see for now
type Peek struct {
	BeforeChars helpers.NullInt
	BeforeLines helpers.NullInt
	AfterChars  helpers.NullInt
	AfterLines  helpers.NullInt
}

// Pos acts as a cursor
type CodeStream struct {
	text       string
	Pos        int // "Hello, 世界, X" <- pos of o is 4, pos of 界 is 10
	lineStarts []int
}

func NewCodeStream(text string) *CodeStream {
	lineStarts := []int{0}
	for pos, rune := range text {
		if rune == '\n' {
			lineStarts = append(lineStarts, pos+1)
		}
	}
	for _, v := range lineStarts {
		fmt.Println(v)
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
	fmt.Printf("index=%d return=%s∎\n", index, code.text[oldpos:code.Pos])
	return code.text[oldpos:code.Pos]
}

func (code *CodeStream) Peek(o Peek) string {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	if !o.BeforeLines.IsSet && !o.BeforeChars.IsSet {
		o.BeforeChars = helpers.NewNullInt(0)
	}
	if !o.AfterLines.IsSet && !o.AfterChars.IsSet {
		o.AfterChars = helpers.NewNullInt(0)
	}
	if o.BeforeChars.IsSet && o.BeforeChars.Int == 0 && o.AfterChars.IsSet && o.AfterChars.Int == 0 {
		return ""
	}
	start := 0
	end := 0
	if o.BeforeLines.IsSet {
		startLine := max(0, code.Line()-o.BeforeLines.Int)
		start = code.lineStarts[startLine]
	} else {
		start = code.Pos - o.BeforeChars.Int
	}
	if o.AfterLines.IsSet {
		endLine := min(len(code.lineStarts)-1, code.Line()+o.AfterLines.Int)
		if endLine < len(code.lineStarts)-1 {
			end = code.lineStarts[endLine+1] - 1
		} else {
			end = len(code.text)
		}
	} else {
		end = code.Pos + o.AfterChars.Int
	}
	return code.text[start:end]
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
	return code.text[oldpos:code.Pos]
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
	if s != code.text[code.Pos:code.Pos+len(s)] {
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
		if firstMatchLoc[0] != code.Pos {
			return false, ""
		} else {
			// TODO test it!
			// original code is
			/*
			   # Regex returns null if match failed,
			   # otherwise returns match[0] which may be ''
			   regex.lastIndex = @pos
			   match = regex.exec(@text)
			   return null if not match or match.index != @pos
			   @pos = regex.lastIndex
			   return match[0]
			*/
			code.Pos = firstMatchLoc[1]
			return true, code.text[firstMatchLoc[0]:firstMatchLoc[1]]
		}
	}
}
