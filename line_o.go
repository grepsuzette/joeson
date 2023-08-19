package joeson

import (
	"regexp"
	"strconv"

	"github.com/grepsuzette/joeson/helpers"
)

// "OR" rule. Inside a rank, "OR" rules (AKA OLine) are parsed one after the
// other until one returns something other than nil. Some of them are named,
// but they usually aren't, as it's more the point of an ILine to be
// referenced.
type OLine struct {
	name    string // "" unless provided by Named()
	content Line
	*ParseOptions
}

// In a rank, with a parent "foo", unnamed rules get automatically called
// "foo[0]", "foo[1]"...
// This is the struct transmitting the naming decision to OLine.toRule
type oLineNaming struct {
	name  string
	index int // -1 if unset
}

/*
- O(Named("EXPR", Rules(....)))  // First argument is string (a rule name) and goes to `name`, second is []Line (subrules)
- O("CHOICE _")  // The argument here is considered a rule (sLine) because there is no rules() array. `name` will be ""
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Ast) Ast { return new Choice it}) // same as above, with a cb
- ..... func(it Ast, ctx *ParseContext) Ast { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false } // callbacks long form
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
   // A handcompiled rule with which the joeson grammar is initially defined (see ast/handcompiled.go)
*/

// "OR" rule. Inside a rank, "OR" rules (AKA OLine) are parsed one after the
// other until one returns something other than nil. Some of them are named,
// but they usually aren't, as it's more the point of an ILine to be
// referenced.
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
