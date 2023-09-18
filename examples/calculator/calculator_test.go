package main

import (
	"strings"
	"testing"

	j "github.com/grepsuzette/joeson"
)

// A classical example for PEG parser is that
// of a calculator. It can parse and calculate
// expressions such as "((5 + 6 * 7)-1) / 2".

// named() creates a rule with a name
// :E.g. o(named("Alpha", `'a'`)) // named rule, item must match 'a'
// :     o(`'a'`)                 // the same rule, anonymous
func named(name string, v any) j.NamedRule { return j.Named(name, v) }

// to make i and o rules
func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }

var (
	rules = []j.Line{
		o(named("Expr", `_ first:Term rest:( _ AddOp _ Term )* _`), xx),
		i(named("Term", `first:Factor rest:( _ MulOp _ Factor )*`), xx),
		i(named("Factor", `'(' Expr ')' | Integer`)),
		i(named("AddOp", `'+' | '-'`)),
		i(named("MulOp", `'*' | '/'`)),
		i(named("_", `[ \t]*`)),
		i(named("Integer", `/^-?[0-9]+/`), func(it j.Ast) j.Ast {
			return j.NativeIntFrom(it)
		}),
	}
	ops = map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}
)

func xx(it j.Ast) j.Ast {
	return eval(
		it.(*j.NativeMap).Get("first"),
		it.(*j.NativeMap).Get("rest"),
	)
}

func eval(first j.Ast, rest j.Ast) j.Ast {
	operand1 := first.(j.NativeInt).Int()
	a := rest.(*j.NativeArray)
	// compute iteratively
	for _, v := range *a {
		// [first [+ x] [+ y] [- z] ...]
		//         ^^^ 2 elements by 2 elements
		pair := v.(*j.NativeArray)
		if pair.Length() != 2 {
			panic("assert")
		}
		op := pair.Get(0).(j.NativeString).String()
		operand2 := pair.Get(1).(j.NativeInt).Int()
		operand1 = ops[op](operand1, operand2)
	}
	return j.NewNativeInt(operand1)
}

func Test_calc(t *testing.T) {
	gm := j.GrammarFromLines("calc", rules)
	h := map[string]any{
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
		"90 (6090)":                        "Error parsing at char:3",
		"-(7)":                             "Error parsing at char:0",
	}
	for s, expectedAny := range h {
		res := gm.ParseString(s)
		switch expected := expectedAny.(type) {
		case int:
			// expect success (int)
			if j.IsParseError(res) {
				t.Error(res.String())
			} else {
				if res.(j.NativeInt).Int() != expected {
					t.Fail()
				}
			}
		case string:
			// expect failure (string)
			if !j.IsParseError(res) {
				t.Errorf("expected error but got none, for: %q. Res: %s",
					s,
					res.String(),
				)
			} else {
				sError := res.(j.ParseError).ErrorString
				if strings.Index(sError, expected) == 0 {
					// no worry, expected case
				} else {
					t.Errorf("expected error %q for %s but got %s instead",
						expected,
						s,
						sError,
					)
				}
			}
		default:
			panic("assert")
		}
	}
}
