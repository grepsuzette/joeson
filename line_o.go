package joeson

import (
	"regexp"
	"strconv"

	"github.com/grepsuzette/joeson/helpers"
)

type OLine struct {
	name    string // "" unless provided by Named()
	content Line
	*ParseOptions
}

func O(a ...any) OLine {
	name, content, attrs := lineInit(a)
	return OLine{name, content, attrs}
}

func (ol OLine) lineType() string { return "o" }
func (ol OLine) stringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += ol.lineType()
	s += " " + ol.name + " "
	re := regexp.MustCompile("^ *o *")
	s += re.ReplaceAllString(ol.content.stringIndent(nIndent), "o ")
	if ol.ParseOptions.CbBuilder != nil {
		s += " " + Yellow("𝘧")
	}
	return s
}

// Convert a O line to a rule Parser.
// lazyGrammar: normally nil. A lazy grammar can be provided here. nil uses joeson_handcompiled grammar.
func (ol OLine) toRule(
	rank_ *rank,
	parentRule Parser,
	by oLineNaming,
	opts *TraceOptions,
	lazyGrammar *helpers.Lazy[*Grammar],
) Parser {
	// figure out the name for this rule
	var name string
	var content Line = ol.content
	if ol.name != "" {
		name = ol.name // named rule
	} else if by.name != "" {
		name = by.name
	} else if by.index > -1 && parentRule != nil { // unamed OLine by index, will produce names like foo[0], foo[1], foo[2]
		name = parentRule.GetRuleName() + "[" + strconv.Itoa(by.index) + "]"
	} else {
		panic("assert")
	}
	rule := getRule(rank_, name, content, parentRule, ol.ParseOptions, opts, lazyGrammar)
	rule.gnode().parent = parentRule
	return rule
}

// specifies how to reference a rule
// (either by name, or by index)
type oLineNaming struct {
	name  string
	index int // -1 if unset
}
