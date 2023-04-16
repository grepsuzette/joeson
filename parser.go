package joeson

// Parser objects are able to to parse other grammars.
// The joeson grammar produce Ast nodes that satisfy both
// interfaces Ast and Parser.
//
// These nodes are the usual suspects for PEG grammars,
// namely Choice, Existential, Lookahead, Not, Pattern, Ref,
// Regex, Sequence, and joeson also have Rank and Str.
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
