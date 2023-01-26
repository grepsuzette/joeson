package joeson

// Each ParseContext have ParseOptions.
// They can originate from different places:
// - grammar.ParseString("<text to parse>", <optionalParseOptions>))
// - or from a rule:
//    i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true}, func(it Ast, ctx *ParseContext) Ast { return it })),
//    in which case they will be stored in GNode.ParseOptions
type ParseOptions struct {
	SkipLog   bool
	SkipCache bool
	Debug     bool
	/*
	 `cbBuilder` represents optional callbacks declared within inlined rules.
	 E.g. the func in `o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?",
	 		   func(result Ast) Ast { return ast.NewPattern(result) }),`

	 Since this example have labels, `result` will be of type NativeMap (which
	 implements Ast) with the 3 keys "value", "join" and "@". Otherwise
	 it will be a NativeArray.

	 Second arg `...*ParseContext` is rarely passed in practice,
	 see a rare use in joescript.coffee:660.

	 Third arg `Ast` is the caller Ast (see joeson.js:455
	 or joeson.coffee:278) and represents the bounded `this` in javascript.
	*/
	CbBuilder func(Ast, *ParseContext, Ast) Ast
}
