package joeson

import (
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// ALine embeds a []Line.
type ALine struct {
	Array []Line
}

func NewALine(a []Line) ALine     { return ALine{a} }
func (al ALine) lineType() string { return "a" }
func (al ALine) stringIndent(nIndent int) string {
	return boldBlue("[\n") + strings.Join(
		helpers.AMap(al.Array, func(line Line) string {
			return line.stringIndent(nIndent + 1)
		}),
		"\n",
	) + "\n" + helpers.Indent(nIndent) + boldBlue("]")
}
