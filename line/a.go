package line

import (
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/helpers"
	"grepsuzette/joeson/lambda"
	"strings"
)

type ALine struct {
	Array []Line
}

func NewALine(a []Line) ALine     { return ALine{a} }
func (al ALine) Content() Line    { panic("uncallable") }
func (al ALine) LineType() string { return "a" }
func (al ALine) StringIndent(nIndent int) string {
	return BoldBlue("[\n") + strings.Join(
		lambda.Map(al.Array, func(line Line) string {
			return line.StringIndent(nIndent + 1)
		}),
		"\n",
	) + "\n" + helpers.Indent(nIndent) + BoldBlue("]")
}
