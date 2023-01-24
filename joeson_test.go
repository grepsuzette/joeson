package main

import (
	"fmt"
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"grepsuzette/joeson/line"
	"testing"
	"time"
)

// useful aliases
func o(a ...any) line.OLine { return line.O(a...) }
func i(a ...any) line.ILine { return line.I(a...) }

func Rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }
func Named(name string, lineStringOrAstnode any) line.NamedRule {
	return line.Named(name, lineStringOrAstnode)
}

// Basic test, doesn't do much
func TestHandcompiled(t *testing.T) {
	gm := line.NewJoeson()
	if gm.GetGNode().Name != line.JoesonGrammarName {
		t.Fail()
	}
	if gm.NumRules != line.JoesonNbRules {
		t.Errorf("Expected %d rules, got %d\n", line.JoesonNbRules, gm.NumRules)
	}
	if !gm.IsReady() {
		t.Fail()
	}
}

// Main test
// Parse joeson_intention using joeson_handcompiled
// It's similar to joeson_test.coffee
func TestParseIntention(t *testing.T) {
	gmIntention := line.GrammarFromLines(line.IntentionRules(), "gmIntention")
	if !gmIntention.IsReady() || gmIntention.NumRules != line.JoesonNbRules {
		t.Fail()
	}
}

// This test bootstraps the intention grammar:
// __joeson__ -> __intention__ -> an arbitrary grammar -> parses a string
func TestBootstrap(t *testing.T) {
	gmJoeson := line.NewJoeson()
	gmIntention := line.GrammarFromLines(
		line.IntentionRules(),
		"gmIntention",
		line.GrammarOptions{LazyGrammar: helpers.NewLazyFromValue[*ast.Grammar](gmJoeson)},
	)
	gmJoeson.Bomb() // destroy the grammar! to make sure it plays no part below
	gmDebuglabel := line.GrammarFromLines(
		[]line.Line{
			o(Named("In", "l:Br")),
			i(Named("Br", "'Toy' | 'BZ'")),
		},
		"dbglbl/bootst",
		line.GrammarOptions{LazyGrammar: helpers.NewLazyFromValue[*ast.Grammar](gmIntention)},
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

// TODO benchmarks later
func TestManyTimes(t *testing.T) {
	// this test replicates joeson_test.coffee
	start := time.Now()
	nbIter := 10
	parsedGrammar := line.GrammarFromLines(line.IntentionRules(), "gmIntention", line.GrammarOptions{TraceOptions: Mute()})
	var frecurse func(rule line.Line, indent int, name string)
	frecurse = func(rule line.Line, indent int, name string) {
		switch v := rule.(type) {
		case line.ALine:
			// fmt.Println("aLINE n=" + name)
			if name != "" {
				fmt.Printf("%s%s\n", helpers.Indent(indent), Red(name+":"))
			}
			for _, subline := range v.Array {
				frecurse(subline, indent+1, "")
			}
		case line.OLine:
			// fmt.Println("Oline: name=" + name + " VNAme:" + v.Name() + "  stringindent=" + v.StringIndent(indent))
			if name == "" {
				name = v.Name()
			}
			frecurse(v.Content(), indent, name)
		case line.ILine:
			// fmt.Println("Iline: " + v.Name() + "  " + v.StringIndent(indent))
			frecurse(v.Content(), indent, v.Name())
		case line.CLine:
			// fmt.Println("CLINE")
			fmt.Printf("%s%s\n", helpers.Indent(indent), String(v.Ast))
		case line.SLine:
			// fmt.Println("Sline: " + v.Str + " name= " + name)
			// parse the rules of the intention grammar, one line at a time
			if it, err := parsedGrammar.ParseString(v.Str, ParseOptions{Debug: false}); err != nil {
				panic(err)
			} else {
				sName := ""
				if name != "" {
					sName = Red(helpers.PadLeft(name+":", 10-indent*2))
				}
				sResult := Red("null")
				if it != nil {
					sResult = Yellow(String(it))
				}
				fmt.Printf("%s%s%s\n", helpers.Indent(indent), sName, sResult)
			}

		default:
			fmt.Printf("unknown -----%#v %T\n", v, v)
		}
	}
	for i := 0; i < nbIter; i++ {
		frecurse(line.NewALine(line.IntentionRules()), 0, "")
	}
	fmt.Printf("Duration for %d iterations: %d ms\n", nbIter, time.Now().Sub(start).Milliseconds())
}

// short grammar was useful for debugging. Kept for the good memories
// __joeson__ -> an arbitrary grammar -> parse a string
func TestDebugLabel(t *testing.T) {
	debuglabel := line.GrammarFromLines(
		[]line.Line{
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
