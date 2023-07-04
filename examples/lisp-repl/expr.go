package main

import (
	"fmt"
	"strings"
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
		Kind     exprKind
		String   string
		Number   float64
		List     List
		Operator string
	}
)

func empty() Expr              { return Expr{kindList, "", 0, list(), ""} }
func number(f float64) Expr    { return Expr{kindNumber, "", f, nil, ""} }
func str(s string) Expr        { return Expr{kindString, s, 0, nil, ""} }
func operator(fun string) Expr { return Expr{kindOperator, "", 0, nil, fun} }
func True() Expr               { return number(1) }
func False() Expr              { return number(0) }
func Bool(b bool) Expr {
	if b {
		return True()
	} else {
		return False()
	}
}

func (o Expr) MustString() string {
	switch o.Kind {
	case kindString:
		return o.String
	case kindOperator:
		panic("Expected a string, got an operator instead: " + o.ContentString() + ". Did you mean MustStringOrOperator()?")
	default:
		fmt.Println(o)
		panic(E)
	}
}

func (o Expr) MustStringOrOperator() string {
	switch o.Kind {
	case kindString:
		return o.String
	case kindOperator:
		// fmt.Println("warn: Expected a string, got an operator instead: " + o.ContentString() + ". Using as string.")
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
		panic("Expected a List, got: " + o.ContentString())
	}
}

func (o Expr) ContentString() string {
	switch o.Kind {
	case kindString:
		return quoted(o.String)
	case kindNumber:
		return fmt.Sprintf(bold_magenta("%f"), o.Number)
	case kindOperator:
		return bold_cyan(o.Operator)
	case kindList:
		a := []string{}
		for _, expr := range o.List {
			a = append(a, expr.ContentString()) // beware, cycles will produce infinite loops here
		}
		return colorParen("(") + strings.Join(a, " ") + colorParen(")")
	}
	panic("Expr.ContentString() missing case")
}
