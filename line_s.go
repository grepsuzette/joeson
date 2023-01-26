package joeson

import (
	"strings"
)

// SLine wraps a string. That string is intended to be parsed
type SLine struct {
	Str string
}

func NewSLine(s string) SLine     { return SLine{s} }
func (sl SLine) Name() string     { panic("uncallable") }
func (sl SLine) Content() Line    { panic("uncallable") }
func (sl SLine) LineType() string { return "s" }
func (sl SLine) StringIndent(nIndent int) string {
	return strings.Replace(sl.Str, "\n", "\\n", -1)
}
