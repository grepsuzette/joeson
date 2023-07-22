package main

import "github.com/grepsuzette/joeson"

// function x() helps to quickly write a grammar.
// Calling x("foo") returns a callback `func(τ Ast) Ast`.
// Calling cb.String() gives "<foo:" + τ.ContentString() + ">"
//
// For example:
//
// var rules_tokens = rules(
//
//	o(named("token", "( keyword | identifier | operator | punctuation | literal )"), x("token")),
//	i(named("keyword", "( 'break' | 'default' | 'func' | 'interface' | 'select' | 'case' | 'defer' | 'go' | 'map' | 'struct' | 'chan' | 'else' | 'goto' | 'package' | 'switch' | 'const' | 'fallthrough' | 'if' | 'range' | 'type' | 'continue' | 'for' | 'import' | 'return' | 'var' )"), x("keyword")),
//	i(named("identifier", "[a-zA-Z_][a-zA-Z0-9_]*"), x("identifier")), // letter { letter | unicode_digit } .   We rewrite it so to accelerate parsing
//	i(named("operator", "( '+' | '&' | '+=' | '&=' | '&&' | '==' | '!=' | '(' | ')' | '-' | '|' | '-=' | '|=' | '||' | '<' | '<=' | '[' |  ']' | '*' | '^' | '*=' | '^=' | '<-' | '>' | '>=' | '{' | '}' | '/' | '<<' | '/=' | '<<=' | '++' | '=' | ':=' | '%' | '>>' | '%=' | '>>=' | '--' | '!' | '...' | '&^' | '&^=' | '~' )"), x("operator")),
//
// ...
// )
//
// Here, whichever of keyword, identifier etc gets built,
// its String() will be like "<token:keyword>", "<token:identifier>" etc.
func x(typename string) func(joeson.Ast) joeson.Ast {
	return func(ast joeson.Ast) joeson.Ast {
		return dumb{typename, ast, ast.GetOrigin()}
	}
}

// ParseError is an ast denoting parse errors
type ParseError struct{ string }

func NewParseError(s string) ParseError { return ParseError{s} }
func (e ParseError) String() string     { return e.string }

// type dumb is used by x(). As the name hints, it's nothing too exciting
type dumb struct {
	typename string
	ast      joeson.Ast
	origin   joeson.Origin
}

func (dumb dumb) assertNode()                                     {}
func (dumb dumb) String() string                                  { return "<" + dumb.typename + ":" + dumb.ast.String() + ">" }
func (dumb dumb) HasAttribute(key interface{}) bool               { return false }
func (dumb dumb) GetAttribute(key interface{}) interface{}        { panic("not implemented") }
func (dumb dumb) SetAttribute(key interface{}, value interface{}) { panic("not implemented") }
func (dumb dumb) GetOrigin() joeson.Origin                        { return dumb.origin }
func (dumb dumb) SetOrigin(o joeson.Origin)                       { dumb.origin = o }
func (dumb dumb) GetLine() int                                    { return dumb.origin.Line }
func (dumb dumb) SetLine(n int) {
	dumb.origin.Line = n
	dumb.origin.Start = 0
	dumb.origin.End = 0
}
