package joeson

import (
	"fmt"
	"testing"
	"time"

	"github.com/grepsuzette/joeson/helpers"
)

// this allows tracing and diffing,
// it does not do more than compiling the intention grammar
func TestParseIntention(t *testing.T) {
	gmIntention := GrammarFromLines(IntentionRules(), "gmIntention")
	if !gmIntention.IsReady() || gmIntention.numrules != JoesonNbRules {
		t.Fail()
	}
}

// This test bootstraps the intention grammar:
// __joeson__ -> __intention__ -> an arbitrary grammar -> parses a string
func TestBootstrap(t *testing.T) {
	gmJoeson := NewJoeson()
	gmIntention := GrammarFromLines(
		IntentionRules(),
		"gmIntention",
		GrammarOptions{LazyGrammar: helpers.NewLazyFromValue[*Grammar](gmJoeson)},
	)
	gmJoeson.Bomb() // destroy the grammar to make sure it plays no part below
	gmDebuglabel := GrammarFromLines(
		[]Line{
			o(Named("In", "l:Br")),
			i(Named("Br", "'Toy' | 'BZ'")),
		},
		"gmFoo",
		GrammarOptions{LazyGrammar: helpers.NewLazyFromValue[*Grammar](gmIntention)},
	)
	gmDebuglabel.PrintRules()
	ast := gmDebuglabel.ParseString("Toy")
	if IsParseError(ast) {
		t.Error(ast.(ParseError).String())
	} else {
		var nm NativeMap = ast.(NativeMap)
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

func TestManyTimes(t *testing.T) {
	// this test replicates the original joeson_test.coffee
	start := time.Now()
	nbIter := 100
	parsedGrammar := GrammarFromLines(IntentionRules(), "gmIntention")
	var frecurse func(rule Line, indent int, name string)
	frecurse = func(rule Line, indent int, name string) {
		switch v := rule.(type) {
		case ALine:
			if name != "" {
				fmt.Printf("%s%s\n", helpers.Indent(indent), Red(name+":"))
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
			fmt.Printf("%s%s\n", helpers.Indent(indent), String(v.Parser))
		case sLine:
			// parse the rules of the intention grammar, one line at a time
			ast := parsedGrammar.ParseString(v.Str)
			if IsParseError(ast) {
				panic(ast.String())
			} else {
				sName := ""
				if name != "" {
					sName = Red(helpers.PadLeft(name+":", 10-indent*2))
				}
				sResult := Red("null")
				if ast != nil {
					sResult = Yellow(String(ast))
				}
				fmt.Printf("%s%s%s %s\n", helpers.Indent(indent), sName, sResult, White(""))
			}

		default:
			fmt.Printf("unknown -----%#v %T\n", v, v)
		}
	}
	for i := 0; i < nbIter; i++ {
		frecurse(NewALine(IntentionRules()), 0, "")
	}
	fmt.Printf("Duration for %d iterations: %d ms\n", nbIter, time.Since(start).Milliseconds())
}

// short grammar was useful for debugging. Kept for the good memories
// __joeson__ -> an arbitrary grammar -> parse a string
func TestDebugLabel(t *testing.T) {
	debuglabel := GrammarFromLines(
		[]Line{
			o(Named("In", "l:Br")),
			i(Named("Br", "'Toy' | 'BZ'")),
		},
		"gmDebugLabel",
	)
	ast := debuglabel.ParseString("Toy")
	if IsParseError(ast) {
		t.Error(ast.String())
	} else {
		if nm, isNativeMap := ast.(NativeMap); !isNativeMap {
			t.Errorf("expected NativeMap, got %T. String: %s\n", ast, ast.String())
		} else {
			// in two operations...
			if label, exists := nm.GetExists("l"); !exists {
				t.Fail()
			} else if label.(NativeString).Str != "Toy" {
				t.Fail()
			}
			// ...or in 1 operation
			if label, exists := nm.GetStringExists("l"); !exists || label != "Toy" {
				t.Fail()
			}
		}
	}
}

func TestSquareroot(t *testing.T) {
	gm := GrammarFromLines(
		[]Line{
			o(Named("sqr", "w:word '(' n:int ')'")),
			i(Named("word", "[a-z]{1,}")),
			i(Named("int", "/-?[0-9]{1,}/"), func(it Ast) Ast { return NewNativeIntFrom(it) }),
		},
		"gmSqr",
	)
	ast := gm.ParseString("squareroot(-1)")
	if IsParseError(ast) {
		t.Error(ast.String())
	} else {
		nmap := ast.(NativeMap)
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
