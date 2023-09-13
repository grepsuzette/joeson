package joeson

import (
	"errors"
	"fmt"
	goscanner "go/scanner"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// TokenStream helps parsing tokenized texts.
// Suppose you want to parse a grammar from some pre-tokenized text.
// Tokenization will simplify your grammar, but we would now need
// two systems of source coordinates, e.g. to report usable errors
// to the user.
//
// Suppose a sequence of `=` represents tokens:
//
//	 pos=0    pos=12
//	/        /
//	=======  ===N=====
//	====  ======= ===
//	===========  =====
//
// When parsing fails at offset N, TokenStream can find which offset
// of which token, and translate back to byte offset in original
// text.
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
		// meta interface{} // useless for now
	}
	coord struct {
		token          Token // token found at requested position
		nToken         int   // number of that token, ∈ [0; len(tokens)[
		offsetInToken  int   // offset of requested position in token.Repr, ∈ [0; len(tokens[nToken].Repr)[
		workOffset     int   // offset relative to `work` text. ∈ [0; len(work)[ . Use toWorkOffset()
		originalOffset int   // offset relative to `original` text
	}
)

func (t Token) String() string {
	return fmt.Sprintf("OriginalOffset: %d\tWorkOffset: %d\t%q",
		t.OriginalOffset,
		t.WorkOffset,
		t.Repr,
	)
}

// Create a new token stream.
// `text` is the original, untokenized text.
// `tokens` must have been generated from `text`.
// For go code, you could use for instance "go/scanner".
func NewTokenStream(text string, tokens []Token) *TokenStream {
	lineStarts := []int{0}
	for pos, rune := range text {
		if rune == '\n' {
			lineStarts = append(lineStarts, pos+1)
		}
	}
	var b strings.Builder
	for _, token := range tokens {
		b.WriteString(token.Repr)
	}
	return &TokenStream{tokens, text, b.String(), 0, lineStarts}
}

// `Pos()` returns the byte offset relative to workOffset (the tokenized string)
func (ts *TokenStream) Pos() int {
	return ts.workOffset
}

// `SetPos()` sets the the byte offset relative to workOffset (the tokenized string)
// A value of `n` equals to the length of `work` represents the end of the
// stream (nothing to parse anymore).
func (ts *TokenStream) SetPos(n int) {
	if n < 0 || n > len(ts.work) {
		panic(fmt.Sprintf("%d is out of bound", n))
	}
	ts.workOffset = n
}

// Get current line (in the original text), first line is 1.
func (code *TokenStream) Line() int { return code.PosToLine(code.workOffset) }

// Get current column (in the original text), first column is 1.
func (code *TokenStream) Col() int { return code.PosToCol(code.workOffset) }

// Convert a certain position (byte offset relative to workOffset) to Line.
// Line refers to the original text, and starts at 1.
func (code *TokenStream) PosToLine(workOffset int) int {
	return helpers.BisectRight(
		code.lineStarts,
		code.coords(workOffset).originalOffset,
	) - 1
}

// Convert a certain position (byte offset relative to workOffset) to Column.
// Column refers to the original text, and start at 1.
func (code *TokenStream) PosToCol(workOffset int) int {
	return code.coords(workOffset).originalOffset -
		code.lineStarts[code.PosToLine(workOffset)]
}

func (code *TokenStream) Code() string { return code.original }

// Get the length in bytes of the original text
func (code *TokenStream) Length() int { return len(code.original) }

// Get the length in bytes of the tokenized text
func (code *TokenStream) workLength() int { return len(code.work) }

// Get the string from current position until the start of string `needle` is found.
// Update current position accordingly (**after** `needle` if found).
func (code *TokenStream) GetUntil(needle string) string {
	offset := strings.Index(code.work[code.workOffset:], needle)
	if offset == -1 {
		offset = len(code.work)
	} else {
		offset += code.workOffset // because we searched from this pos
		offset += len(needle)     // what we're after is the length in bytes
	}
	oldWorkOffset := code.workOffset
	code.workOffset = offset
	return code.work[oldWorkOffset:offset]
}

// Take a look `n` runes backwards or forwards, depending on the sign of n,
// return the string contained in the interval made with the current position.
// don't update position
func (code *TokenStream) PeekRunes(n int) string {
	if n <= 0 {
		return helpers.LastNRunes(code.work[:code.workOffset], -n)
	} else {
		var b strings.Builder
		i := 0
		for _, rune := range code.work[code.workOffset:] {
			b.WriteRune(rune)
			i++
			if i >= n {
				break
			}
		}
		return b.String()
	}
}

// Extract the string contained at lines [least(n...)+currentLine; greatest(n...)+currentLine], backwards or forwards,
// When only 1 value is given, a second value of 0 is implied to create a range.
// For TokenStream, PeekLines() is mostly meant for printing purposes;
// it responds with the original text, not the tokenized one.
func (code *TokenStream) PeekLines(n ...int) string {
	if len(n) == 0 {
		n = []int{0}
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
		end = len(code.original)
	}
	return helpers.SliceString(code.original, start, end)
}

