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
	var b strings.Builder
	b.WriteString(BoldBlue("[\n"))
	for _, line := range al.Array {
		b.WriteString(line.stringIndent(nIndent + 1))
		b.WriteString("\n")
	}
	b.WriteString(helpers.Indent(nIndent))
	b.WriteString(BoldBlue("]"))
	return b.String()
}
