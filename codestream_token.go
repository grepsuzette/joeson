package joeson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// The idea of TokenStream is to match against a tokenkized grammar
// that would standardize for example blanks or comments, whilst converting
// back error locations when there is a problem.
//
// There are therefore two views of the same thing:
//
// - Some user-provided `original` text, to be parsed.
// - Some user-provided tokens from that same text. A `work` string
// is joined from those tokens, along with the positions of each token in
// the original text.
//
// Two systems of source coordinates exist then.
// An illustration is provided to explain graphically:
//
//	 pos=0    pos=12
//	/        /
//	=======  ===N===== <--- tokens
//	====  ======= ===
//	===========  =====
//
// When fail at offset N,
// we can find which offset
// of which token, and
// translate back to original
// source space position (pos).
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

func (ts *TokenStream) Pos() int {
	return ts.workOffset
}

func (ts *TokenStream) SetPos(n int) {
	if n > len(ts.work) {
		panic("assert")
	}
	ts.workOffset = n
}

func (code *TokenStream) PosToLine(workOffset int) int {
	return helpers.BisectRight(
		code.lineStarts,
		code.coords(workOffset).originalOffset,
	) - 1
}

func (code *TokenStream) PosToCol(workOffset int) int {
	return code.coords(workOffset).originalOffset -
		code.lineStarts[code.PosToLine(workOffset)]
}

func (code *TokenStream) Line() int   { return code.PosToLine(code.workOffset) }
func (code *TokenStream) Col() int    { return code.PosToCol(code.workOffset) }
func (code *TokenStream) Length() int { return len(code.original) }

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

func (code *TokenStream) PeekRunes(n int) string {
	panic("TODO")
}

func (code *TokenStream) PeekLines(n int) string {
	panic("TODO")
}

// the reverse operation can be obtained with coords()
// func (code *TokenStream) calcWorkOffset(originalOffset int) int {
// 	for n, token := range code.tokens {
// 		token
// 	}
// }

// Match string `s` against current code.pos.
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

func (code *TokenStream) MatchRegexp(re regexp.Regexp) (didMatch bool, m string) { panic("TODO") }

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
	s += "Work text (tokenized):\n"
	s += code.work + "\n\n"
	s += "Tokens:\n"
	for _, token := range code.tokens {
		s += BoldBlack("[o=" + strconv.Itoa(token.OriginalOffset) + ", w=" + strconv.Itoa(token.WorkOffset) + "]")
		s += token.Repr
	}
	return s
}

// if overflow, panic since it would be an impossible coordinate
// the reverse operation can be obtain with calcWorkOffset()
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
