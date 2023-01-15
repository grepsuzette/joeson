package main

import (
	"fmt"
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/grammars"
	line "grepsuzette/joeson/line"
	"strconv"
	"strings"
	"testing"
)

// this example is the typical example with a calculator
// and is inspired from the mna-pigeon example

var ops = map[string]func(int, int) int{
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
}

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }
func mul(a, b int) int { return a * b }
func div(a, b int) int { return a / b }

func o(a ...any) line.OLine               { return line.O(a...) }
func i(a ...any) line.ILine               { return line.I(a...) }
func Rules(lines ...line.Line) line.ALine { return line.NewALine(lines) }
func Named(name string, lineStringOrAstnode any) line.NamedRule {
	return line.Named(name, lineStringOrAstnode)
}

// extract the "expr" key from a result `x` to an int
// or, failing to do that, call FailNow()
func extractResult(t *testing.T, x Ast) int {
	t.Helper()
	if n, exists := x.(NativeMap).GetIntExists("expr"); exists {
		return n
	} else {
		t.Errorf("Failed to find a result like NativeMap{expr:<INT>} in " + x.ContentString())
	}
	t.FailNow()
	return 0 // so it compiles
}

func eval(first Ast, rest Ast) Ast {
	var lhs Ast = first.(NativeInt)
	if a, isArray := rest.(*NativeArray); isArray {
		for _, v := range a.Array {
			aFirstRest := v.(*NativeArray)
			if aFirstRest.Length() != 2 {
				panic("assert")
			}
			rhs := aFirstRest.Get(1).(NativeInt)
			op := aFirstRest.Get(0).(NativeString).Str
			lhs = NewNativeInt(ops[op](lhs.(NativeInt).Int(), rhs.Int()))
		}
	} else {
		panic("expected NativeArray, got " + rest.ContentString())
	}
	return lhs
}

// A working grammar for the calculator
// The ast nodes are special in that they will evaluate
// the numerical expression to eventually respond with an int
func grammar() *ast.Grammar {
	joeson := grammars.NewJoeson()
	CALC := []line.Line{
		o(Named("Input", "expr:Expression")),
		i(Named("Expression", "_ first:Term rest:( _ AddOp _ Term )* _"), func(it Ast) Ast {
			return eval(it.(NativeMap).Get("first"), it.(NativeMap).Get("rest"))
		}),
		i(Named("Term", "first:Factor rest:( _ MulOp _ Factor )*"), func(it Ast) Ast {
			return eval(it.(NativeMap).Get("first"), it.(NativeMap).Get("rest"))
		}),
		i(Named("Factor", "'(' expr:Expression ')' | integer:Integer"), func(it Ast) Ast {
			// alternation example
			var nm NativeMap = it.(NativeMap)
			if n, exists := nm.GetExists("integer"); exists {
				return n
			} else if expr, exists := nm.GetExists("expr"); exists {
				return expr
			} else {
				panic("assert")
			}
		}),
		i(Named("AddOp", "'+' | '-'")),
		i(Named("MulOp", "'*' | '/'")),
		i(Named("Integer", "/^-?[0-9]+/"), func(it Ast) Ast { return NewNativeIntFrom(it) }),
		i(Named("_", "[ \t]*")),
	}
	return line.NewGrammarFromLines("calc", CALC, joeson)
}

func assertResultIs(t *testing.T, sExpression string, nExpectedResult int) {
	t.Helper()
	if res, error := grammar().ParseString(sExpression); error == nil {
		fmt.Println(
			Cyan(sExpression),
			" --> ",
			Yellow(res.ContentString()),
			" --> ",
			BoldYellow(strconv.Itoa(extractResult(t, res))),
		)
		if extractResult(t, res) != nExpectedResult {
			t.Fail()
		}
	} else {
		t.Error(error)
	}
}

func Test_minus7(t *testing.T)         { assertResultIs(t, "-7", -7) }
func Test_73_plus_4(t *testing.T)      { assertResultIs(t, "73 + 4", 77) }
func Test_73_plus_minus4(t *testing.T) { assertResultIs(t, "73 +(-4)", 69) }
func Test_36(t *testing.T)             { assertResultIs(t, "4  *( (2 +1 )*3 )", 36) }
func Test_12(t *testing.T)             { assertResultIs(t, "-4 * ((-2+1) *3)", 12) }
func Test_minus11849(t *testing.T)     { assertResultIs(t, "241+513* -24 +((1934-192*2)/7)+1", -11849) }
func Test_many(t *testing.T) {
	gm := grammar()
	var h = map[string]int{
		"60/6/5":                2,
		"1 + 2 + 3 + 4 + 5 + 6": 21,
		"1 - 2 + 3 - 4 + 5 - 6": -3,
		"0 + 1":                 1,
		"0 - 1":                 -1,
		"0 * 1":                 0,
		"0 / 1":                 0,
	}
	for s, nExpectedResult := range h {
		res, error := gm.ParseString(s)
		if error == nil {
			if extractResult(t, res) != nExpectedResult {
				t.Fail()
			}
		} else {
			t.Error(error)
		}
	}
}

func Test_failing(t *testing.T) {
	gm := grammar()
	var h = map[string]string{
		"90 (6090)": "Error parsing at char:3",
		"-(7)":      "Error parsing at char:0",
	}
	for s, sExpectedError := range h {
		res, error := gm.ParseString(s)
		if error == nil {
			t.Error("expected error but got none, for: " + s + ". Res: " + res.ContentString())
		} else if strings.Index(error.Error(), sExpectedError) == 0 {
			// expected
		} else {
			t.Error("expected error " + sExpectedError + " for " + s + " but got " + error.Error() + " instead")
		}
	}
}
