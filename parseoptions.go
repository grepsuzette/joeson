package joeson

// ParseContext can react to ParseOptions.
// Those options should be considered legacy,
// they come from the original implementation.
// but were never used during the developement of this version.
//
// ParseOptions can originate from:
//   - grammar.ParseString("<text to parse>", <optionalParseOptions>))
//   - a rule:
//     i(Named("EXAMPLE", "/regex/", ParseOptions{SkipLog: false, SkipCache: true}, func(it Ast, ctx *ParseContext) Ast { return it })),
//     in that case they will be stored in gnodeimpl.ParseOptions
type ParseOptions struct {
	SkipLog   bool
	SkipCache bool
	Debug     bool
	CbBuilder parseCallback
}

func newParseOptions() *ParseOptions {
	return &ParseOptions{
		SkipLog:   false,
		SkipCache: false,
		Debug:     false,
		CbBuilder: nil,
	}
}

// The second arg `*ParseContext` allows usually used to build ParseError:
//
//	func(it j.Ast, ctx *ParseContext) j.Ast { return ctx.Error("oops") }
//
// The 3rd arg, `Ast` is the caller Ast and represents the bounded `this` in
// javascript. It is not used, in most cases.
type parseCallback func(Ast, *ParseContext, Ast) Ast
