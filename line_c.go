package joeson

// cLine wraps an Ast (more precisely a Parser) and therefore represents
// a rule which has already been compiled.
type cLine struct {
	Parser
}

func newCLine(x Parser) cLine     { return cLine{x} }
func (cl cLine) lineType() string { return "c" }
func (cl cLine) stringIndent(nIndent int) string {
	return cl.Parser.ContentString()
}
