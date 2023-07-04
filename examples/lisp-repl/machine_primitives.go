package main

import (
	"fmt"
	"math"

	"github.com/grepsuzette/joeson/helpers"
)

// (car (1 2 3)) is 1
func car(m Machine, rest List) Expr {
	expr := unnestListEval(m, rest)
	if expr.Kind != kindList {
		panic("car must operates on a list, got " + expr.ContentString())
	}
	list := expr.MustList()
	if len(list) > 0 {
		return list[0]
	} else {
		panic("car of empty list")
	}
}

// (cdr (1 2 3)) is (2 3)
// (cdr (cdr (2 3))) is (3)
func cdr(m Machine, rest List) List {
	expr := unnestListEval(m, rest)
	if expr.Kind != kindList {
		panic("cdr must operates on a list, got " + expr.ContentString())
	}
	l := expr.MustList()
	if len(l) <= 1 {
		return list()
	} else {
		return l[1:]
	}
}

func add(m Machine, a List) Expr {
	r := 0.
	for _, item := range a {
		r += numberFromExpr(m, item)
	}
	return number(r)
}

func sub(m Machine, a List) Expr {
	if len(a) < 2 {
		panic("not enough args for sub (-). Needs at least 2")
	} else {
		r := numberFromExpr(m, a[0])
		for _, item := range a[1:] {
			r -= numberFromExpr(m, item)
		}
		return number(r)
	}
}

func mul(m Machine, a List) Expr {
	r := 1.
	for _, item := range a {
		r *= numberFromExpr(m, item)
	}
	return number(r)
}

func remainder(m Machine, a List) Expr {
	if len(a) != 2 {
		panic("% (remainder) expects 2 args")
	} else {
		return number(math.Mod(numberFromExpr(m, a[0]), numberFromExpr(m, a[1])))
	}
}

// we allow things like a < b < c < d < e < f
// (>= 5 4 4 3 2 1 0 -1 -1000) is true  etc
func eq(m Machine, a List) Expr  { return cmpNum(m, a, func(x, y float64) bool { return x == y }) }
func neq(m Machine, a List) Expr { return cmpNum(m, a, func(x, y float64) bool { return x != y }) }
func lt(m Machine, a List) Expr  { return cmpNum(m, a, func(x, y float64) bool { return x < y }) }
func le(m Machine, a List) Expr  { return cmpNum(m, a, func(x, y float64) bool { return x <= y }) }
func gt(m Machine, a List) Expr  { return cmpNum(m, a, func(x, y float64) bool { return x > y }) }
func ge(m Machine, a List) Expr  { return cmpNum(m, a, func(x, y float64) bool { return x >= y }) }

// if there is one false, it's false. Otherwise true. Even (and ()) is true.
// TODO (and (true true false) (false true true)) should be legit and would be like
// (and false false) meaning it evalutes to false. But it's TODO
// We just evaluate the (and a b c d ) form for now
func and(m Machine, a List) Expr {
	for _, v := range a {
		if !boolFromExpr(m, v) {
			return False()
		}
	}
	return True()
}

// if there is one true, it's true. Otherwise false. Even (or ()) is false.
func or(m Machine, a List) Expr {
	for _, v := range a {
		if boolFromExpr(m, v) {
			return True()
		}
	}
	return False()
}

// (not (eq 7)) is like (neq 7). Contract: exactly 1 argument.
func not(m Machine, rest List) Expr {
	if len(rest) != 1 {
		panic("`not` admits only 1 arg")
	}
	return Bool(boolFromExpr(m, rest[0]))
}

// (if <predicate> <then-expr> <else-expr>)
// <predicate> is evaluated as boolean,
func _if(m Machine, rest List) Expr {
	if len(rest) != 3 {
		panic("if expects exactly 3 args")
	}
	pred := rest[0]
	thenexpr := rest[1]
	elseexpr := rest[2]
	if boolFromExpr(m, pred) {
		return m.Eval(thenexpr)
	} else {
		return m.Eval(elseexpr)
	}
}

