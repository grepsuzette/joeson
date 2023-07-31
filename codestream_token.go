package joeson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// TokenStream allows matching against tokenkized texts.
// User can provide tokens produced from an original text.
// Two systems of source coordinates exist then (token-space, original-space).
// An illustration is provided to explain graphically,
// although it's relatively straightforward:
//
//	 pos=0    pos=12
//	/        /
//	=======  ===N===== <--- tokens
//	====  ======= ===
//	===========  =====
//
// When fail at offset N, we can find which offset
// of which token, and translate back to original
// source space position (pos).
//
// Grammars are intended to parse against the tokenized text (`work`)
// and errors are meant to show the `original` text.
type (
	TokenStream struct {
		tokens     []Token // tokens with their position
		original   string  // original text (not tokenized)
		work       string  // work text (tokenized + joined upon " ")
		workOffset int     // the current working position relative to `work`
		lineStarts []int   // line starts in the `original` text
	}
	Token struct {
		Repr           string
		OriginalOffset int // relative to original
		WorkOffset     int // relative to work
		// meta interface{} // for now is useless
	}
	coord struct {
		token          Token // token found at requested position
		nToken         int   // number that token, ∈ [0; len(tokens)[
		offsetInToken  int   // offset of requested position in token.Repr, ∈ [0; len(tokens[nToken].Repr)[
		workOffset     int   // offset relative to `work` text. ∈ [0; len(work)[ . Use toWorkOffset()
		originalOffset int   // offset relative to `original` text
	}
)

// `tokens` must have been generated from `text` using a scanner/lexer/tokenizer,
// however you want to call it. An example is "go/scanner".
// If `tokens` have position starting beyond `text` it will panic
func NewTokenStream(text string, tokens []Token) *TokenStream {
	lineStarts := []int{0}
	for pos, rune := range text {
		if rune == '\n' {
			lineStarts = append(lineStarts, pos+1)
		}
	}
	var b strings.Builder
	for _, token := range tokens {
		if token.OriginalOffset > len(text) {
			panic("tokens reference text outside the original text")
		}
		b.WriteString(token.Repr)
	}
	return &TokenStream{tokens, text, b.String(), 0, lineStarts}
}

// `Pos` here means the offset in the tokenized string (AKA workOffset)
func (ts *TokenStream) Pos() int {
	return ts.workOffset
}

// `Pos` here means the offset in the tokenized string (AKA workOffset)
func (ts *TokenStream) SetPos(n int) {
	if n > len(ts.work) {
		panic("assert")
	}
	ts.workOffset = n
}

// `Pos` here means the offset in the tokenized string (AKA workOffset)
// Line refers to the original text, and start at 1.
func (code *TokenStream) PosToLine(workOffset int) int {
	return helpers.BisectRight(
		code.lineStarts,
		code.coords(workOffset).originalOffset,
	) - 1
}

// `Pos` here means the offset in the tokenized string (AKA workOffset)
// Col refers to the original text, and start at 1.
func (code *TokenStream) PosToCol(workOffset int) int {
	return code.coords(workOffset).originalOffset -
		code.lineStarts[code.PosToLine(workOffset)]
}

// Current line (in the original text), starting counting at 1.
func (code *TokenStream) Line() int { return code.PosToLine(code.workOffset) }

// Current column (in the original text), starting counting at 1.
func (code *TokenStream) Col() int { return code.PosToCol(code.workOffset) }

<<<<<<< HEAD
// Length of the original text (since exported function are for external usage)
func (code *TokenStream) Length() int { return len(code.original) }
=======
func (code *TokenStream) Code() string { return code.original }
func (code *TokenStream) Length() int  { return len(code.original) }

func (code *TokenStream) workLength() int { return len(code.work) }
>>>>>>> d68d265 (feat: add grammar ParseTokens)

// Get until the string `end` is encountered.
// Change workingpos accordingly, including the string
func (code *TokenStream) GetUntil(end string) string {
	return code.GetUntilWithIgnoreEOF(end, true)
}

// Get until the string `end` is encountered.
// Change workingpos accordingly, including the string
func (code *TokenStream) GetUntilWithIgnoreEOF(end string, ignoreEOF bool) string {
	index := strings.Index(code.work[code.workOffset:], end)
	if index == -1 {
		if ignoreEOF {
			index = len(code.work)
		} else {
			panic("EOFError")
		}
	} else {
		index += code.workOffset // because we searched from this pos
		index += len(end)        // what we're after is length in bytes
	}
	oldWorkOffset := code.workOffset
	code.workOffset = index
	s := helpers.SliceString(code.work, oldWorkOffset, code.workOffset)
	return s
}

