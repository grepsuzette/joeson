package main

import (
	"fmt"
	"strings"

	j "github.com/grepsuzette/joeson"
)

// -- uLisp AST

type exprKind int

const (
	kindString   exprKind = 0
	kindNumber            = 1
	kindList              = 2
	kindOperator          = 3 // such as a func name
)

type (
	Expr struct {
		*j.Attr
		Kind     exprKind
		Str      string
		Number   float64
		List     List
		Operator string
	}
)

func empty() Expr              { return Expr{j.NewAttr(), kindList, "", 0, list(), ""} }
func number(f float64) Expr    { return Expr{j.NewAttr(), kindNumber, "", f, nilList(), ""} }
func str(s string) Expr        { return Expr{j.NewAttr(), kindString, s, 0, nilList(), ""} }
func operator(fun string) Expr { return Expr{j.NewAttr(), kindOperator, "", 0, nilList(), fun} }
func True() Expr               { return number(1) }
func False() Expr              { return number(0) }
func Bool(b bool) Expr {
	if b {
		return True()
	} else {
		return False()
	}
}

func (o Expr) assertNode() {}

func (o Expr) MustString() string {
	switch o.Kind {
	case kindString:
		return o.Str
	case kindOperator:
		panic("Expected a string, got an operator instead: " + o.String() + ". Did you mean MustStringOrOperator()?")
	default:
		fmt.Println(o)
		panic(E)
	}
}

func (o Expr) MustStringOrOperator() string {
	switch o.Kind {
	case kindString:
		return o.Str
	case kindOperator:
		// fmt.Println("warn: Expected a string, got an operator instead: " + o.String() + ". Using as string.")
		return o.Operator
	default:
		fmt.Println(o)
		panic(E)
	}
}

func (o Expr) MustList() List {
	switch o.Kind {
	case kindList:
		return o.List
	default:
		panic("Expected a List, got: " + o.String())
	}
}

func (o Expr) String() string {
	switch o.Kind {
	case kindString:
		return quoted(o.Str)
	case kindNumber:
		return fmt.Sprintf(j.BoldMagenta("%f"), o.Number)
	case kindOperator:
		return j.BoldCyan(o.Operator)
	case kindList:
		a := []string{}
		for _, expr := range o.List.List {
			a = append(a, expr.String()) // beware, cycles will produce infinite loops here
		}
		return colorParen("(") + strings.Join(a, " ") + colorParen(")")
	}
	panic("Expr.String() missing case")
}
