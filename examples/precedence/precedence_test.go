package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

// A study on precedence.
//
// strung := "a + b"
// expect := "[+ a b]"
//
// strung := "a + b - c"
// expect := "[- [+ a b] c]"
//
// strung := "a*m/n + x - y + z + b*c"
// expect := "[+ [+ [- [+ [/ [* a m] n] x] y] z] [* b c]]"
//
// Original joeson has this mention:
//
// Some known limitations of this parser;
// Parses '3-2-1' right-recursively, which is unexpected.
// [source](http://tratt.net/laurie/research/publications/html/tratt__direct_left_recursive_parsing_expression_grammars/)
// ```
//
//	START:   "EXPR"
//	EXPR:
//	  EXPR1: "EXPR '-' EXPR"
//	  EXPR2: "NUM"
//	NUM:     "/[0-9]+/"
//
// This is a study to show what can be done.
type (
	operation struct {
		*j.Attr
		operator string
		operand  j.Ast
	}
	binaryexpr struct {
		*j.Attr
		Lx j.Ast
		Op string
		Rx j.Ast
	}
)

func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }
func named(name string, lineStringOrAst any) j.NamedRule {
	return j.Named(name, lineStringOrAst)
}

// growMoss consumes a string of operations of the same precedence
// and their operands, as in `a + b - c + d - e`. It returns a tree of
// binary expressions, e.g.: `[- [+ [- [+ a b] c] d] e]`.
//
// `it` is the product of rules such as `T (Op T)*`,
// which means `it` is a NativeArray of size 2.
// it.Get(0) is described as `first` below,
// it.Get(1) is a NativeArray such as [[+ a], [- b], [+ c], [- d]]
//
// Each level of precedence must call this function in their own rule.
// Then something like "a*m/n + x - y + z + b*c", with +- and */ being
// 2 different precendence level, must give
// "[+ [+ [- [+ [/ [* a m] n] x] y] z] [* b c]]"
//
// The moss grows laterally:
// fMoss(first op1 e1 op2 e2 op3 e3 ...) grows like this:
//
// :	       [op1 first e1]
// :  	  [op2 [op1 first e1] e2]
// : [op3 [op2 [op1 first e1] e2] e3]       etc.
func growMoss(it j.Ast) j.Ast {
	switch v := it.(type) {
	case j.NativeString:
		return v
	}

	// extract `first` and `operations` from `it`
	first := it.(*j.NativeArray).Get(0)
	operations := []operation{}
	for _, v := range it.(*j.NativeArray).Get(1).(*j.NativeArray).Array() {
		a := v.(*j.NativeArray).Array()
		operations = append(operations, operation{
			j.NewAttr(),
			a[0].String(),
			a[1],
		})
	}
	// no operations means result simply is `first`
	if len(operations) == 0 {
		return first
	}
	// The moss grows laterally:
	//
	//           [op1 first e1]
	//      [op2 [op1 first e1] e2]
	// [op3 [op2 [op1 first e1] e2] e3]       etc.
	moss := first
	for _, operation := range operations {
		moss = binaryexpr{
			j.NewAttr(),
			moss,
			operation.operator,
			operation.operand,
		}
	}
	return moss
}

func (o operation) String() string {
	return o.operator + o.operand.String()
}

func (be binaryexpr) String() string {
	return "[" + be.Op + " " + be.Lx.String() + " " + be.Rx.String() + "]"
}

func assertWorks(t *testing.T, strung, expect string, rules []j.Line) {
	t.Helper()
	gm := j.GrammarFromLines("at hand", rules)
	r := gm.ParseString(strung).String()
	if r != expect {
		t.Errorf("%q should have parsed as %q, obtained %q instead\n", strung, expect, r)
	}
}

func assertFails(t *testing.T, strung, expect string, rules []j.Line) {
	t.Helper()
	gm := j.GrammarFromLines("at hand", rules)
	r := gm.ParseString(strung).String()
	if r == expect {
		t.Errorf("%q has parsed as %q, we expected a fail with those rules\n", strung, expect)
	}
}

func assert321Works(t *testing.T, rules []j.Line) {
	t.Helper()
	assertWorks(t, "3-2-1", "[- [- 3 2] 1]", rules)
}

func assert321Fails(t *testing.T, rules []j.Line) {
	t.Helper()
	assertFails(t, "3-2-1", "[- [- 3 2] 1]", rules)
}

