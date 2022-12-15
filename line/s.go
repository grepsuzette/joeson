package line

import (
	"strings"
)

// SLine wraps a string, that is intended to be parsed
type SLine struct {
	Str string
}

func NewSLine(s string) SLine     { return SLine{s} }
func (sl SLine) Content() Line    { panic("uncallable") }
func (sl SLine) LineType() string { return "s" }
func (sl SLine) StringIndent(nIndent int) string {
	return strings.Replace(sl.Str, "\n", "\\n", -1)
}
