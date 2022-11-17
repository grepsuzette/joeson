package line

type ALine struct {
	Array []Line
}

func NewALine(a []Line) ALine                    { return ALine{a} }
func (al ALine) Content() Line                   { panic("uncallable") }
func (al ALine) LineType() string                { return "a" }
func (al ALine) String() string                  { return "aline" }
func (al ALine) StringIndent(nIndent int) string { return "aline" }
