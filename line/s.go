package line

// SLine wraps a string, that is intended to be Parse-d
//  and become a CLine (which wraps Astnode)
type SLine struct {
	Str string
}

func NewSLine(s string) SLine                    { return SLine{s} }
func (sl SLine) Content() Line                   { panic("uncallable") }
func (sl SLine) LineType() string                { return "s" }
func (sl SLine) String() string                  { return "sline:" + sl.Str }
func (sl SLine) StringIndent(nIndent int) string { return "sline" }
