package main

import (
	"fmt"
	. "grepsuzette/joeson/ast/handcompiled"
	fake "grepsuzette/joeson/fake"
	line "grepsuzette/joeson/line"
	"strconv"
	// . "grepsuzette/joeson/core"
	. "grepsuzette/joeson/colors"
)

var QUOTE string = "'\\''"

// function aliases: fake/ package
func rules(lines ...fake.Line) []fake.Line { return lines }
func i(_ ...any) []fake.Line               { return []fake.Line{} }
func o(_ ...any) []fake.Line               { return []fake.Line{} }

// function aliases: line/ package

func I(a ...any) line.ILine { return line.I(a...) }
func O(a ...any) line.OLine { return line.O(a...) }

func main() {
	fmt.Println(JOESON_GRAMMAR_RULES)
	fmt.Println("--- new grammar from above rules... ---")
	gm := line.NewGrammarFromLines("joeson from handcompiled", JOESON_GRAMMAR_RULES)
	fmt.Println("done, made grammar from JOESON_GRAMMAR_RULES")
	fmt.Println("Name: " + Cyan(gm.GetGNode().Name))
	fmt.Println("Rules: " + BoldYellow(strconv.Itoa(gm.NumRules)))
	fmt.Println("CountRules: " + BoldYellow(strconv.Itoa(gm.CountRules())))
	// fmt.Println(grammar.ToString())

	/* grammar.ParseCode(```
		  (343+32) * 392 - 1
	```)
	*/
}
