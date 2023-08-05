package joeson

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// What is tested here:
//
// | test name                 | counters, line, col.. | Grammar.ParseTokens | other...
// | TestTokenStreamInternals  | 1                     | 0                   |
// | TestLongerSample          | 0                     | 0                   | prints tokenized version (ultimately will be compared to some expected large tokenized version, when finalized)
// | TestMiniatures            | 0                     | 1                   |
// | TestExpectedTokenization  | 0                     | 0                   | tests short sample expected tokenization

type found struct {
	search   string // text to search in tokenstream.work
	work     int    // expected offset of `search` text in TokenStream.work, untested if <0
	original int    // expected matching offset in TokenStream.original, untested if <0
	line     int    // expected line in TokenStream.original, untested if <0
	col      int    // expected col in TokenStream.original, untested if <0
}

const source = `
	// RuneStream is a very simple code holder, cursor, matcher.
	type RuneStream struct {
		text       string
		pos        int // "Hello, 世界, X" <- Pos of o is 4, Pos of 界 is 10
		lineStarts []int
	}
	`

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

func TestTokenStreamInternals(t *testing.T) {
	tokens, e := TokenStreamFromGoCode(source)
	if e != nil {
		t.Errorf("Failed to tokenize %q: %s", source, e.Error())
	}
	// test initial counters
	if tokens.Pos() != 0 {
		t.Errorf("at start Pos() == 0: %d\n", tokens.Pos())
	}
	tokens.SetPos(0)
	if tokens.Pos() != 0 {
		t.Errorf("After SetPos(0), Pos() == 0: %d\n", tokens.Pos())
	}
	// --- test private functions -------------
	{
		originalOffset := tokens.coords(0).originalOffset // @ "type "
		expect := 64
		if originalOffset != expect {
			t.Errorf("coords(0).originalOffset should be %d, got %d\n", expect, originalOffset)
		}
	}
	{
		originalOffset := tokens.coords(16).originalOffset // @ "struct "
		expect := 81
		if originalOffset != expect {
			t.Errorf("coords(16).originalOffset should be %d, got %d\n", expect, originalOffset)
		}
	}
	{
		originalOffset := tokens.coords(26).originalOffset // @ "t|ext " (at tokenOffset 1 of token "text ")
		expect := 94
		if originalOffset != expect {
			t.Errorf("coords(26).originalOffset should be %d, got %d\n", expect, originalOffset)
		}
	}
	if tokens.PosToLine(0) != 2 {
		t.Errorf("PosToLine(0) == 2: %d\n", tokens.PosToLine(0))
	}
	// --- test PosToLine PosToCol Line Col Length ------------------
	posType := found{search: "type", work: 0, original: 64, line: 2, col: 1}
	posText := found{search: "text", work: 23, original: 91, line: 3, col: 2}
	pos_ext := found{search: "ext", work: 24, original: 92, line: 3, col: 3}
	posStrn := found{search: "string", work: 28, original: 103, line: 3, col: 14}
	posRBrk := found{search: "]", work: 56, original: 197, line: 5, col: 14}
	testHas(t, tokens, posType)
	testHas(t, tokens, posText)
	testHas(t, tokens, pos_ext)
	testHas(t, tokens, posStrn)
	testHas(t, tokens, posRBrk)

	// --- test MatchString MatchRegexp PeekRunes PeekLines ---------
	// jump to "string"
	tokens.workOffset = posStrn.work
	// now it should match "string", and advance
	if ok, m := tokens.MatchString("string"); !ok || m != "string" {
		t.Error("Failed to match string \"string\"")
	}
	re1 := regexp.MustCompilePOSIX(`[ ;\t\n\r]*pos`)
	if ok, _ := tokens.MatchRegexp(*re1); !ok {
		t.Errorf("failed to match regexp %q\nThe TokenStream.Print(): %s", re1.String(), tokens.Print())
	}
	re2 := regexp.MustCompilePOSIX(`NO`)
	if ok, _ := tokens.MatchRegexp(*re2); ok {
		t.Error("Should not have matched regexp " + re2.String())
	}

	fmt.Println(tokens.Print())
}