// (cond
//   ((>= x 9) <action1>)
//   ((>= x 5) <action2>)
//   (else <action3>)
// )
// note `else` will simply be substituted by a true condition
// an empty list is returned if no branch is matched.
func cond(m Machine, rest List) Expr {
	for _, branch := range rest {
		a := branch.MustList()
		if len(a) != 2 {
			panic("(cond (<pred> <action>) ...): each branch requires 2 args, but branch is " + branch.ContentString())
		}
		pred := a[0]
		action := a[1]
		if pred.Kind == kindOperator && pred.Operator == "else" {
			return m.Eval(action)
		} else if pred.Kind == kindList {
			// e.g. (>= x 438)
			// if boolFromExpr(m, unnestListEval(m, pred.List)) {
			// 	return m.Eval(action)
			// }
			if boolFromExpr(m, m.Eval(pred)) {
				return m.Eval(action)
			}
		} else {
			if boolFromExpr(m, pred) {
				return m.Eval(action)
			}
		}
	}
	return empty()
}

// (list? 43) -> false
// (list? (1 2 3)) -> true
// (list? (+ 1 2 3)) -> false
// TODO (list? '(+ 1 2 3)) -> true
func isList(m Machine, rest List) Expr {
	if len(rest) != 1 {
		panic("list? requires 1 arg")
	}
	return Bool(rest[0].Kind == kindList)
}

// Define function aliases. Modifies the machine. Returns empty().
// An alias must point to an operator name that is alrady appearing in m.funcs.
// E.g. in lisp: `alias == eq?`.
// E.g. in go: alias(m, list("==", "eq?"))
// Then (== 4 2) will be executed as (eq? 4 2)
func alias(m Machine, rest List) Expr {
	a := rest
	if len(a) != 2 {
		panic("alias arg1 arg2")
	}
	alias := a[0].MustStringOrOperator()
	aliased := a[1].MustStringOrOperator()
	if _, has := (*m.funcs)[aliased]; has {
		(*m.aliases)[alias] = aliased
	} else {
		panic(aliased + " is not a defined function")
	}
	return empty()
}

// Define functions. Modifies the machine. Returns empty().
// A function can be redefined (be careful).
// Contract: If an alias exists with this name, it will be deleted before.
// E.g. in lisp: `(define ("square" n) (* n n))`
// Thus:
//	a[0] is a list with: the new function name, followed by named arguments.
//	a[1] is an expression, that can contain symbols declared as args in a[0].
func define(m Machine, a List) Expr {
	syntax := "define (funName [arg1 [arg2 [...]]]) (expr)"
	if len(a) != 2 {
		panic(syntax)
	}
	def := a[0].MustList()
	if len(def) < 1 {
		panic(syntax + ": missing funName")
	}
	funName := def[0].MustStringOrOperator()
	argNames := helpers.AMap(def[1:], func(x Expr) string { return x.MustStringOrOperator() })
	expr := a[1]
	delete((*m.aliases), funName) // Contract: If an alias exists with this name, it will be deleted before.
	// when funName is called, we are to return evaluation of expr,
	// substituting named parameters
	(*m.funcs)[funName] = func(_ Machine, rest List) Expr {
		// rest is given by the caller.
		// first need to check if number of given arguments is that of argNames
		if len(rest) != len(argNames) {
			panic(fmt.Sprintf("%s expecting args %v but got (%v)\n", funName, argNames, rest))
		}
		args := map[string]Expr{}
		for i := range argNames {
			args[argNames[i]] = rest[i]
		}
		return m.Eval(substArgsInExpr(args, expr))
	}
	return empty()
}

func substArgsInExpr(args map[string]Expr, expr Expr) Expr {
	switch expr.Kind {
	case kindString, kindNumber:
		return expr
	case kindOperator:
		if v, has := args[expr.Operator]; has {
			return v
		} else {
			return expr
		}
	case kindList:
		return Expr{kindList, "", 0, helpers.AMap(expr.List, func(subexpr Expr) Expr {
			return substArgsInExpr(args, subexpr)
		}), ""}
	default:
		panic("unhandled kind")
	}
}
