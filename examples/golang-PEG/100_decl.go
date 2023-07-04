package main

var rules_decl = rules(
	// o(named("TopLevelDecl", "Declaration | FunctionDecl | MethodDecl")),
	o(named("TopLevelDecl", "Declaration")),
	o(named("__expressions", partial_rules_expressions)), // import previous rules
	o(named("__types", partial_rules_types)),
	o(named("Declaration", rules(
		// o("ConstDecl | TypeDecl | VarDecl"),
		o("ConstDecl"),
		o(named("ConstDecl", "'const' ( '(' ConstSpec*';' ')' | ConstSpec )")),
		i(named("ConstSpec", "IdentifierList ( Type? '=' ExpressionList )?")),
		o(named("IdentifierList", "identifier ( ',' identifier )*")),
		o(named("ExpressionList", "Expression ( ',' Expression )")),
	))),
)