func TestLongerSample(t *testing.T) {
	source := `
package p
import fmt "fmt"
const pi = 3.14
type T struct{
	a int
	 b string
	  c float
}
var x int
func f() { L: }

var (
	_ int = 23
	_ string = "abc"
)

type (
	Alpha struct {
		a int
		b string
		 c float
	}
	Beta struct {
		a []string{256} // comment
	}
	 Gamma struct {}
)

func f() {
	if true {
		if false {
			// after 1 below, there should have an automatic ; inserted
			n := 1
			fmt.Println("no")
		}
	}
	a := []string{
		"foo",
		 "bar",
	 "baz",
	}
}
	`
	fmt.Println(source)
	if tokens, e := TokenStreamFromGoCode(source); e != nil {
		t.Error(e.Error())
	} else {
		work := tokens.PrintWorkText()
		fmt.Println(work)

		// no test yet during development phase
		// it's just printing what we obtain.
		// if work != expectWork {
		// 	t.Fail()
		// }
	}
}

const expectWork string = `package p;
import fmt "fmt";
const pi= 3.14;
type T struct{a int;
b string;
c float;
};
var x int;
func f(){L:};
var(_ int= 23;
_ string= "abc";
);
type(Alpha struct{a int;
b string;
c float;
};
Beta struct{a[]string{256};
};
Gamma struct{};
);
func f(){if true{if false{n:= 1;
fmt.Println("no");
};
};
a:=[]string{"foo","bar","baz",};
};
`

func TestExpectedTokenization(t *testing.T) {
	for _, a := range [][]string{
		{"a", "a;\n"},
		{"1234+  (-321)", "1234+(-321);\n"},
	} {
		if tokens, e := TokenStreamFromGoCode(a[0]); e != nil {
			t.Error(e.Error())
		} else {
			joined := ""
			for _, token := range tokens.Tokens() {
				joined += token.Repr
			}
			if joined != a[1] {
				// work := tokens.PrintWorkText()
				// work := tokens.Print()
				// fmt.Println(work)
				t.Errorf("tokenizing %q should have produced %q, not %q",
					a[0], a[1], joined,
				)
			}
		}
	}
}

// have a small grammar, parse many small tokenized go expressions parsed
func TestMiniatures(t *testing.T) {
	gm := GrammarFromLines(
		"miniatures",
		[]Line{
			o(named("Input", rules(
				o(named("Number", "[0-9]+ term")),
				o(named("LnChar", "/'[^']+'/ term")),
			))),
			i(named("term", "';' '\n'*")),
		})
	ter := ";\n"
	for _, a := range [][]string{
		{"1234", "1234" + ter, "1234"},
		{`'\f'`, `'\f'` + ter, `'\f'`},
		{`'\n'`, `'\n'` + ter, `'\n'`},
	} {
		if len(a) != 3 {
			t.Errorf("Expected array of len 3, got len %d for %v", len(a), a)
			continue
		}

		miniature := a[0]
		tokenized := a[1]
		expectation := a[2] // stringified expected parse result

		if tokens, err := TokenStreamFromGoCode(miniature); err != nil {
			t.Errorf("Fail to tokenize %q: %s", miniature, err.Error())
		} else {
			if tokens.work != tokenized {
				t.Errorf("%q should have been tokenized as %q, got %q", miniature, tokenized, tokens.work)
			} else {
				ast := gm.ParseTokens(tokens)
				reality := ast.String()
				if IsParseError(ast) {
					if !strings.HasPrefix(expectation, "ERROR ") {
						t.Errorf("%q parsed as unexpected error %s", miniature, reality)
					}
					if !strings.HasPrefix(expectation, reality) {
						t.Errorf("ParseError when parsing %q. Expected %q, got %s", miniature, expectation, reality)
					}
				} else {
					if strings.HasPrefix(expectation, "ERROR") {
						t.Errorf("%q parsed as %s but expected %s", miniature, reality, tokenized)
					}
				}
			}
		}
	}
}
