package line

import (
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"grepsuzette/joeson/lambda"
	"reflect"
	"strconv"
	"strings"
)

type OLine struct {
	name    string // can be empty, or provided by Named()
	content Line
	attrs   core.ParseOptions
}

type OLineByIndexOrByName struct {
	name  string
	index helpers.NullInt
}

// func OEmpty() OLine { return OLine{[]any{}} }

/*
O() is a variadic function which allows a variety of declarations, for example:
- O("EXPR", Rules(....))       // "EXPR" is a rule name
- O("CHOICE _")                // "CHOICE _" is a rule desc because there is no rules() array
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Astnode) Astnode { return new Choice it})
- O("<name>", func(it Astnode, ctx *ParseContext) Astnode { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false }
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
The last one is a handcompiled rule with which the joeson grammar is initially defined as in ast/handcompiled
*/
func O(a ...any) OLine {
	name, content, attrs := lineInit(a)
	return OLine{name, content, attrs}
}

func (ol OLine) LineType() string { return "o" }
func (ol OLine) Content() Line    { return ol.content }
func (ol OLine) String() string   { return ol.StringIndent(0) }
func (ol OLine) StringIndent(nIndent int) string {
	s := helpers.Indent(nIndent) + ol.LineType() + " "
	switch v := ol.content.(type) {
	case SLine:
		s += v.Str
	case CLine:
		s += v.Astnode.ContentString()
	case ALine:
		s += BoldBlue("[\n") + strings.Join(
			lambda.Map(v.Array, func(line Line) string { return line.StringIndent(nIndent + 1) }),
			"\n",
		) + "\n" + helpers.Indent(nIndent) + BoldBlue("]")
	default:
		s += "?" + reflect.TypeOf(v).String() + "?"
	}
	if ol.attrs.CbBuilder != nil {
		s += Green(", ") + Yellow("𝘧")
	}
	return s
}

// note TODO think parentRule could almost simply be GNode. but anyway
func (ol OLine) ToRule(grammar *ast.Grammar, parentRule core.Astnode, by OLineByIndexOrByName) core.Astnode {
	// figure out the name for this rule
	if ol.name != "" {
		by.name = ol.name
	} else if by.name == "" && by.index.IsSet && parentRule != nil {
		by.name = parentRule.GetGNode().Name + "[" + strconv.Itoa(by.index.Int) + "]"
	} else if by.name == "" {
		panic("Name undefined for 'o' line")
	}
	rule := getRule(grammar, by.name, ol.content, parentRule, ol.attrs)
	rule.GetGNode().Parent = parentRule
	// TODO is the following commented line really useful? I am not sure yet
	// rule.GetGNode().Index = by.index
	return rule
}
