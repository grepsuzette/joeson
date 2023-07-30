package joeson

import (
	"fmt"
	"go/scanner"
	"go/token"
	"strings"
	"testing"
)

type found struct {
	search   string // text to search in tokenstream.work
	work     int    // expected offset of `search` text in TokenStream.work, untested if <0
	original int    // expected matching offset in TokenStream.original, untested if <0
	line     int    // expected line in TokenStream.original, untested if <0
	col      int    // expected col in TokenStream.original, untested if <0
}

func TestTokenStream(t *testing.T) {
	// let's use go scanner for this example
	source := `
	// RuneStream is a very simple code holder, cursor, matcher.
	type RuneStream struct {
		text       string
		pos        int // "Hello, 世界, X" <- Pos of o is 4, Pos of 界 is 10
		lineStarts []int
	}
	`
	// build some tokenized text
	var tokenizer scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(source))
	tokenizer.Init(file, []byte(source), nil, 0) // scanner.ScanComments)
	tokens := []Token{}
	workOffset := 0
	prevTokenLen := 0
	for {
		pos, tok, lit := tokenizer.Scan()
		if tok == token.EOF {
			break
		}
		s := ""
		if lit != "" {
			s = lit + " "
		} else {
			s = tok.String() + " "
		}
		workOffset += prevTokenLen
		prevTokenLen = len(s)
		tokens = append(tokens, Token{s, int(pos), workOffset})
	}
	code := NewTokenStream(source, tokens)
	// test initial counters
	if code.Pos() != 0 {
		t.Errorf("at start Pos() == 0: %d\n", code.Pos())
	}
	code.SetPos(0)
	if code.Pos() != 0 {
		t.Errorf("After SetPos(0), Pos() == 0: %d\n", code.Pos())
	}
	// --- test private functions -------------
	{
		originalOffset := code.coords(0).originalOffset // @ "type "
		if originalOffset != 65 {
			t.Errorf("coords(0).originalOffset == 0 should be 65, got %d\n", originalOffset)
		}
	}
	{
		originalOffset := code.coords(16).originalOffset // @ "struct "
		if originalOffset != 81 {
			t.Errorf("coords(16).originalOffset should be 81, got %d\n", originalOffset)
		}
	}
	{
		originalOffset := code.coords(26).originalOffset // @ "t|ext " (at tokenOffset 1 of token "text ")
		if originalOffset != 93 {
			t.Errorf("coords(26).originalOffset should be 93, got %d\n", originalOffset)
		}
	}
	if code.PosToLine(0) != 2 {
		t.Errorf("PosToLine(0) == 2: %d\n", code.PosToLine(0))
	}
	// --- test PosToLine PosToCol Line Col Length ------------------
	testHas(t, code, found{search: "type", work: 0, original: 65, line: 2, col: 2})
	testHas(t, code, found{search: "text", work: 25, original: 92, line: 3, col: 3})
	testHas(t, code, found{search: "ext", work: 26, original: 93, line: 3, col: 4})
	testHas(t, code, found{search: "string", work: 30, original: 103, line: 3, col: 14})
	testHas(t, code, found{search: "]", work: 62, original: 198, line: 5, col: 15})
	fmt.Println(code.Print())
}

func testHas(t *testing.T, code *TokenStream, f found) {
	t.Helper()
	s := "Searching '" + f.search + "' "
	offset := strings.Index(code.work, f.search) // work index of <search string>
	original := code.coords(offset).originalOffset
	line := code.PosToLine(offset)
	col := code.PosToCol(offset)
	if f.work >= 0 && offset != f.work {
		t.Errorf(s+"string.Index(ts.work, \"%s\") expected %d, got %d\n", f.search, f.work, offset)
	}
	if f.original >= 0 && original != f.original {
		t.Errorf(s+"original offset expected %d, got %d\n", f.original, original)
	}
	if f.line >= 0 && line != f.line {
		t.Errorf(s+"PosToLine(%d) expected %d, got %d\n", offset, f.line, line)
	}
	if f.col >= 0 && col != f.col {
		t.Errorf(s+"PosToCol(%d) expected %d, got %d\n", offset, f.col, col)
	}
}
