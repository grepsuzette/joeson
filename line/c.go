package line

import . "grepsuzette/joeson/core"

// CLine wraps an Astnode and therefore represents
// a rule which was already compiled in a Line tree.
type CLine struct {
	Astnode
}

func NewCLine(x Astnode) CLine                   { return CLine{x} }
func (cl CLine) Content() Line                   { panic("uncallable") }
func (cl CLine) LineType() string                { return "c" }
func (cl CLine) String() string                  { return "cline" }
func (cl CLine) StringIndent(nIndent int) string { return "cline" }
