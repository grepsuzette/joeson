package joeson

// Parser objects are able to parse a ParseContext producing Ast.
//
// Choice, Existential, Lookahead, Not, Pattern, Ref,
// Regex, Sequence, and joeson also have Rank and Str are built-in Parsers
// in joeson, though arguably and in a way any rule or any callback rule is
// a kind of parser.
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
	GNode
	Parse(ctx *ParseContext) Ast
	Prepare()
	HandlesChildLabel() bool
	ForEachChild(f func(Parser) Parser) Parser // depth-first walk enabler
}

func IsRule(parser Parser) bool {
	return parser.GetGNode().rule == parser
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
