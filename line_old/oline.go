package line

import (
	"fmt"
	"grepsuzette/joeson/ast"
	"grepsuzette/joeson/helpers"
	"grepsuzette/joeson/lambda"
	"reflect"
	"strconv"
	"strings"

	. "grepsuzette/joeson/core"

	. "grepsuzette/joeson/colors"
)

type OLine struct {
	args []any
}

func OEmpty() OLine { return OLine{[]any{}} }

/*
O() is a variadic function which allows a variety of declarations, for example:
- O("EXPR", Rules(....))       // "EXPR" is a rule name
- O("CHOICE _")                // "CHOICE _" is a rule desc because there is no rules() array
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Astnode) Astnode { return new Choice it})
- O("<name>", func(it Astnode, ctx *ParseContext) Astnode { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false }
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
The last one is a handcompiled rule with which the joeson grammar is initially defined as in ast/handcompiled
*/
func O(a ...any) OLine { return OLine{a} }

func (ol OLine) IsO() bool        { return true }
func (ol OLine) LineType() string { return "o" }
func (ol OLine) String() string   { return ol.StringIndent(0) }
func (ol OLine) StringIndent(nIndent int) string {
	return helpers.Indent(nIndent) + ol.LineType() +
		" " +
		strings.Join(
			lambda.Map(
				ol.args,
				func(arg any) string {
					switch v := arg.(type) {
					case Astnode:
						return v.ContentString()
					case string:
						return v
					case []Line:
						return BoldBlue("[\n") + strings.Join(
							lambda.Map(v, func(line Line) string { return line.StringIndent(nIndent + 1) }),
							"\n",
						) + "\n" + helpers.Indent(nIndent) + BoldBlue("]")
					case func(Astnode) Astnode:
						return Yellow("𝘧")
					case func(Astnode, *ParseContext) Astnode:
						return Yellow("func (Astnode, *ParseContext) Astnode")
					default:
						return "?" + reflect.TypeOf(v).String() + "?"
					}
				},
			),
			BoldGreen(", "),
		)
	return Green("OLine(") + strings.Join(lambda.Map(
		ol.args,
		func(arg any) string {
			if s, ok := arg.(string); ok {
				return Red("\"") + s + Red("\"")
			} else if il, ok := arg.(ILine); ok {
				return il.String()
			} else if ol, ok := arg.(OLine); ok {
				return ol.String()
			} else if _, ok := arg.(func(Astnode) Astnode); ok {
				return Yellow("func (Astnode) Astnode")
			} else if _, ok := arg.(func(Astnode, *ParseContext) Astnode); ok {
				return Yellow("func (Astnode, *ParseContext) Astnode")
			} else {
				return Magenta("?")
			}
		},
	), Green(", ")) + Green(")")
}

// Since OLine is constructed with variadic ...any, `getArgs()` helps get those
// args in a structured way.
//
// Original coffee implementation allowed:
//   SOMEKEY: "<def>"
//   SOMEKEY: [ <rules> ]
// We disallow objects and thus `str` is more ambiguous.
// Thus, interpretation of str left to the caller (see ToRule() below):
// - O("CHOICE _"),            <- since no subrules, a nameless definition (nameOrDef is a def)
// - O("CHOICE", rules(...))   <- since rules, "CHOICE" must be considered a name (nameOrDef is a name)
// - O("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", func(it Astnode) Astnode { return ast.NewPattern(it) }),
//          <- ?
//
func (ol OLine) getArgs() (str string, rules []Line, attrs ParseOptions) {
	// unlike joeson.coffee, outside Rules() gives a [OLine, ...], not an actual OLine, therefore:
	for i, arg := range ol.args {
		// fmt.Printf("%d %v", i, arg)
		if i == 0 {
			if s, ok := arg.(string); ok {
				str = s
			} else if _, ok := arg.(Astnode); ok {
				// While Astnode is legit as a first arg of OLine,
				// ToRuleWithName is the expected call path in that case, and
				// thus we should panic
				panic("logic")
			} else {
				panic("OLine first argument must be a string, while reading: " + ol.String())
			}
		} else {
			if f, ok := arg.(func(Astnode) Astnode); ok {
				attrs.CbBuilder = func(z Astnode, _ *ParseContext) Astnode { return f(z) }
			} else if f, ok := arg.(func(Astnode, *ParseContext) Astnode); ok {
				attrs.CbBuilder = f
			} else if passedAttrs, ok := arg.(ParseOptions); ok {
				if passedAttrs.CbBuilder != nil {
					attrs.CbBuilder = passedAttrs.CbBuilder
				}
				attrs.SkipCache = passedAttrs.SkipCache
				attrs.SkipLog = passedAttrs.SkipLog
			} else if aRules, ok := arg.([]Line); ok {
				rules = aRules
				// if it's an array, in original impl.
				// it would be NAME: [...]. This means
				// previously filled `nameOrDef` is actually the "name" of
				// a rule and not a textual "def" of it.
			} else {
				fmt.Println("Ignoring arg " + strconv.Itoa(i) + ": " + fmt.Sprintf("%v", arg) + " in OLine " + fmt.Sprintf("%v", ol.args))
			}
		}
	}
	return
}

