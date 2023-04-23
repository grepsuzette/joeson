package main

import j "github.com/grepsuzette/joeson"

// x("foo") has only one purpose.
// It creates a callback to prepare a dummy Ast.
// This dummy Ast boasts a single method ContentString()
//   which returns "foo:" followed by its ast.ContentString().
//
// This helps to quickly write a grammar.
//
// For example:
// var rules_tokens = rules(
// 	o(named("token", "( keyword | identifier | operator | punctuation | literal )"), x("token")),
// 	i(named("keyword", "( 'break' | 'default' | 'func' | 'interface' | 'select' | 'case' | 'defer' | 'go' | 'map' | 'struct' | 'chan' | 'else' | 'goto' | 'package' | 'switch' | 'const' | 'fallthrough' | 'if' | 'range' | 'type' | 'continue' | 'for' | 'import' | 'return' | 'var' )"), x("keyword")),
// 	i(named("identifier", "[a-zA-Z_][a-zA-Z0-9_]*"), x("identifier")), // letter { letter | unicode_digit } .   We rewrite it so to accelerate parsing
// 	i(named("operator", "( '+' | '&' | '+=' | '&=' | '&&' | '==' | '!=' | '(' | ')' | '-' | '|' | '-=' | '|=' | '||' | '<' | '<=' | '[' |  ']' | '*' | '^' | '*=' | '^=' | '<-' | '>' | '>=' | '{' | '}' | '/' | '<<' | '/=' | '<<=' | '++' | '=' | ':=' | '%' | '>>' | '%=' | '>>=' | '--' | '!' | '...' | '&^' | '&^=' | '~' )"), x("operator")),
// ...
// )

type dumb struct {
	typename string
	ast      j.Ast
}

func (dumb dumb) ContentString() string {
	return "<" + dumb.typename + ":" + dumb.ast.ContentString() + ">"
}

func x(typename string) func(j.Ast) j.Ast {
	return func(ast j.Ast) j.Ast {
		return dumb{typename, ast}
	}
}