// Match func(rune) bool against rune at current position.
// didMatch indicates whether is succeeded. If so the rune is m and position is
// updated. When at EOF it never match.
func (code *TokenStream) MatchRune(f func(rune) bool) (didMatch bool, m rune) {
	if code.workOffset >= code.workLength() {
		return false, '\x00' // never match at EOF
	}
	var ret rune
	newPos := code.workOffset
	iter := 0
	for offset, rune := range code.work[code.workOffset:] {
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
	if newPos == code.workOffset {
		// when not updated, it means rune matched was the last in text
		code.SetPos(len(code.work))
	} else {
		code.SetPos(newPos)
	}
	return true, ret
}

// Match string `s` against current position.
// didMatch indicates whether is succeeded
// in which case the match is in `m`
func (code *TokenStream) MatchString(s string) (didMatch bool, m string) {
	if s != code.work[code.workOffset:helpers.Min(code.workOffset+len(s), len(code.work))] {
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
			from := code.workOffset + firstMatchLoc[0]
			to := helpers.Min(code.workOffset+firstMatchLoc[1], len(code.work))
			s := code.work[from:to]
			code.workOffset += firstMatchLoc[1]
			return true, s
		}
	}
}

// Single line information to be included in parse errors
func (code *TokenStream) Print() string {
	var o strings.Builder
	originalOffset := code.coords(code.Pos()).originalOffset
	o.WriteString("Code at offset ")
	o.WriteString(BoldYellow(strconv.Itoa(originalOffset)))
	o.WriteString("/")
	o.WriteString(BoldYellow(strconv.Itoa(len(code.original))))
	o.WriteString(": '")
	o.WriteString(Cyan(helpers.SliceString(code.original, helpers.Max(0, originalOffset-20), originalOffset)))
	o.WriteString(BoldCyan("|"))
	o.WriteString(BoldWhite(helpers.SliceString(code.original, originalOffset, originalOffset+40)) + "'")
	return o.String()
}

// multiline print, for debugging purposes
func (code *TokenStream) PrintDebug() string {
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

func (code *TokenStream) PrintWorkText() string {
	return "Work text (tokenized):\n" + code.work + "\n"
}

// Get all possible coordinates (originalOffset, line, col)
// from provided workOffset (current byte offset relative to the tokenized
// string).
func (code *TokenStream) coords(workOffset int) coord {
	if len(code.tokens) == 0 {
		return coord{}
	}
	// find most advanced token number, such that the following token would begin
	// after workOffset.
	nToken := 0
	var token Token
	for {
		if nToken >= len(code.tokens)-1 {
			break // don't panic, it can make sense when a token was inserted
		}
		token = code.tokens[nToken]
		if workOffset < token.WorkOffset+len(token.Repr) {
			break
		}
		nToken++
	}
	if nToken >= len(code.tokens) {
		panic(fmt.Sprintf("nToken=%d goes beyond len(code.tokens)=%d\n"+
			"%s", nToken, len(code.tokens), code.Print()))
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

// use by error handler below
var scanErrors []error

func handleErrors(pos token.Position, msg string) {
	scanErrors = append(scanErrors, scannerError{pos, msg})
}

// TokenStreamFromGoCode is a special function transforming
// some go code into a TokenStream. You can then call
// `yourGrammar.ParseTokens(ts TokenStream)` directly.
func TokenStreamFromGoCode(source string) (*TokenStream, error) {
	var scan goscanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(source))
	scan.Init(file, []byte(source), handleErrors, 0 /*goscanner.ScanComments*/)
	if scan.ErrorCount > 0 {
		if scan.ErrorCount != len(scanErrors) {
			panic("assert") // errors must have been collected
		}
		return nil, errors.Join(scanErrors...)
	}
	tokens := []Token{}
	workOffset := 0
	prev := ""

	// Go lexer adds an automatic semicolon when the line's last token is:
	// * an identifier
	// * an integer, floating-point, imaginary, rune, or string literal
	// * one of the keywords break, continue, fallthrough, or return
	// * one of the operators and delimiters ++, --, ), ], or }
	var b strings.Builder
	mustInsertSpaceAfter := regexp.MustCompile("[a-zA-Z0-9_=]$")
	for {
		pos, tok, lit := scan.Scan()
		if tok == token.EOF {
			break
		}
		s := ""
		tokStr := tok.String()
		if tokStr == ";" && lit == "\n" {
			s = ";\n"
		} else if lit != "" {
			if mustInsertSpaceAfter.MatchString(prev) {
				s = " " + lit
			} else {
				s = lit
			}
		} else {
			if mustInsertSpaceAfter.MatchString(prev) &&
				(tok.IsOperator() && tok == token.COMMA) {
				s = " " + tokStr
			} else {
				s = tokStr
			}
		}
		workOffset += len(prev)
		prev = s
		tokens = append(tokens, Token{s, int(pos) - 1, workOffset})
		b.WriteString(s)
	}
	return NewTokenStream(source, tokens), nil
}

type scannerError struct {
	pos token.Position
	msg string
}

func (se scannerError) Error() string {
	return fmt.Sprintf("there was an error at %s: %s", se.pos.String(), se.msg)
}
