package main

import (
	"fmt"
	"math"
)

// (car (1 2 3)) is 1
func car(m Machine, rest List) Expr {
	expr := unnestListEval(m, rest)
	if expr.Kind != kindList {
		panic("car must operates on a list, got " + expr.String())
	}
	list := expr.MustList()
	if list.Length() > 0 {
		return list.List[0]
	} else {
		panic("car of empty list")
	}
}

// (cdr (1 2 3)) is (2 3)
// (cdr (cdr (2 3))) is (3)
func cdr(m Machine, rest List) List {
	expr := unnestListEval(m, rest)
	if expr.Kind != kindList {
		panic("cdr must operates on a list, got " + expr.String())
	}
	l := expr.MustList()
	if l.Length() <= 1 {
		return list()
	} else {
		return list(l.List[1:]...)
	}
}

func add(m Machine, a List) Expr {
	r := 0.
	for _, item := range a.List {
		r += numberFromExpr(m, item)
	}
	return number(r)
}

func sub(m Machine, a List) Expr {
	if a.Length() < 2 {
		panic("not enough args for sub (-). Needs at least 2")
	} else {
		r := numberFromExpr(m, a.List[0])
		for _, item := range a.List[1:] {
			r -= numberFromExpr(m, item)
		}
		return number(r)
	}
}

func mul(m Machine, a List) Expr {
	r := 1.
	for _, item := range a.List {
		r *= numberFromExpr(m, item)
	}
	return number(r)
}

func remainder(m Machine, a List) Expr {
	if a.Length() != 2 {
		panic("% (remainder) expects 2 args")
	} else {
		return number(math.Mod(numberFromExpr(m, a.List[0]), numberFromExpr(m, a.List[1])))
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
	for _, v := range a.List {
		if !boolFromExpr(m, v) {
			return False()
		}
	}
	return True()
}

// if there is one true, it's true. Otherwise false. Even (or ()) is false.
func or(m Machine, a List) Expr {
	for _, v := range a.List {
		if boolFromExpr(m, v) {
			return True()
		}
	}
	return False()
}

// (not (eq 7)) is like (neq 7). Contract: exactly 1 argument.
func not(m Machine, rest List) Expr {
	if rest.Length() != 1 {
		panic("`not` admits only 1 arg")
	}
	return Bool(boolFromExpr(m, rest.List[0]))
}

// (if <predicate> <then-expr> <else-expr>)
// <predicate> is evaluated as boolean,
func _if(m Machine, rest List) Expr {
	if rest.Length() != 3 {
		panic("if expects exactly 3 args")
	}
	pred := rest.List[0]
	thenexpr := rest.List[1]
	elseexpr := rest.List[2]
	if boolFromExpr(m, pred) {
		return m.Eval(thenexpr)
	} else {
		return m.Eval(elseexpr)
	}
}

// (cond
//
//	((>= x 9) <action1>)
//	((>= x 5) <action2>)
//	(else <action3>)
//
// )
// note `else` will simply be substituted by a true condition
// an empty list is returned if no branch is matched.
func cond(m Machine, rest List) Expr {
	for _, branch := range rest.List {
		a := branch.MustList()
		if a.Length() != 2 {
			panic("(cond (<pred> <action>) ...): each branch requires 2 args, but branch is " + branch.String())
		}
		pred := a.List[0]
		action := a.List[1]
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
	if rest.Length() != 1 {
		panic("list? requires 1 arg")
	}
	return Bool(rest.List[0].Kind == kindList)
}

// Define function aliases. Modifies the machine. Returns empty().
// An alias must point to an operator name that is alrady appearing in m.funcs.
// E.g. in lisp: `alias == eq?`.
// E.g. in go: alias(m, list("==", "eq?"))
// Then (== 4 2) will be executed as (eq? 4 2)
func alias(m Machine, rest List) Expr {
	a := rest
	if a.Length() != 2 {
		panic("alias arg1 arg2")
	}
	alias := a.List[0].MustStringOrOperator()
	aliased := a.List[1].MustStringOrOperator()
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
//
//	a[0] is a list with: the new function name, followed by named arguments.
//	a[1] is an expression, that can contain symbols declared as args in a[0].
func define(m Machine, a List) Expr {
	syntax := "define (funName [arg1 [arg2 [...]]]) (expr)"
	if a.Length() != 2 {
		panic(syntax)
	}
	def := a.List[0].MustList()
	if def.Length() < 1 {
		panic(syntax + ": missing funName")
	}
	funName := def.List[0].MustStringOrOperator()
	var argNames []string
	for _, x := range def.List[1:] {
		argNames = append(argNames, x.MustStringOrOperator())
	}
	expr := a.List[1]
	delete((*m.aliases), funName) // Contract: If an alias exists with this name, it will be deleted before.
	// when funName is called, we are to return evaluation of expr,
	// substituting named parameters
	(*m.funcs)[funName] = func(_ Machine, rest List) Expr {
		// rest is given by the caller.
		// first need to check if number of given arguments is that of argNames
		if rest.Length() != len(argNames) {
			panic(fmt.Sprintf("%s expecting args %v but got (%v)\n", funName, argNames, rest))
		}
		args := map[string]Expr{}
		for i := range argNames {
			args[argNames[i]] = rest.List[i]
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
		var us []Expr
		for _, subexpr := range expr.List.List {
			us = append(us, substArgsInExpr(args, subexpr))
		}
		return Expr{attr{}, kindList, "", 0, list(us...), ""}
	default:
		panic("unhandled kind")
	}
}
