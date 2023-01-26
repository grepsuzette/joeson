package main

import (
	"fmt"
	"grepsuzette/joeson"
	"strconv"
	"strings"
	"testing"
)

// this example is the typical example with a calculator
// and is inspired from the mna-pigeon example

func xx(it joeson.Ast) joeson.Ast {
	return eval(it.(joeson.NativeMap).Get("first"), it.(joeson.NativeMap).Get("rest"))
}
func i(a ...any) joeson.ILine { return joeson.I(a...) }
func o(a ...any) joeson.OLine { return joeson.O(a...) }
func named(name string, lineStringOrAst any) joeson.NamedRule {
	return joeson.Named(name, lineStringOrAst)
}

var linesCalc = []joeson.Line{
	o(named("Input", "expr:Expression")),
	i(named("Expression", "_ first:Term rest:( _ AddOp _ Term )* _"), xx),
	i(named("Term", "first:Factor rest:( _ MulOp _ Factor )*"), xx),
	i(named("Factor", "'(' expr:Expression ')' | integer:Integer"), func(it joeson.Ast) joeson.Ast {
		// --- example of an alternation ------------
		var nm joeson.NativeMap = it.(joeson.NativeMap)
		if n, exists := nm.GetExists("integer"); exists {
			return n
		} else if expr, exists := nm.GetExists("expr"); exists {
			return expr
		} else {
			panic("assert")
		}
	}),
	i(named("AddOp", "'+' | '-'")),
	i(named("MulOp", "'*' | '/'")),
	i(named("Integer", "/^-?[0-9]+/"), func(it joeson.Ast) joeson.Ast { return joeson.NewNativeIntFrom(it) }),
	i(named("_", "[ \t]*")),
}

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

// extract the "expr" key from a result `x` to an int
// or, failing to do that, call FailNow()
func extractResult(t *testing.T, x joeson.Ast) int {
	t.Helper()
	if n, exists := x.(joeson.NativeMap).GetIntExists("expr"); exists {
		return n
	} else {
		t.Errorf("Failed to find a result like NativeMap{expr:<INT>} in " + x.ContentString())
	}
	t.FailNow()
	return 0 // so it compiles
}

func eval(first joeson.Ast, rest joeson.Ast) joeson.Ast {
	var lhs joeson.Ast = first.(joeson.NativeInt)
	if a, isArray := rest.(*joeson.NativeArray); isArray {
		for _, v := range a.Array {
			aFirstRest := v.(*joeson.NativeArray)
			if aFirstRest.Length() != 2 {
				panic("assert")
			}
			rhs := aFirstRest.Get(1).(joeson.NativeInt)
			op := aFirstRest.Get(0).(joeson.NativeString).Str
			lhs = joeson.NewNativeInt(ops[op](lhs.(joeson.NativeInt).Int(), rhs.Int()))
		}
	} else {
		panic("expected NativeArray, got " + rest.ContentString())
	}
	return lhs
}

// A working grammar for the calculator
// The ast nodes are special in that they will evaluate
// the numerical expression to eventually respond with an int
func grammar() *joeson.Grammar {
	return joeson.GrammarFromLines(linesCalc, "calc")
}

const esc string = ""
const reset string = esc + "[0m"

func cyan(s string) string       { return esc + "[36m" + s + reset }
func yellow(s string) string     { return esc + "[33m" + s + reset }
func boldYellow(s string) string { return esc + "[1;33m" + s + reset }

func assertResultIs(t *testing.T, sExpression string, nExpectedResult int) {
	t.Helper()
	if res, error := grammar().ParseString(sExpression); error == nil {
		fmt.Println(
			cyan(sExpression),
			" --> ",
			yellow(res.ContentString()),
			" --> ",
			boldYellow(strconv.Itoa(extractResult(t, res))),
		)
		if extractResult(t, res) != nExpectedResult {
			t.Fail()
		}
	} else {
		t.Error(error)
	}
}

func Test_73_plus_4(t *testing.T)      { assertResultIs(t, "73 + 4", 77) }
func Test_minus7(t *testing.T)         { assertResultIs(t, "-7", -7) }
func Test_73_plus_minus4(t *testing.T) { assertResultIs(t, "73 +(-4)", 69) }
func Test_36(t *testing.T)             { assertResultIs(t, "4  *( (2 +1 )*3 )", 36) }
func Test_12(t *testing.T)             { assertResultIs(t, "-4 * ((-2+1) *3)", 12) }
func Test_minus11849(t *testing.T)     { assertResultIs(t, "241+513* -24 +((1934-192*2)/7)+1", -11849) }
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

func Test_calc(t *testing.T) {
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
