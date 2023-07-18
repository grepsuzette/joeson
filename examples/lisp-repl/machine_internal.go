package main

import (
	"fmt"

	j "github.com/grepsuzette/joeson"
)

// for internal usage

func numberFromExpr(m Machine, item Expr) float64 {
	switch item.Kind {
	case kindNumber:
		return item.Number
	case kindString:
		panic("m.numberFromExpr should not have string arg")
	case kindOperator:
		panic("m.numberFromExpr should have numbers or list only")
	case kindList:
		return numberFromExpr(m, m.Eval(item))
	}
	panic(E)
}

func boolFromExpr(m Machine, item Expr) bool { return numberFromExpr(m, item) != 0. }

// TODO change the name of this one
func apply(f func(Machine, List) List) func(Machine, List) Expr {
	return func(m Machine, list List) Expr {
		var a List = f(m, list)
		return Expr{&j.Attributes{}, kindList, "", 0, a, ""}
	}
}

// (( x y z )) -> (x y z)
// (( + 1 1 1 )) -> 3
// Contract: rest is a List of one element
// Contract: that element is a List
// Contract: that List is returned
func unnestListEval(m Machine, rest List) Expr {
	if rest.Length() != 1 {
		fmt.Println(rest.String())
		fmt.Println(rest.Length())
		panic(E)
	}
	switch rest.List[0].Kind {
	case kindList:
		switch ((rest.List[0]).List.List)[0].Kind {
		case kindOperator:
			return m.Eval(rest.List[0])
		default:
		}
		return rest.List[0]
	default:
		fmt.Println(rest.List[0])
		panic("unexpected kind, expected List")
	}
}

func cmpNum(m Machine, a List, pred func(x, y float64) bool) Expr {
	if a.Length() < 2 {
		panic("needs at least 2 args")
	}
	x := numberFromExpr(m, a.List[0])
	for _, v := range a.List[1:] {
		y := numberFromExpr(m, v)
		if !pred(x, y) {
			return False()
		}
		x = y
	}
	return True()
}