// Take a look `n` runes backwards or forwards, depending on the sign of n,
// return the string contained in the interval made with the current position.
func (code *TokenStream) PeekRunes(n int) string {
	start := code.workOffset
	end := code.workOffset
	if n < 0 {
		start += n
	} else {
		end += n
	}
	return helpers.SliceString(code.work, start, end)
}

// Take a look `n` lines backwards or forwards, depending on the sign of n,
// return the string contained in the interval made with the current position.
// PeekLines() is meant for printing purposes only.
// It responds with the original text, not the tokenized one.
func (code *TokenStream) PeekLines(n int) string {
	pos := code.coords(code.workOffset).originalOffset
	start := pos
	end := pos
	if n < 0 {
		start = code.lineStarts[helpers.Max(0, code.Line()+n)]
	} else {
		endLine := helpers.Min(len(code.lineStarts)-1, code.Line()+n)
		if endLine < len(code.lineStarts)-1 {
			end = code.lineStarts[endLine+1] - 1
		} else {
			end = len(code.original)
		}
	}
	return helpers.SliceString(code.original, start, end) // respond w/ original text
}

// Match string `s` against current position.
// didMatch indicates whether is succeeded
// in which case the match is in `m`
func (code *TokenStream) MatchString(s string) (didMatch bool, m string) {
	if s != helpers.SliceString(code.work, code.workOffset, code.workOffset+len(s)) {
		return false, ""
	} else {
	}
	code.workOffset += len(s)
	return true, s
}

// Match regex `re` against current position.
// didMatch indicates whether is succeeded.
// If so the full text for the match (usually called match[0]) is in m.
func (code *TokenStream) MatchRegexp(re regexp.Regexp) (didMatch bool, m string) {
	if firstMatchLoc := re.FindStringIndex(code.work[code.workOffset:]); firstMatchLoc == nil {
		return false, ""
	} else {
		if firstMatchLoc[0] != 0 {
			return false, ""
		} else {
			s := helpers.SliceString(code.work, code.workOffset+firstMatchLoc[0], code.workOffset+firstMatchLoc[1])
			code.workOffset += firstMatchLoc[1]
			return true, s
		}
	}
}

// This somewhat lengthy printage is really meant for debugging/developping tests.  Don't use it for anything else, or prepare to meet your destiny.
func (code *TokenStream) Print() string {
	pos := code.workOffset
	s := "Code at offset " + BoldYellow(strconv.Itoa(pos)) + "/" + BoldYellow(strconv.Itoa(len(code.original))) + ": '"
	s += Cyan(helpers.SliceString(code.original, helpers.Max(0, pos-20), pos))
	s += BoldCyan("|")
	s += BoldWhite(helpers.SliceString(code.original, pos, pos+40)) + "'"
	s += "\n"
	s += "lineStarts:\n"
	s += fmt.Sprintf("%v\n", code.lineStarts)
	s += "Original text:\n"
	s += strings.NewReplacer("\n", "<CR>\n", "\t", "<TAB>", " ", "_").Replace(code.original) + "\n\n"
	s += code.PrintWorkText()
	s += "\n"
	s += "Tokens:\n"
	for _, token := range code.tokens {
		s += BoldBlack("[o=" + strconv.Itoa(token.OriginalOffset) + ", w=" + strconv.Itoa(token.WorkOffset) + "]")
		s += token.Repr
	}
	return s
}

<<<<<<< HEAD
=======
func (code *TokenStream) PrintWorkText() string {
	return "Work text (tokenized):\n" + code.work + "\n"
}

>>>>>>> d68d265 (feat: add grammar ParseTokens)
// Given an arbitrary work offset (as given by Pos()),
// get all possible coordinates (i.e. originalOffset, line, col).
//
// if overflow, panic since it would be an impossible coordinate.
// the reverse operation can be obtain with calcWorkOffset().
func (code *TokenStream) coords(workOffset int) coord {
	// find most advanced token number, such that the following token would begin
	// after workOffset.
	nToken := 0
	var token Token
	for {
		if nToken >= len(code.tokens) {
			panic(fmt.Sprintf("workOffset %d would overflow, what you ask makes no sense", workOffset))
		}
		token = code.tokens[nToken]
		if workOffset < token.WorkOffset+len(token.Repr) {
			break
		}
		nToken++
	}
	offsetInToken := workOffset - token.WorkOffset
	originalOffset := code.tokens[nToken].OriginalOffset + offsetInToken
	return coord{
		token:          token,
		nToken:         nToken,
		offsetInToken:  offsetInToken,
		workOffset:     workOffset,
		originalOffset: originalOffset,
	}
}
