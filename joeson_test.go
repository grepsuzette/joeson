package main

import (
	"fmt"
	// . "grepsuzette/joeson/ast/handcompiled"
	// . "grepsuzette/joeson/ast/raw"
	// . "grepsuzette/joeson/ast/handcompiled"
	// . "grepsuzette/joeson/ast/raw"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/grammars"
	line "grepsuzette/joeson/line"
	"testing"
	"time"
)

func o(a ...any) line.OLine               { return line.O(a...) }
func i(a ...any) line.ILine               { return line.I(a...) }
func Rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }
func Named(name string, lineStringOrAstnode any) line.NamedRule {
	return line.Named(name, lineStringOrAstnode)
}

// NewJoeson() is the native, manually compiled joeson grammar defined in handcompiled.go
func TestHandcompiled(t *testing.T) {
	gm := grammars.NewJoeson()
	if gm.GetGNode().Name != grammars.JOESON_GRAMMAR_NAME {
		t.Fail()
	}
	if gm.CountRules() != gm.NumRules || gm.CountRules() != grammars.JoesonNbRules {
		t.Errorf("Expected %d rules, got %d\n", grammars.JoesonNbRules, gm.CountRules())
	}
	if !gm.IsReady() {
		t.Fail()
	}
}

// bootstrapped grammar, using joeson to define itself,
// similar to joeson_test.coffee
func TestRaw(t *testing.T) {
	raw := line.NewGrammarFromLines(
		"bootstrapped grammar",
		grammars.RAW_GRAMMAR(),
		grammars.NewJoeson(),
	)
	if !raw.IsReady() {
		t.Fail()
	}
}

func Test100Times(t *testing.T) {
	// this test comes directly from joeson_test.coffee
	start := time.Now()
	iter := 100
	for i := 0; i < iter; i++ {
		// testGrammar(line.NewALine(RAW_GRAMMAR()), 0, "")
		fmt.Println(line.NewALine(grammars.RAW_GRAMMAR()).StringIndent(0))
		fmt.Println("-------------")
	}
	fmt.Printf("Duration for %d iterations: %d ms\n", iter, time.Now().Sub(start).Milliseconds())
}

// this short grammar was useful for debugging
func TestDebugLabel(t *testing.T) {
	debuglabel := line.NewGrammarFromLines(
		"gm_DebugLabel",
		[]line.Line{
			o(Named("In", "l:Br")),
			i(Named("Br", "'Toy' | 'BZ'")),
		},
		grammars.NewJoeson(),
	)
	debuglabel.PrintRules()
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
