package joeson

import (
	"fmt"
	"github.com/grepsuzette/joeson/helpers"
	"testing"
	"time"
)

// this allows tracing and diffing,
// it does not do more than compiling the intention grammar
func TestParseIntention(t *testing.T) {
	gmIntention := GrammarFromLines(IntentionRules(), "gmIntention")
	if !gmIntention.IsReady() || gmIntention.NumRules != JoesonNbRules {
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
	if x, err := gmDebuglabel.ParseString("Toy"); err != nil {
		t.Error(err)
	} else {
		var nm NativeMap = x.(NativeMap)
		fmt.Println(x.ContentString())
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
				fmt.Printf("%s%s\n", helpers.Indent(indent), red(name+":"))
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
			if it, err := parsedGrammar.ParseString(v.Str, ParseOptions{Debug: false}); err != nil {
				panic(err)
			} else {
				sName := ""
				if name != "" {
					sName = red(helpers.PadLeft(name+":", 10-indent*2))
				}
				sResult := red("null")
				if it != nil {
					sResult = yellow(String(it))
				}
				fmt.Printf("%s%s%s %s\n", helpers.Indent(indent), sName, sResult, white(""))
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
	if x, error := debuglabel.ParseString("Toy"); error == nil {
		if nm, isNativeMap := x.(NativeMap); !isNativeMap {
			t.Errorf("expected NativeMap, got %T. ContentString: %s\n", x, x.ContentString())
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
	} else {
		t.Error(error)
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
	if x, error := gm.ParseString("squareroot(-1)"); error != nil {
		t.Error(error)
	} else {
		nmap := x.(NativeMap)
		if w, exists := nmap.GetStringExists("w"); !exists || w != "squareroot" {
			t.Error("was expecting w == squareroot")
		} else if n, exists := nmap.GetIntExists("n"); !exists || n != -1 {
			if !exists {
				t.Error("label n not found")
			} else {
				t.Error(fmt.Sprintf("was expecting n == -1, but got %d\n", n))
			}
		}
	}
}

func TestChoice2(t *testing.T) {
	gm := GrammarFromLines(
		[]Line{
			// this is a meaningful test that will eventually go away
			o(Named("CHOICE", rules(
				o("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Ast) Ast {
					fmt.Println(it.ContentString())
					return it
				}),
				o(Named("SEQUENCE", "WORD _ '_'")),
			))),
			i(Named("_PIPE", "_ '|'")),
			i(Named("WORD", "[A-Z]{1,}")),
			i(Named("_", "(' ' | '\n')*")),
		},
		"gmChoice",
	)
	if x, error := gm.ParseString("CHOICE _"); error != nil {
		t.Error(error)
	} else {
		fmt.Println(x.ContentString())
	}
}
