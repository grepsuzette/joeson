package joeson

// Parsers are built ("compiled") by a joeson grammar.
// Each rule's Parser is a complex recursive dynamic arrangement.
// A Parser can parse a ParseContext, producing Ast nodes.
//
// Note although Parser is exported, there is no established way
// to really work with it yet outside from this package.
type Parser interface {
	Ast
	parse(ctx *ParseContext) Ast
	getRule() *rule
	prepare()
	handlesChildLabel() bool                 // prevents collecting labelled children as part of a NativeMap in packrat.go prepareResult()
	forEachChild(func(Parser) Parser) Parser // depth-first walk mapper
}

var (
	_ Parser = &Grammar{}
	_ Parser = &choice{}
	_ Parser = &existential{}
	_ Parser = &lookahead{}
	_ Parser = &not{}
	_ Parser = &pattern{}
	_ Parser = &rank{}
	_ Parser = &regex{}
	_ Parser = &sequence{}
	_ Parser = &str{}
	_ Parser = &cLine{}
	_ Parser = &NativeUndefined{}
)

// Return a prefix consisting of a name or a label when appropriate.
func prefix(parser Parser) string {
	if IsRule(parser) {
		return Red(parser.getRule().name + ": ")
	} else if parser.getRule().label != "" {
		return Cyan(parser.getRule().label + ":")
	} else {
		return ""
	}
}

func IsRule(parser Parser) bool {
	return parser.getRule().parser == parser
}
