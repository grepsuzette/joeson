package joeson

// cLine wraps an Ast and therefore represents
// a rule which was already compiled.
type cLine struct {
	Ast
}

func newCLine(x Ast) cLine        { return cLine{x} }
func (cl cLine) lineType() string { return "c" }
func (cl cLine) stringIndent(nIndent int) string {
	return cl.Ast.ContentString()
}
