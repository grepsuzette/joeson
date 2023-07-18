package joeson

// ParseError implements Ast and indicates a fatal parse error.
//
// It doesn't panic. After `ast := myGrammar.ParseString("FOO")`,
// the correct way to check for error is with `IsParseError(ast)`.
//
// When a parser returns nil, it indicates to backtrack and parse
// in another way. When a ParseError Ast is returned instead,
// the parsing will fail immediately. See examples.
type ParseError struct {
	ctx         *ParseContext
	ErrorString string
}

func (pe ParseError) String() string {
	return "ERROR " + pe.ErrorString + " " + pe.ctx.String()
}

func NewParseError(ctx *ParseContext, s string) ParseError {
	return ParseError{ctx, s}
}

func IsParseError(ast Ast) bool {
	switch ast.(type) {
	case ParseError:
		return true
	default:
	}
	return false
}
