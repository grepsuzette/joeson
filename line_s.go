package joeson

import (
	"strings"
)

// sLine wraps a string. That string is intended to be parsed
// and represents a rule that has not been compiled yet.
// It is never entered directly however, it's a transient state.
type sLine struct {
	Str string
}

func newSLine(s string) sLine     { return sLine{s} }
func (sl sLine) lineType() string { return "s" }
func (sl sLine) stringIndent(nIndent int) string {
	return strings.Replace(sl.Str, "\n", "\\n", -1)
}
