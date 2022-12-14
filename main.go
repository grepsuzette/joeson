package main

import (
	"fmt"
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/ast/handcompiled"
	. "grepsuzette/joeson/colors"
	// "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	line "grepsuzette/joeson/line"

	"strconv"
)

func o(a ...any) line.OLine               { return line.O(a...) }
func i(a ...any) line.ILine               { return line.I(a...) }
func Rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }
func Named(name string, lineStringOrAstnode any) line.NamedRule {
	return line.Named(name, lineStringOrAstnode)
}

func main() {
	// fmt.Println("readying hand compiled grammar")
	// fmt.Println(JOESON_GRAMMAR_RULES())
	// fmt.Println("--- new grammar from above rules... ---")
	gm := line.NewGrammarFromLines("joeson from handcompiled", JOESON_GRAMMAR_RULES(), ast.NewEmptyGrammarNamed("empty grammar"))
	fmt.Println("done, made grammar from JOESON_GRAMMAR_RULES()")
	fmt.Println("Name: " + Cyan(gm.GetGNode().Name))
	fmt.Println("Rules: " + BoldYellow(strconv.Itoa(gm.NumRules)))
	fmt.Println("CountRules: " + BoldYellow(strconv.Itoa(gm.CountRules())))
	keys := helpers.SortIntKeys(gm.Id2Rule)
	for _, i := range keys {
		fmt.Printf("%d %s\n", i, gm.Id2Rule[i].ContentString())
	}

	// fmt.Printf("is grammar ready? %v\n", gm.IsReady())
	// ast := gm.ParseString("FOO", core.ParseOptions{})
	// fmt.Println(ast.ContentString())
	// grammar.ParseCode(```
	// 	  (343+32) * 392 - 1
	// ```)
	// example_20_calculator_no_cb()
}

func example_20_calculator_no_cb() {
	x := []line.Line{
		o(Named("Input", Rules(
			o("expr:Expr EOF"),
			o(Named("Expr", Rules(
				o("_ first:Term rest:( _ AddOp _ Term )* _"),
				o(Named("Term", Rules(
					o("first:Factor rest:( _ MulOp _ Factor )*"),
					o(Named("Factor", Rules(
						o("ExprInParen | Integer"),
						o(Named("ExprInParen", "'(' expr:Expr ')'")),
					))),
					o(Named("MulOp", "'*' | '/'")),
				))),
				o(Named("AddOp", "'+' | '-'")),
			))),
		))),
		i(Named("Integer", "'-'? [0-9]+")),
		i(Named("_", "[ \n\t\r]*")),
		i(Named("EOF", "!.")),
	}
	fmt.Println(x)
}

// func example_2_calculator_eval() {
// 	x := []Line{
// 		o(Named("Input", Rules(
// 			o("expr:Expr EOF"),
// 			o(Named("Expr", Rules(
// 				o("_ first:Term rest:( _ AddOp _ Term )* _"),
// 				o(Named("Term", Rules(
// 					o("first:Factor rest:( _ MulOp _ Factor )*"),
// 					o(Named("Factor", Rules(
// 						o("ExprInParen | Integer"),
// 						o(Named("ExprInParen", "'(' expr:Expr ')'")),
// 					),
// 					o(Named("MulOp", "'*' | '/'")),
// 				))),
// 				o(Named("AddOp", "'+' | '-'")),
// 			))),
// 		))),
// 		i(Named("Integer", "'-'? [0-9]+")),
// 		i(Named("_", "[ \n\t\r]*")),
// 		i(Named("EOF", "!.")),
// 	}
// 	fmt.Println(ast.ContentString())
// }
