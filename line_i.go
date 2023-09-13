package joeson

import (
	"regexp"

	"github.com/grepsuzette/joeson/helpers"
)

type ILine struct {
	name    string // name must be set for ILine
	content Line
	*ParseOptions
}

func I(a ...any) ILine {
	name, content, attrs := lineInit(a)
	if name == "" {
		panic("ILine must always be named")
	}
	return ILine{name, content, attrs}
}

func (il ILine) lineType() string { return "i" }
func (il ILine) stringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += il.lineType()
	s += " " + il.name + " "
	re := regexp.MustCompile("^ *o *")
	s += re.ReplaceAllString(il.content.stringIndent(nIndent), "o ")
	if il.ParseOptions.CbBuilder != nil {
		s += " " + Yellow("𝘧")
	}
	return s
}

// Convert a ILine to a rule Parser.
func (il ILine) toRule(
	rank_ *rank,
	parentRule Parser,
	opts *TraceOptions,
	lazyGrammar *helpers.Lazy[*Grammar],
) (name string, rule Parser) {
	return il.name, getRule(rank_, il.name, il.content, parentRule, il.ParseOptions, opts, lazyGrammar)
}
