package main

import (
	"testing"

	j "github.com/grepsuzette/joeson"
)

// A study on precedence.
//
// This is based on example/calculator; however instead of
// evaluating a result, we focus on building a tree of binary
// expressions.
//
// strung := "a + b"
// expect := "[+ a b]"
//
// strung := "a + b - c"
// expect := "[- [+ a b] c]"
//
// strung := "a*m/n + x - y + z + b*c"
// expect := "[+ [+ [- [+ [/ [* a m] n] x] y] z] [* b c]]"
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

var linesCalc = []j.Line{
	o(`Expr`),
	// We have 2 level of precedence:
	i(named("Expr", `Term ( _ AddOp _ Term )*`), growMoss),
	i(named("Term", `Factor ( _ MulOp _ Factor )*`), growMoss),
	i(named("Factor", `'(' Expr ')' | Glyph`)),
	i(named("AddOp", `'+' | '-'`)),
	i(named("MulOp", `'*' | '/'`)),
	i(named("Glyph", `[a-zA-Z0-9]`)),
	i(named("_", `[ \t]*`)),
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
// +- is one specie of moss that grows laterally,
// */ is another specie of moss that grows laterally but on another level.
// The moss species don't intermix
//
// The moss grows laterally:
// fMoss(first op1 e1 op2 e2 op3 e3 ...) grows like this:
//
// :	       [op1 first e1]
// :  	  [op2 [op1 first e1] e2]
// : [op3 [op2 [op1 first e1] e2] e3]       etc.
func growMoss(it j.Ast) j.Ast {
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

func TestExpr(t *testing.T) {
	gm := j.GrammarFromLines("priorities", linesCalc)
	strung := "a*m/n + x - y + z + b*c"
	expect := "[+ [+ [- [+ [/ [* a m] n] x] y] z] [* b c]]"
	r := gm.ParseString(strung).String()
	if r != expect {
		t.Errorf("%q should have parsed as %q, obtained %q instead\n", strung, expect, r)
	}
}
