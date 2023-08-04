package joeson

// ParseContext can react to ParseOptions.
// ParseOptions can originate from:
//   - grammar.ParseString("<text to parse>", <optionalParseOptions>))
//   - from a rule:
//     i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true}, func(it Ast, ctx *ParseContext) Ast { return it })),
//     in that case they will be stored in gnodeimpl.ParseOptions
type ParseOptions struct {
	SkipLog   bool
	SkipCache bool
	Debug     bool
	CbBuilder ParseCallback
}

// Rules often have callbacks. ParseCallback are never directly used,
// instead they are built by O() and I(). Even though ParseCallback has
// 3 arguments, usually one 1 is passed when writing a callback for O() and I()
// as in `func(it j.Ast) j.Ast { return it }`.
//
// Let's consider a few examples below:
//
//   - o(named("foo", "INT"))
//     Declare that rule "foo" must parse rule named "INT".
//     Since this rule has no callback, it will implicitely return the parse
//     result of the parser for rule named "INT" (probably an int).
//     This is equivalent to `func(it j.Ast) j.Ast { return it }`
//
//   - o(named("bar", "m:INT op:('+'|'-') n:INT"))
//     Similar, this time it has 3 labels. A NativeMap with keys m, op
//     and n will be returned here (by parser_sequence.go, to be precise).
//
//   - o(named("bar", "m:INT op:('+'|'-') n:INT"), func(it j.Ast) j.Ast {
//     m := it.(*NativeMap);
//     return newBar(
//     m.GetOrPanic("m"),
//     m.GetOrPanic("op"),
//     m.GetOrPanic("n")
//     })
//
//     A callback is defined.
//     It will build a new object using the 3 fields m, op, and n.
//     This object must implement j.Ast (or j.Parser).
//
// The second arg `*ParseContext` is used usually to build ParseError:
//
//	func(it j.Ast, ctx *ParseContext) j.Ast { return ctx.Error("oops") }
//
// The 3rd arg, `Ast` is the caller Ast and represents the bounded `this` in
// javascript. It is not used, in most cases.
type ParseCallback func(Ast, *ParseContext, Ast) Ast