// This test exposes the problem of strings like 3-2-1 parsing as [- 3 [- 2 1]]
// with certain grammars. It is a well-known problem aknowledged by original
// joeson implementation in doc/limitations.md, and described by Laurence Tratt
// in her 2010 paper "Direct left recursive PEG".
func TestPrecedence(t *testing.T) {
	// Everything seems good in this example.
	// It is inspired from the calculator example but
	// doesn't evaluate anything, keeping symbolic expressions.
	// Torturing the grammar however can show the problems explained by Laurence Tratt
	// in her paper "Direct left recursive PEG" and documented in
	// original joeson doc/limitations.md
	//
	// "3-2-1" will parse as [- [- 3 2] 1].
	assert321Works(t, []j.Line{
		o(`Expr`),
		i(named("Expr", `Term ( AddOp Term )*`), growMoss),
		i(named("Term", `Factor ( MulOp Factor )*`), growMoss),
		i(named("Factor", `'(' Expr ')' | Glyph`)),
		i(named("AddOp", `'+' | '-'`)),
		i(named("MulOp", `'*' | '/'`)),
		i(named("Glyph", `[a-zA-Z0-9]`)),
	})

	// Right-handside of first choice is a Right Recursion.
	// This fails, parsing "3-2-1" as [- 3 [- 2 1]] instead of [- [- 3 2] 1].
	// In other words, it parses as (3-(2-(1))) instead of (((3)-2)-1)
	assert321Fails(t, []j.Line{
		o(`Expr`),
		i(named("Expr", `Expr ( AddOp Expr )* | Glyph`), growMoss),
		i(named("AddOp", `'+' | '-'`)),
		i(named("Glyph", `[a-zA-Z0-9]`)),
	})

	// Changing Expr to Glyph in the first choice fixes the previous grammar.
	// It is as stated above in Laurie's paper RR in "EXPR" <- "Expr '-' Expr" is a problem.
	// But is this a generic enough solution?
	assert321Works(t, []j.Line{
		o(`Expr`),
		i(named("Expr", `Expr ( _ AddOp _ Glyph )* | Glyph`), growMoss),
		i(named("AddOp", `'+' | '-'`)),
		i(named("Glyph", `[a-zA-Z0-9]`)),
		i(named("_", `[ \t]*`)),
	})

	// Works. This is totally equivalent to the
	// previous grammar, but expressed in a more
	// realistic way (large grammars would choose to
	// separate choices in several O lines rather
	// than one very long, also allowing callbacks
	// for each alternation). Will use this notation below
	assert321Works(t, []j.Line{
		o(named("Expr", []j.Line{
			o(`Expr ( _ AddOp _ Glyph )*`, growMoss),
			o(`Glyph`, growMoss),
			i(named("AddOp", `'+' | '-'`)),
			i(named("Glyph", `[a-zA-Z0-9]`)),
			i(named("_", `[ \t]*`)),
		},
		)),
	})

	// Factor rule has deep recursion (of Expr), alternation.
	assert321Works(t, []j.Line{
		o(named("Expr", []j.Line{
			o(`Term ( _ AddOp _ Term )*`, growMoss),
			i(named("Term", `Factor ( _ MulOp _ Factor  )*`), growMoss),
			i(named("Factor", `'(' Expr ')' | Glyph`), growMoss),
			i(named("AddOp", `'+' | '-'`)),
			i(named("MulOp", `'*' | '/'`)),
			i(named("Glyph", `[a-zA-Z0-9]`)),
			i(named("_", `[ \t]*`)),
		},
		)),
	})

	// A more realistic example, with deep recursion, alternation,
	// both unary and binary expression allowing '-'. It still works.
	// This will reuse these `rules` for a few other tests below.
	rules := []j.Line{
		o(named("Expr", []j.Line{
			o(named("BinaryExpr", []j.Line{
				o(`Term ( _ AddOp _ Term )*`),
				i(named("Term", `Factor ( _ MulOp _ Factor  )*`), growMoss),
				i(named("Factor", `'(' Expr ')' | Glyph`), growMoss),
				i(named("AddOp", `'+' | '-'`)),
				i(named("MulOp", `'*' | '/'`)),
			}), growMoss),
			o(named("UnaryExpr", `unary_op Glyph`)),
			i(named("unary_op", `'!' | '^' | '-'`)),
			i(named("Glyph", `[a-zA-Z0-9]`)),
			i(named("_", `[ \t]*`)),
		},
		)),
	}
	assert321Works(t, rules)

	{
		strung := "a + b"
		expect := "[+ a b]"
		assertWorks(t, strung, expect, rules)
	}
	{
		strung := "a + b - c"
		expect := "[- [+ a b] c]"
		assertWorks(t, strung, expect, rules)
	}
	{
		strung := "a*m/n + x - y + z + b*c"
		expect := "[+ [+ [- [+ [/ [* a m] n] x] y] z] [* b c]]"
		assertWorks(t, strung, expect, rules)
	}

	// Show how to make it fail
	assert321Fails(t, []j.Line{
		o(named("Expr", []j.Line{
			o(named("BinaryExpr", []j.Line{
				o(`Term ( _ AddOp _ Expr )*`), // <- Indirect right recursion
				i(named("Term", `Factor ( _ MulOp _ Factor  )*`), growMoss),
				i(named("Factor", `'(' Expr ')' | Glyph`), growMoss),
				i(named("AddOp", `'+' | '-'`)),
				i(named("MulOp", `'*' | '/'`)),
			}), growMoss),
			o(named("UnaryExpr", `unary_op Glyph`)),
			i(named("unary_op", `'!' | '^' | '-'`)),
			i(named("Glyph", `[a-zA-Z0-9]`)),
			i(named("_", `[ \t]*`)),
		},
		)),
	})

	// in short, problems shown by Laurie in her paper
	// are real. Our recomm. is to read it, and work around it
	// as in this test, rather than to create a perfect parser.
}
