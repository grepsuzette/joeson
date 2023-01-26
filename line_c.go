package joeson

// CLine wraps an Ast and therefore represents
// a rule which was already compiled in a Line tree.
type CLine struct {
	Ast
}

func NewCLine(x Ast) CLine        { return CLine{x} }
func (cl CLine) Name() string     { panic("uncallable") }
func (cl CLine) Content() Line    { panic("uncallable") }
func (cl CLine) LineType() string { return "c" }
func (cl CLine) StringIndent(nIndent int) string {
	return cl.Ast.ContentString()
}
