package main

// a lisp like REPL.
// Try
// (+ 29 52)
// (+ (- 439 2) (+ 92 2))

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	j "github.com/grepsuzette/joeson"
)

const E = "should not happen" // panic(E)

func colorKeyword(s string) string { return (s) }
func colorParen(s string) string   { return j.BoldCyan(s) }
func colorComma(s string) string   { return j.BoldGreen(s) }
func quoted(s string) string       { return j.Red(`"`) + j.Blue(s) + j.Red(`"`) }

func grammar() *j.Grammar { return j.GrammarFromLines("lisp", grammarRules) }

// REPL using joeson grammar to parse inputs,
// : and feeding a VM with that ast
func main() {
	gm := grammar()
	vm := NewMachine()
	fmt.Println("pseudo lisp REPL. This will try to interpret some lisp like expr")
	for true {
		fmt.Print("Eval: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Bye!")
				break
			} else {
				panic(err)
			}
		} else {
			s := strings.TrimSuffix(input, "\n")
			s = strings.Trim(s, " ")
			if !strings.HasPrefix(s, "(") {
				s = "(" + s + ")"
			}
			ast := gm.ParseString(s)
			if j.IsParseError(ast) {
				fmt.Println("Parse error: " + ast.String())
			} else {
				if expr, ok := ast.(Expr); ok {
					evaluatedExpr := vm.Eval(expr)
					fmt.Println(evaluatedExpr.String())
				} else {
					panic(E)
				}
			}
		}
	}
}
