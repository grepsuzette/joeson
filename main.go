package main

import (
	"fmt"
	. "grepsuzette/joeson/ast/handcompiled"
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/helpers"
	line "grepsuzette/joeson/line"
	"strconv"
)

func o(a ...any) line.OLine               { return line.O(a...) }
func i(a ...any) line.ILine               { return line.I(a...) }
func Rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }

func main() {
	fmt.Println("readying hand compiled grammar")
	fmt.Println(JOESON_GRAMMAR_RULES)
	fmt.Println("--- new grammar from above rules... ---")
	gm := line.NewGrammarFromLines("joeson from handcompiled", JOESON_GRAMMAR_RULES)
	fmt.Println("done, made grammar from JOESON_GRAMMAR_RULES")
	fmt.Println("Name: " + Cyan(gm.GetGNode().Name))
	fmt.Println("Rules: " + BoldYellow(strconv.Itoa(gm.NumRules)))
	fmt.Println("CountRules: " + BoldYellow(strconv.Itoa(gm.CountRules())))
	keys := helpers.SortIntKeys(gm.Id2Rule)
	for _, i := range keys {
		fmt.Printf("%d %s\n", i-1, gm.Id2Rule[i].ContentString())
	}
	// grammar.ParseCode(```
	// 	  (343+32) * 392 - 1
	// ```)
}
