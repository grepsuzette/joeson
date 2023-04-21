package main

import (
	"fmt"
	"testing"

	// "strings"
	"github.com/grepsuzette/joeson/helpers"
)

type exprExpectation struct {
	expr        string
	expectation string
}

func duo(s string, s2 string) exprExpectation {
	// return exprExpectation{strings.Replace(s, "\n", "", -1), s2}
	return exprExpectation{s, s2}
}

const (
	TRUE  = "1.000000"
	FALSE = "0.000000"
)

var tests = []exprExpectation{
	// ------------ arith -----------------------------
	duo("(+ 1 1)", "2.000000"),
	duo("(+ 1 1 1)", "3.000000"),
	duo("(+ 1 1 1 1 1 1 1 1)", "8.000000"),
	duo("(+ (+ 1 1) 1)", "3.000000"),
	duo(" ( + ( + 1 1 ) 1 ) ", "3.000000"),
	duo("(+ (+ 1 1) (+ 1 1))", "4.000000"),
	duo("(- (+ 1 1) (+ 1 1))", "0.000000"),
	duo("(- 10 12)", "-2.000000"),
	duo("(- 12 (- 2 1))", "11.000000"),
	duo("(% 12 5)", "2.000000"),

	// ------------ car, cdr -----------------------------
	duo("(car (10 9 8 7 6 5))", "10.000000"),
	duo("(cdr (10 9 8))", "(9.000000 8.000000)"),
	duo("(car (cdr (10 9 8)))", "9.000000"), // TODO

	// ------------ cmp, logic ----------------------------- // {{{1
	duo("(eq? 1 2)", FALSE),
	duo("(== 2 2)", TRUE),
	duo("(< 1 10)", TRUE),
	duo("(<= -10 -1)", TRUE),
	duo("(< 1 2 3 4 5)", TRUE),
	duo("(< 10 1)", FALSE),
	duo("(< -1 -10)", FALSE),
	duo("(lt 1 2 1 4 5)", FALSE),
	duo("(>= 5 4 4 3 2 1 0 -1 -1000)", TRUE),
	duo("(and 1 1)", TRUE),
	duo("(and (> 9 8 7) ( <= 1 1 2 2 3 3 4 ) )", TRUE),
	duo("(and (eq 7 7) (neq 7 7))", FALSE),
	duo("(or (eq 7 7) (neq 7 7))", TRUE),
	// }}}1

	// ------------ type --------------------
	duo("(list? 3)", FALSE),
	duo("(list? (1 2 3))", TRUE),

	// ------------ aliases  -----------------------------
	duo("(alias foo add)", "()"),
	duo("(foo 1 2 3)", "6.000000"),

	// ------------ user-defined  -----------------------------
	duo("(define (double x) (+ x x))", "()"),
	duo("(double 15)", "30.000000"),
	duo("(define (quadruple x) (double (double x)))", "()"),
	duo("(quadruple 15)", "60.000000"),
	duo("(define (cadr l) (car (cdr l)))", "()"),
	duo(`(cadr ("a" "b" "c"))`, `"b"`),
	duo("(define (caddr l) (car (cdr (cdr l))))", "()"),
	duo(`(caddr ("a" "b" "c"))`, `"c"`),

	// ------------ alternation -------------
	duo(`(if (== 4 4) "ok" "nok")`, `"ok"`),
	duo(`(if (!= 4 4) "ok" "nok")`, `"nok"`),
	duo(`(define (codeSize l)
		 (cond
		   ((>= l 9) "L")
		   ((>= l 5) "M")
		   (else "S")
		 ))`, `()`),
	duo(`(codeSize 10)`, `"L"`),
	duo(`(codeSize 6)`, `"M"`),
	duo(`(codeSize 2)`, `"S"`),

	// ----------- some fun -------------
	duo(`(define (divisible? n m) (== (% n m) 0))`, `()`),
	duo(`(divisible? 1000000 10)`, TRUE),
	duo(`(divisible? 1000000 7)`, FALSE),
	duo(`(define (fact n) (if (<= n 0) 1 (* n (fact (- n 1)))))`, `()`),
	duo(`(fact 0)`, `1.000000`),
	duo(`(fact 4)`, `24.000000`),
	duo(`(fact 5)`, `120.000000`),
}

func Test(t *testing.T) {
	gm := grammar()
	m := NewMachine()
	// simply compare m.Eval(gm.ParseString(`k`)) with `v`.
	for _, o := range tests {
		k := o.expr
		v := o.expectation
		t.Run(fmt.Sprintf("eval %s expected to give %s", k, v), func(t *testing.T) {
			if ast, e := gm.ParseString(k); e == nil {
				s := helpers.StripAnsi(m.Eval(ast.(Expr)).ContentString())
				if s != v {
					t.Errorf("%s expected to eval as %s gave %s instead\n", k, v, s)
				}
			} else {
				t.Errorf("%s expect to eval as %s did not even parse! error = %s\n", k, v, e)
			}
		})
	}
}

// vim: fdm=marker fdl=0
