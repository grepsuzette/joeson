package main

import (
	"reflect"

	j "github.com/grepsuzette/joeson"
)

// -- The parsing grammar

func i(a ...any) j.ILine                   { return j.I(a...) }
func o(a ...any) j.OLine                   { return j.O(a...) }
func named(name string, v any) j.NamedRule { return j.Named(name, v) }
func rules(a ...j.Line) []j.Line           { return a }

var grammarRules = rules(
	o(`_ expr _`),
	i(named("expr", rules(
		o(named("list", `'(' _  expr*__  _ ')'`), exprFromList),
		o(named("string", `'"' [^"]* '"'`), exprFromString),
		o(named("number", "/-?[0-9]+/"), exprFromNumber),
		o(named("operator", `word | '+' | '-' | '*' | '/' | '%' | '>=' | '<=' | '!=' | '=='| '<' | '>' |  '=' `), exprFromOperator),
		i(named("word", `/[a-zA-Z\._][a-zA-Z\._0-9?]*/`)),
	))),
	i(named("_", `_SPACECHAR*`)),
	i(named("__", `_SPACECHAR+`)),
	i(named("_SPACECHAR", "' ' | '\t' | '\n'")),
	i(named(".", `'\s' | '\S'`)),
	i(named("ESC1", `'\\' .`)),
)

func exprFromList(it j.Ast) j.Ast {
	return Expr{j.NewAttr(), kindList, "", 0, parseList(it).(List), ""}
}

func exprFromString(it j.Ast) j.Ast {
	return Expr{j.NewAttr(), kindString, it.(*j.NativeArray).Concat().String(), 0, nilList(), ""}
}

func exprFromNumber(it j.Ast) j.Ast {
	return Expr{j.NewAttr(), kindNumber, "", float64(j.NewNativeIntFrom(it).Int()), nilList(), ""}
}

func exprFromOperator(it j.Ast) j.Ast {
	return Expr{j.NewAttr(), kindOperator, "", 0, nilList(), it.(j.NativeString).String()}
}

func parseList(it j.Ast) j.Ast {
	exprs := []Expr{}
	for _, v := range *it.(*j.NativeArray) {
		switch v := v.(type) {
		case Expr:
			exprs = append(exprs, v)
		default:
			panic("expected Expr, got " + reflect.TypeOf(v).String())
		}
	}
	return list(exprs...)
}
