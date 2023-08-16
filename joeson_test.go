package joeson

import (
	"fmt"
	"testing"
	"time"

	"github.com/grepsuzette/joeson/helpers"
)

// just compile the intention grammar
func TestParseIntention(t *testing.T) {
	gmIntention := GrammarFromLines(
		"gmIntention",
		IntentionRules(),
	)
	if !gmIntention.IsReady() || gmIntention.numrules != JoesonNbRules {
		t.Fail()
	}
}

// bootstrap the intention grammar:
// __joeson__ -> __intention__ -> an arbitrary grammar -> parse a string
func TestBootstrap(t *testing.T) {
	gmJoeson := NewJoeson()
	gmIntention := GrammarWithOptionsFromLines(
		"gmIntention",
		GrammarOptions{LazyGrammar: helpers.LazyFromValue[*Grammar](gmJoeson)},
		IntentionRules(),
	)
	gmJoeson.Bomb() // destroy the grammar to make sure it plays no part below
	gmDebuglabel := GrammarWithOptionsFromLines(
		"gmFoo",
		GrammarOptions{LazyGrammar: helpers.LazyFromValue[*Grammar](gmIntention)},
		[]Line{
			o(Named("In", "l:Br")),
			i(Named("Br", "'Toy' | 'BZ'")),
		},
	)
	gmDebuglabel.PrintRules()
	ast := gmDebuglabel.ParseString("Toy")
	if IsParseError(ast) {
		t.Error(ast.(ParseError).String())
	} else {
		nm := ast.(*NativeMap)
		fmt.Println(ast.String())
		if s, exists := nm.GetStringExists("l"); exists {
			if s != "Toy" {
				t.Errorf("expected NativeMap with 'l' label containing 'Toy', got %s\n", s)
			}
		} else {
			t.Errorf("expected NativeMap with 'l' label")
		}
	}
}

// This benchmark replicates the original joeson_test.coffee test.
// It used to be called TestManyTimes with manual timing, but
// those are now replaced by using go benchmark.
//
// As of August 16, 2023 commit 9894c1 result of `go test -bench=Intention -run=^noother -count=3` is
// goos: linux
// goarch: amd64
// pkg: github.com/grepsuzette/joeson
// cpu: Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
// BenchmarkIntention/joeson.Parse(intentionGrammarRules)-8                1000000000               0.004663 ns/op
// BenchmarkIntention/joeson.Parse(intentionGrammarRules)-8                1000000000               0.004493 ns/op
// BenchmarkIntention/joeson.Parse(intentionGrammarRules)-8                1000000000               0.004490 ns/op
// PASS
func BenchmarkIntention(b *testing.B) {
	parsedGrammar := GrammarFromLines(
		"gmIntention",
		IntentionRules(),
	)
	var frecurse func(rule Line, indent int, name string)
	frecurse = func(rule Line, indent int, name string) {
		switch v := rule.(type) {
		case ALine:
			if name != "" {
				// b.Logf("%s%s\n", helpers.Indent(indent), Red(name+":"))
			}
			for _, subline := range v.Array {
				frecurse(subline, indent+1, "")
			}
		case OLine:
			if name == "" {
				name = v.name
			}
			frecurse(v.content, indent, name)
		case ILine:
			frecurse(v.content, indent, v.name)
		case cLine:
			// b.Logf("%s%s\n", helpers.Indent(indent), String(v.Parser))
		case sLine:
			// parse the rules of the intention grammar, one line at a time
			ast := parsedGrammar.ParseString(v.Str)
			if IsParseError(ast) {
				panic(ast.String())
			} else if false { // set to true to enable as in original joeson implementation
				sName := ""
				if name != "" {
					sName = Red(helpers.PadLeft(name+":", 10-indent*2))
				}
				sResult := Red("null")
				if ast != nil {
					sResult = Yellow(String(ast))
				}
				b.Logf("%s%s%s %s\n", helpers.Indent(indent), sName, sResult, White(""))
			}
		default:
			b.Logf("unknown -----%#v %T\n", v, v)
		}
	}
	b.Run("joeson.Parse(intentionGrammarRules)", func(b *testing.B) {
		frecurse(NewALine(IntentionRules()), 0, "")
	})
}

// Meaningless test that can now be deleted
func TestSquareroot(t *testing.T) {
	gm := GrammarFromLines(
		"gmSqr",
		[]Line{
			o(Named("sqr", "w:word '(' n:int ')'")),
			i(Named("word", "[a-z]{1,}")),
			i(Named("int", "/-?[0-9]{1,}/"), func(it Ast) Ast { return NewNativeIntFrom(it) }),
		})
	ast := gm.ParseString("squareroot(-1)")
	if IsParseError(ast) {
		t.Error(ast.String())
	} else {
		nmap := ast.(*NativeMap)
		if w, exists := nmap.GetStringExists("w"); !exists || w != "squareroot" {
			t.Error("was expecting w == squareroot")
		} else if n, exists := nmap.GetIntExists("n"); !exists || n != -1 {
			if !exists {
				t.Error("label n not found")
			} else {
				t.Errorf("was expecting n == -1, but got %d\n", n)
			}
		}
	}
}

func Test_LeftRecursion(t *testing.T) {
	gm := GrammarWithOptionsFromLines(
		"leftRecursion",
		GrammarOptions{TraceOptions: Mute()},
		[]Line{
			o(named("Input", "expr:Expression")),
			i(named("Expression", "Expression _ binary_op _ Expression | UnaryExpr")),
			i(named("binary_op", "'+'")),
			i(named("UnaryExpr", "[0-9]+")),
			i(named("_", "[ \t]*")),
		},
	)
	res := gm.ParseString("123 + 456")
	fmt.Println(res.String())
}

// Showing difference between captured and uncaptured Str
func TestCapturingStr(t *testing.T) {
	// "'0x' [0-9a-f]{2,2}" parsing "0x7d" will only capture "7d"
	// To capture "0x7d" you can have a label: "prefix:'0x' [0-9a-f]{2,2}"
	// This time it should capture all of it.
	{
		// str not captured
		ast := GrammarFromLines("gm",
			[]Line{o(named("Input", "'0x' [0-9a-f]{2,2}"))}).ParseString("0x7d")
		if s := ast.(*NativeArray).Concat(); s != "7d" {
			t.Errorf("for test 1 unexpected result %s", s)
		}
	}
	{
		// now captured using labels
		ast := GrammarFromLines("gm", []Line{o(named("Input", "captureMe:'0x' captureMeToo:[0-9a-f]{2,2}"))}).ParseString("0x7d")
		if s := ast.(*NativeMap).Concat(); s != "0x7d" {
			t.Errorf("for test 2 unexpected result %s", s)
		}
	}
}

func TestNativeIntFromBool(t *testing.T) {
	{
		n := NewNativeIntFromBool(true)
		if !n.Bool() {
			t.Error()
		}
	}
	{
		n := NewNativeIntFromBool(false)
		if n.Bool() {
			t.Error()
		}
	}
}
