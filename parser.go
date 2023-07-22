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
	/*
		// FIXED Hum...
		// as of July 22, 2023 some changes made
		// this runtime error happen on str type:
		// "comparing uncomparable type joeson.str"
		// This happens when both types are str, and
		// even if they are similar.
		// Think it might be related to our refactoring
		// (we now use Attr, which ~is~was using
		// map[interface{}]interface{}.
		// Turns out it makes all our nodes not comparable,
		// and our algorithm right now needs them to be.
		//
		// As confirmed by this blog
		// https://go.dev/blog/comparable
		//
		// So let's remove our map in attr.go until
		// a good-enough solution presents itself.
		//
		if parser.gnode().rule == nil {
			return false
		}
		t1 := reflect.TypeOf(parser.gnode().rule).String()
		t2 := reflect.TypeOf(parser).String()
		if t1 != t2 {
			return false
		}
		// if t1 == "joeson.str" && t2 == "joeson.str" {
		// 	return ............  // let's not do this
		// }
		return parser.gnode().rule == parser
	*/
	return parser.gnode().rule == parser
}

// Return a prefix consisting of a name or a label when appropriate.
func prefix(parser Parser) string {
	if IsRule(parser) {
		return Red(parser.GetRuleName() + ": ")
	} else if parser.GetRuleLabel() != "" {
		return Cyan(parser.GetRuleLabel() + ":")
	} else {
		return ""
	}
}
