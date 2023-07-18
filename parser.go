package joeson

// Parser objects are normally built by the joeson grammar
// and are able in turn to parse a ParseContext, producing Ast nodes.
//
// = Errors =
// Parse() should return nil when the current parser failed to recognize
// anything,
//
//	but it should return an AstError when current context IS the correct one
//	handled by this parser but the context somehow is erroneous. For example,
//	parsing "0o9" for an octal parser may return NewAstError(parser, "0o9" is
//	invalid octal").
type Parser interface {
	Ast
	gnode
	Parse(ctx *ParseContext) Ast
	ForEachChild(func(Parser) Parser) Parser // depth-first walk mapper
	prepare()
	handlesChildLabel() bool
}

var (
	_ Parser = &choice{}
	_ Parser = &existential{}
	_ Parser = &lookahead{}
	_ Parser = &not{}
	_ Parser = &pattern{}
	_ Parser = &rank{}
	_ Parser = &regex{}
	_ Parser = &sequence{}
	_ Parser = &str{}
	_ Parser = &Grammar{}
	_ Parser = &cLine{}
	_ Parser = &NativeUndefined{}
)

func IsRule(parser Parser) bool {
	return parser.gnode().rule == parser
}

// Return a prefix consisting of a name or a label when appropriate.
func prefix(parser Parser) string {
	if IsRule(parser) {
		return red(parser.Name() + ": ")
	} else if parser.Label() != "" {
		return cyan(parser.Label() + ":")
	} else {
		return ""
	}
}
