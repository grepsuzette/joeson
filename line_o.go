package joeson

import (
	"grepsuzette/joeson/helpers"
	"strconv"
)

type OLine struct {
	name    string // "" unless provided by Named()
	content Line
	attrs   ParseOptions
}

type OLineByIndexOrName struct {
	name  string
	index helpers.NilableInt
}

// TODO fix that doc:

/*
O() is a variadic function, for example:
- O(Named("EXPR", Rules(....)))  // First argument is string (a rule name) and goes to `name`, second is []Line (subrules)
- O("CHOICE _")  // "CHOICE _" here is considered a rule desc because there is no rules() array. `name` will be ""
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Ast) Ast { return new Choice it}) // same as above, with a cb
- ..... func(it Ast, ctx *ParseContext) Ast { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false } // callbacks long form
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
   // A handcompiled rule with which the joeson grammar is initially defined (see ast/handcompiled.go)
*/
func O(a ...any) OLine {
	name, content, attrs := lineInit(a)
	return OLine{name, content, attrs}
}

func (ol OLine) LineType() string { return "o" }
func (ol OLine) Name() string     { return ol.name }
func (ol OLine) Content() Line    { return ol.content }
func (ol OLine) StringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += ol.LineType()
	s += " "
	s += ol.content.StringIndent(nIndent)
	if ol.attrs.CbBuilder != nil {
		s += green(", ") + yellow("𝘧")
	}
	return s
}

// You may provide a `grammar` to attempt to parse the rules with, or leave it nil
// which will use the joeson_handcompiled grammar.
func (ol OLine) toRule(rank *Rank, parentRule Ast, by OLineByIndexOrName, opts TraceOptions, lazyGrammar *helpers.Lazy[*Grammar]) Ast {
	// figure out the name for this rule
	var name string
	var content Line = ol.content
	if ol.name != "" {
		// a named rule, easy
		name = ol.name
	} else if by.name != "" {
		name = by.name
	} else if by.index.IsSet && parentRule != nil {
		name = parentRule.GetGNode().Name + "[" + strconv.Itoa(by.index.Int) + "]"
	} else {
		panic("assert")
	}
	rule := getRule(rank, name, content, parentRule, ol.attrs, opts, lazyGrammar)
	rule.GetGNode().Parent = parentRule
	rule.GetGNode().Index = by.index.Int // 0 by default is fine
	return rule
}
