package main

import (
	"fmt"
)

// machine with just a few built-in primitives

type (
	Funcdef func(Machine, List) Expr
	Machine struct {
		funcs   *map[string]Funcdef
		aliases *map[string]string // alias to a func
	}
)

func NewMachine() Machine {
	m := Machine{}
	// funcs is initialized with built-ins.
	// User defined funcs will also be dynamically added here.
	// Built-ins can be dynamically redefined by users.
	m.funcs = &map[string]Funcdef{
		"add":    add,
		"sub":    sub,
		"mul":    mul,
		"car":    car,
		"cdr":    apply(cdr),
		"alias":  alias,
		"define": define,
		"lt":     lt,
		"le":     le,
		"gt":     gt,
		"ge":     ge,
		"eq":     eq,
		"neq":    neq,
		"and":    and,
		"or":     or,
		"not":    not,
		"if":     _if,
		"cond":   cond,
		"list?":  isList,
		"%":      remainder,
	}
	m.aliases = &map[string]string{
		"+":    "add",
		"-":    "sub",
		"*":    "mul",
		"<":    "lt",
		">":    "gt",
		">=":   "ge",
		"<=":   "le",
		"==":   "eq",
		"eq?":  "eq",
		"neq?": "neq",
		"!=":   "neq",
	}
	return m
}

func (m Machine) Eval(expr Expr) Expr {
	switch expr.Kind {
	case kindString, kindNumber:
		return expr
	case kindOperator:
		// Operator may appear to make no sense if not at the start of a list,
		return expr // but anything goes in FP...
	case kindList:
		if len(expr.List) == 0 {
			return empty()
		}
		first := expr.List[0]
		switch first.Kind {
		case kindList, kindNumber, kindString:
			return expr
		case kindOperator:
			op := first.Operator
			// resolve alias
			if s, exists := (*m.aliases)[op]; exists {
				// TODO rewrite next line shorter
				if _, exists := (*m.funcs)[s]; exists {
					op = s
				} else {
					panic("aliases must point to functions: " + s)
				}
			}
			if f, exists := (*m.funcs)[op]; exists {
				return f(m, expr.List[1:])
			} else {
				funcs := []string{}
				aliases := []string{}
				for k := range *m.funcs {
					funcs = append(funcs, k)
				}
				for k := range *m.aliases {
					aliases = append(aliases, k)
				}
				panic(fmt.Sprintf("undefined function '%s'. Defined functions: %v. Defined aliases: %v ", op, funcs, aliases))
			}
		}
	default:
		panic("TODO other")
	}
	panic(E)
}
