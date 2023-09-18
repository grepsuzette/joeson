package main

import (
	j "github.com/grepsuzette/joeson"
)

// -- The parsing grammar

func i(a ...any) j.ILine                                 { return j.I(a...) }
func o(a ...any) j.OLine                                 { return j.O(a...) }
func named(name string, lineStringOrAst any) j.NamedRule { return j.Named(name, lineStringOrAst) }
func rules(a ...j.Line) []j.Line                         { return a }

var grammarRules = rules(
	o(`_ expr:expr _`, parseTopLevelExpr),
	i(named("expr", `l:list | s:string | n:number | operator:operator`), parseExpr),
	i(named("list", `'(' _ (expr*__) _ ')'`), parseList),
	i(named("operator", `word | '+' | '-' | '*' | '/' | '%' | '>=' | '<=' | '!=' | '=='| '<' | '>' |  '=' `), parseOperator),
	i(named("_", "(' ' | '\t' | '\n')*")),
	i(named("__", "(' ' | '\t' | '\n')+")),
	i(named("string", "'\"' s:([^\"]*) '\"'"), parseString),
	i(named("word", "/[a-zA-Z\\._][a-zA-Z\\._0-9?]*/")),
	i(named("number", "/-?[0-9]+/") /* TODO non-int */, func(it j.Ast) j.Ast { return j.NativeIntFrom(it) }),
	i(named(".", "/[\\s\\S]/")),
	i(named("ESC1", "'\\\\' .")),
)

func parseTopLevelExpr(it j.Ast) j.Ast {
	if ast, ok := it.(*j.NativeMap).GetExists("expr"); ok {
		switch v := ast.(type) {
		case Expr:
			return v
		case List:
			return list()
		default:
			return ast
		}
	} else {
		panic("743")
	}
}

func parseExpr(it j.Ast) j.Ast {
	if h, ok := it.(*j.NativeMap); !ok {
		panic(E)
	} else {
		if v, ok := h.GetExists("s"); ok {
			if ns, ok := v.(j.NativeString); ok {
				return Expr{j.NewAttr(), kindString, string(ns), 0, nilList(), ""}
			} else {
				panic("24942")
			}
		} else if v, ok := h.GetExists("n"); ok {
			if n, ok := v.(j.NativeInt); ok {
				return Expr{j.NewAttr(), kindNumber, "", float64(n.Int()), nilList(), ""}
			} else {
				panic("24942")
			}
		} else if v, ok := h.GetExists("l"); ok {
			return Expr{j.NewAttr(), kindList, "", 0, v.(List), ""}
		} else if v, ok := h.GetExists("operator"); ok {
			return Expr{j.NewAttr(), kindOperator, "", 0, nilList(), string(v.(j.NativeString))}
		} else {
			panic(h.String())
		}
	}
}

func parseList(it j.Ast) j.Ast {
	exprs := []Expr{}
	for _, v := range *it.(*j.NativeArray) {
		switch v := v.(type) {
		case Expr:
			exprs = append(exprs, v)
		default:
			panic("f9439393")
		}
	}
	return list(exprs...)
}

func parseOperator(it j.Ast) j.Ast {
	return it.(j.NativeString)
	// if s, ok := it.(j.NativeString).Str; ok {
	// 	return j.NewNativeString(s)
	// } else {
	// 	panic("38274")
	// }
}

func parseString(it j.Ast) j.Ast {
	if s, ok := it.(*j.NativeMap).GetStringExists("s"); ok {
		return j.NewNativeString(s)
	} else {
		panic("38274")
	}
}
