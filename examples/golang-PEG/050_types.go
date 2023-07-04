package main

// https://go.dev/ref/spec#Types

// also depends on 060_grammar_expressions.go

// Type      = TypeName [ TypeArgs ] | TypeLit | "(" Type ")" .
// TypeName  = identifier | QualifiedIdent .
// TypeArgs  = "[" TypeList [ "," ] "]" .
// TypeList  = Type { "," Type } .
// TypeLit   = ArrayType | StructType | PointerType | FunctionType | InterfaceType |
//             SliceType | MapType | ChannelType .

var partial_rules_types = rules(
	o(named("Type", rules(
		o("TypeName TypeArgs? | TypeLit | '(' Type ')'"),
		o(named("TypeLit", rules(
			// "The length is part of the array's type; it must evaluate to
			// a non-negative constant representable by a value of type int.
			// The length of array a can be discovered using the built-in
			// function len. The elements can be addressed by integer indices
			// 0 through len(a)-1. Array types are always one-dimensional but
			// may be composed to form multi-dimensional types."
			o(named("ArrayType", "'[' length:Expression ']' elementType:Type"), x("ArrayType")),
		// o("StructType"),
		// o("PointerType"),
		// o("FunctionType"),
		// o("InterfaceType"),
		// o("SliceType"),
		// o("MapType"),
		// o("ChannelType"),
		))),
		i(named("TypeName", "QualifiedIdent | identifier")),
		i(named("TypeArgs", "'[' TypeList ','? ']'")),
		i(named("TypeList", "Type*','")),
	))),
	o(named("tokens", rules_tokens)), // extending token rules
)
