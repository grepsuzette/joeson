package line

import (
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/helpers"
	"strings"
)

type ALine struct {
	Array []Line
}

func NewALine(a []Line) ALine     { return ALine{a} }
func (al ALine) Name() string     { panic("uncallable") }
func (al ALine) Content() Line    { panic("uncallable") }
func (al ALine) LineType() string { return "a" }
func (al ALine) StringIndent(nIndent int) string {
	return BoldBlue("[\n") + strings.Join(
		helpers.AMap(al.Array, func(line Line) string {
			return line.StringIndent(nIndent + 1)
		}),
		"\n",
	) + "\n" + helpers.Indent(nIndent) + BoldBlue("]")
}
