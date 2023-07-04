package main

// https://go.dev/ref/spec#Constants
// https://go.dev/ref/spec#Constant_expressions

// ConstDecl      = "const" ( ConstSpec | "(" { ConstSpec ";" } ")" ) .
// ConstSpec      = IdentifierList [ [ Type ] "=" ExpressionList ] .
// IdentifierList = identifier { "," identifier } .
// ExpressionList = Expression { "," Expression } .

// this one is problematic for now:
//
//	depends on partial_rules_types and rules_tokens
//	also it would be better work on a pre-tokenized the source
//	 to simplify the grammar
//
// so let's see in 100_decl.go
var rules_constants = rules(
// o(named("ConstDecl", rules(
// o("'const' _ ( '(' ConstSpec*';' ')' | ConstSpec )"),
// i(named("ConstSpec", "IdentifierList ( Type? __ '=' __ ExpressionList )?")),
// o(named("IdentifierList", "identifier ( ',' identifier )*")),
// o(named("ExpressionList", "Expression ( ',' Expression )")),
// o(named("__type", partial_rules_types)),
// o(named("__identifier", rules_tokens)),
// ))),
)
