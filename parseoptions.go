package joeson

// ParseOption can apply to:
//   - rules: `i(Named("EXAMPLE", "/regex/", Debug{true})),`
//   - individual cases: `grammar.ParseString("<text to parse>", j.Debug{true}))`
//
// You almost always want the later, e.g.: `grammar.ParseString(text, j.Debug{true})`.
// For rules, it will be effective while compiling the grammar itself (not during parsing).
type ParseOption interface {
	apply(*parseOptions) *parseOptions
}

var (
	_ ParseOption = SkipLog{true}
	_ ParseOption = SkipCache{true}
	_ ParseOption = Debug{true}
)

type (
	SkipLog   struct{ Bool bool }
	SkipCache struct{ Bool bool }
	Debug     struct{ Bool bool }
)

func (o SkipLog) apply(opts *parseOptions) *parseOptions   { opts.skipLog = o.Bool; return opts }
func (o SkipCache) apply(opts *parseOptions) *parseOptions { opts.skipCache = o.Bool; return opts }
func (o Debug) apply(opts *parseOptions) *parseOptions     { opts.debug = o.Bool; return opts }

// legacy ParseOption store composed by ParseContext
type parseOptions struct {
	skipLog   bool
	skipCache bool
	debug     bool

	// parse function. Or rule callback.
	// This is the full signature. In practice shorter signature
	// are used, and lineInit() builds this full signature from
	// either `func(Ast) Ast` or `func(Ast, *ParseContext) Ast`.
	//
	// Arguments:
	//
	// - `it Ast`: input Ast to map to the return.
	// - `ctx *ParseContext`: Usually to build ParseError:
	//   `panic(ctx.Error("illegal phone number"))`
	// - `caller Ast`: bounded `this` in JS. Almost never used.
	cb func(it Ast, ctx *ParseContext, caller Ast) Ast
}

func newParseOptions() *parseOptions {
	return &parseOptions{
		skipLog:   false,
		skipCache: false,
		debug:     false,
		cb:        nil,
	}
}
