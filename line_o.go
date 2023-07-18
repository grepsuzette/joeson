package joeson

import (
	"regexp"
	"strconv"

	"github.com/grepsuzette/joeson/helpers"
)

type OLine struct {
	name    string // "" unless provided by Named()
	content Line
	attrs   ParseOptions
}

type oLineByIndexOrName struct {
	name  string
	index helpers.NilableInt
}

/*
- O(Named("EXPR", Rules(....)))  // First argument is string (a rule name) and goes to `name`, second is []Line (subrules)
- O("CHOICE _")  // The argument here is considered a rule (sLine) because there is no rules() array. `name` will be ""
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Ast) Ast { return new Choice it}) // same as above, with a cb
- ..... func(it Ast, ctx *ParseContext) Ast { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false } // callbacks long form
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
   // A handcompiled rule with which the joeson grammar is initially defined (see ast/handcompiled.go)
*/

// O() is a helper to declare non-terminal lines of rules (aka OLine).
// It is better to refer to the readme for this function.
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
	if ol.attrs.CbBuilder != nil {
		s += " " + yellow("𝘧")
	}
	return s
}

// You may provide a `grammar` to attempt to parse the rules with, or leave it nil
// which will use the joeson_handcompiled grammar.
func (ol OLine) toRule(rank_ *rank, parentRule Parser, by oLineByIndexOrName, opts TraceOptions, lazyGrammar *helpers.Lazy[*Grammar]) Parser {
	// figure out the name for this rule
	var name string
	var content Line = ol.content
	if ol.name != "" {
		// a named rule, easy
		name = ol.name
	} else if by.name != "" {
		name = by.name
	} else if by.index.IsSet && parentRule != nil {
		name = parentRule.GetRuleName() + "[" + strconv.Itoa(by.index.Int) + "]"
	} else {
		panic("assert")
	}
	rule := getRule(rank_, name, content, parentRule, ol.attrs, opts, lazyGrammar)
	rule.gnode().parent = parentRule
	return rule
}
