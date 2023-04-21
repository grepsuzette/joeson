package main

// https://go.dev/ref/spec#Expressions

// Operands
// Operand     = Literal | OperandName [ TypeArgs ] | "(" Expression ")" .
// Literal     = BasicLit | CompositeLit | FunctionLit .
// BasicLit    = int_lit | float_lit | imaginary_lit | rune_lit | string_lit .
// OperandName = identifier | QualifiedIdent .

var partial_rules_expressions = rules( // also depends on 050_grammar_types.go
	o(named("Expression", "Operand")),
	o(named("Operand", rules(
		o("'(' Expression ')' | OperandName TypeArgs? | Literal"),
		o(named("Literal", "BasicLit | CompositeLit | FunctionLit")),
		// TODO add float_lit and imaginary_lit
		o(named("BasicLit", "int_lit | rune_lit | string_lit")),
		o(named("OperandName", "QualifiedIdent | identifier")),
		i(named("QualifiedIdent", "PackageName '.' identifier"), x("QualifiedIdent")), // https://go.dev/ref/spec#QualifiedIdent
		i(named("PackageName", "identifier")),                                         // https://go.dev/ref/spec#PackageName

		o(named("Block", "'{' Statement*';' '}'")),
	))),
	o(named("token", rules_tokens)), // import previous rules
)
