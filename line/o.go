package line

import (
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"strconv"
)

type OLine struct {
	name    string // "" unless provided by Named()
	content Line
	attrs   core.ParseOptions
}

type OLineByIndexOrName struct {
	name  string
	index helpers.NilableInt
}

/*
O() is a variadic function, for example:
- O(Named("EXPR", Rules(....)))  // First argument is string (a rule name) and goes to `name`, second is []Line (subrules)
- O("CHOICE _")  // "CHOICE _" here is considered a rule desc because there is no rules() array. `name` will be ""
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Astnode) Astnode { return new Choice it}) // same as above, with a cb
- ..... func(it Astnode, ctx *ParseContext) Astnode { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false } // callbacks long form
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
   // A handcompiled rule with which the joeson grammar is initially defined (see ast/handcompiled.go)
*/
func O(a ...any) OLine {
	name, content, attrs := lineInit(a)
	return OLine{name, content, attrs}
}

func (ol OLine) LineType() string { return "o" }
func (ol OLine) Content() Line    { return ol.content }
func (ol OLine) StringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += ol.LineType()
	s += " "
	s += ol.content.StringIndent(nIndent)
	if ol.attrs.CbBuilder != nil {
		s += Green(", ") + Yellow("𝘧")
	}
	return s
}

func (ol OLine) ToRule(grammar *ast.Grammar, parentRule core.Ast, by OLineByIndexOrName) core.Ast {
	// figure out the name for this rule
	var name string
	var content Line = ol.content
	if ol.name != "" { // A named rule
		name = ol.name
	} else if by.name != "" {
		name = by.name
	} else if by.index.IsSet && parentRule != nil {
		name = parentRule.GetGNode().Name + "[" + strconv.Itoa(by.index.Int) + "]"
	} else {
		panic("assert")
	}
	rule := getRule(grammar, name, content, parentRule, ol.attrs)
	rule.GetGNode().Parent = parentRule
	rule.GetGNode().Index = by.index.Int // 0 by default is fine
	return rule
}