// Called only from NewRankFromLines
func (ol OLine) ToRuleWithIndex(parentRule Astnode, rankname string, index int, grammar *ast.Grammar) Astnode {
	if astnode, ok := ol.args[0].(Astnode); ok {
		astnode.GetGNode().Name = rankname
		return astnode // nothing to do when OLine is a handcompiled rule
	}
	nameOrDef, rules, attrs := ol.getArgs()
	// figure out the name for this rule
	// (note we unroll joeson.coffee's getRule() directly here)
	name := ""
	var rule Astnode = nil
	if rules == nil {
		name = rankname + "[" + strconv.Itoa(index) + "]"
		def := nameOrDef
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("Error in OLine named \"%s\" with def \"%s\":\n%v\n", name, def, e)
				// make it fail again for real this time
				grammar.ParseString(def, attrs)
			}
		}()
		// HACK temporarily halt trace
		oldTrace := Trace
		if Trace.SkipSetup {
			Trace.Stack = false
			Trace.Loop = false
		}
		ctx := NewParseContext(NewCodeStream(def), grammar, attrs)
		// ------------------------
		rule = grammar.Parse(ctx)
		// ------------------------
		Trace = oldTrace // end hack
	} else { // rules is []Line (non-nil)
		name = nameOrDef
		rule = NewRankFromLines(name, rules, grammar)
	}
	if rule.GetGNode().Rule != nil && rule.GetGNode().Rule != rule {
		panic("assert.ok((rule.rule == null) || rule.rule === rule)")
	}
	rule.GetGNode().Rule = rule
	if rule.GetGNode().Name != "" && rule.GetGNode().Name != name {
		panic("assert.ok((rule.name == null) || rule.name === name)")
	}
	rule.GetGNode().Name = name
	rule.GetGNode().SkipCache = attrs.SkipCache
	rule.GetGNode().SkipLog = attrs.SkipLog
	rule.GetGNode().CbBuilder = attrs.CbBuilder
	rule.GetGNode().Debug = attrs.Debug
	return rule
}

// called in one case only: when an I rule's definition is an OLine, as in:
// i "RANGE": o S(St('{'), R("_"), L("min",E(R("INT"))), R("_"), St(','), R("_"), L("max",E(R("INT"))), R("_"), St('}'))
//
// call path in that case:  NewRankFromLines() -> ILine.ToRule() -> ILine.getArgs() -> OLine.ToRuleWithName()
//
// compared to the original coffee implementation, this method is based on
// OLine.toRule with partial inlining of Line.getRules (just keeping the parts
// that matter for this use case)
func (ol OLine) ToRuleWithName(parentRule Astnode, name string) Astnode {
	// since getArgs() gets (nameOrDef, rules, attrs) and we need none of that,
	// let's skip that call. We actually require the 1st arg to be an Astnode.
	rule := ol.args[0].(Astnode) // optionally may also check for attrs in subsequent args
	rule.GetGNode().Parent = parentRule
	return rule
}
