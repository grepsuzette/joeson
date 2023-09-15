package main

import (
	"strings"
	"testing"

	j "github.com/grepsuzette/joeson"
)

// this example is the typical example with a calculator
// and is inspired from the mna-pigeon example
//
// To understand more easily, start from the bottom, with Test_calc()

// "Inline" line of rule AKA ILine. Inline rules are always named().
// An inline rule can be referenced by its name. When it isn't
// it is totally inert and inactive.
func i(a ...any) j.ILine { return j.I(a...) }

// "OR" rule. Inside a rank, "OR" rules (AKA OLine) are parsed one after the
// other until one returns something other than nil. Some of them are named,
// but they usually aren't, as it's more the point of an ILine to be
// referenced.
func o(a ...any) j.OLine { return j.O(a...) }

// A Key-value pair, where Key is the name.
// This is exclusively used with joeson ILine and OLine to name things.
func named(name string, lineStringOrAst any) j.NamedRule {
	return j.Named(name, lineStringOrAst)
}

var linesCalc = []j.Line{
	o(named("Input", "expr:Expression")),
	i(named("Expression", "_ first:Term rest:( _ AddOp _ Term )* _"), xx),
	i(named("Term", "first:Factor rest:( _ MulOp _ Factor )*"), xx),
	i(named("Factor", "'(' expr:Expression ')' | integer:Integer"), func(it j.Ast) j.Ast {
		// --- example of an alternation ------------
		nm := it.(*j.NativeMap)
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
	i(named("Integer", "/^-?[0-9]+/"), func(it j.Ast) j.Ast {
		return j.NativeIntFrom(it)
	}),
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
func xx(it j.Ast) j.Ast {
	return eval(
		it.(*j.NativeMap).Get("first"),
		it.(*j.NativeMap).Get("rest"),
	)
}

// extract the "expr" key from a result `x` to an int
// or, failing to do that, call FailNow()
func extractResult(t *testing.T, x j.Ast) int {
	t.Helper()
	if n, exists := x.(*j.NativeMap).GetIntExists("expr"); exists {
		return n
	} else {
		t.Errorf("Failed to find a result like NativeMap{expr:<INT>} in %s",
			x.String(),
		)
	}
	t.FailNow()
	return 0 // so it compiles
}

func eval(first j.Ast, rest j.Ast) j.Ast {
	var lhs j.Ast = first.(j.NativeInt)
	if a, isArray := rest.(*j.NativeArray); isArray {
		for _, v := range *a {
			aFirstRest := v.(*j.NativeArray)
			if aFirstRest.Length() != 2 {
				panic("assert")
			}
			rhs := aFirstRest.Get(1).(j.NativeInt)
			op := string(aFirstRest.Get(0).(j.NativeString))
			lhs = j.NewNativeInt(ops[op](lhs.(j.NativeInt).Int(), rhs.Int()))
		}
	} else {
		panic("expected NativeArray, got " + rest.String())
	}
	return lhs
}

func Test_failing(t *testing.T) {
	gm := j.GrammarFromLines("calc", linesCalc)
	h := map[string]string{
		"90 (6090)": "Error parsing at char:3",
		"-(7)":      "Error parsing at char:0",
	}
	for s, sExpectedError := range h {
		res := gm.ParseString(s)
		if !j.IsParseError(res) {
			t.Errorf("expected error but got none, for: %q. Res: %s",
				s,
				res.String(),
			)
		} else {
			sError := res.(j.ParseError).ErrorString
			if strings.Index(sError, sExpectedError) == 0 {
				// np, expected case
			} else {
				t.Errorf("expected error %q for %s but got %s instead",
					sExpectedError,
					s,
					sError,
				)
			}
		}
	}
}

func Test_calc(t *testing.T) {
	gm := j.GrammarFromLines("calc", linesCalc)
	h := map[string]int{
		"0 + 1":                            1,
		"0 - 1":                            -1,
		"0 * 1":                            0,
		"0 / 1":                            0,
		"-7":                               -7,
		"73 + 4":                           77,
		"73 +(-4)":                         69,
		"-4 * ((-2+1) *3)":                 12,
		"1 + 2 + 3 + 4 + 5 + 6":            21,
		"1 - 2 + 3 - 4 + 5 - 6":            -3,
		"241+513* -24 +((1934-192*2)/7)+1": -11849,
		"60/6/5":                           2,
		"4*((2+1) * 3 )":                   36,
		"-1-2-3":                           -6,
	}
	for s, nExpectedResult := range h {
		res := gm.ParseString(s)
		if j.IsParseError(res) {
			t.Error(res.String())
		} else {
			if extractResult(t, res) != nExpectedResult {
				t.Fail()
			}
		}
	}
}
