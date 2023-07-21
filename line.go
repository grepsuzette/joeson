package joeson

import (
	"fmt"
	"os"
	"reflect"

	"github.com/grepsuzette/joeson/helpers"
)

type Line interface {
	lineType() string                // i, o, a, s, c
	stringIndent(nIndent int) string // indent with `nIdent` levels (for nested rules)
}

/* -- follow some common functions used by ILine & OLine -- */

// The functions `I(a ...any)` and `O(a ...any)` both call `lineInit(a)`
// to help destructuring `a` into a name, content (Line) and options. This
// is where `Named()` gets decomposed if it was used. This is also
// where parsing callbacks make their way into ParseOptions.
func lineInit(origArgs []any) (name string, lineContent Line, attrs ParseOptions) {
	for i, arg := range origArgs {
		if i == 0 {
			switch v := arg.(type) {
			case NamedRule:
				name = v.name
				lineContent = rule2line(v.line)
			default:
				name = ""
				lineContent = rule2line(v)
			}
		} else {
			switch v := arg.(type) {
			case func(it Ast) Ast:
				attrs.CbBuilder = func(x Ast, _ *ParseContext, _ Ast) Ast {
					return v(x)
				}
			case func(Ast, *ParseContext) Ast:
				attrs.CbBuilder = func(x Ast, ctx *ParseContext, _ Ast) Ast {
					return v(x, ctx)
				}
			case func(Ast, *ParseContext, Ast) Ast:
				attrs.CbBuilder = v
			case ParseOptions:
				attrs = v
			case string:
				fmt.Printf(
					"Error in grammar: O (or I) called lineInit with %v\nSo the second parameter was a string: %s\nRight now this syntax is not supported\nPlease fix your grammar",
					origArgs,
					v,
				)
				os.Exit(1)
			case OLine:
				fmt.Printf("Error in grammar: lineInit called with OLine %s\n", v.stringIndent(0))
				panic("assert")
			case []Line:
				fmt.Printf("Error in grammar: Arrays of rules are expected to arrive as the 1st argument i.e. i=0) but here it came in position i=%d. Here is the faulty rule, you probably forgot to wrap the rule into named():\n%s\n.", i, summarizeRule(origArgs, 2))
				os.Exit(1)
			default:
				fmt.Printf("%s\n", reflect.TypeOf(v).String())
				panic("assert")
			}
		}
	}
	return
}

// instead of displaying a full tree of rules when there is an error,
// just retain the first few children.
func summarizeRule(args []any, max int) string {
	s := ""
	for i := range args {
		if i > max {
			return s + "...\n"
		} else {
			switch v := args[i].(type) {
			case string:
				s += `"` + v + `"`
			case []Line:
				s += "["
				for j, w := range v {
					if j > max {
						s += "...\n"
						break
					} else {
						s += w.stringIndent(i+1) + "\n"
					}
				}
				s += "]"
			default:
				s += "?unhandled_type:" + reflect.TypeOf(v).String() + "', "
				// s += fmt.Sprintf("%v", v)
			}
		}
		s += ", "
	}
	return s
}

// sanitize arbitrary content into a Line
func rule2line(x any) Line {
	switch v := x.(type) {
	case string:
		return newSLine(v)
	case sLine:
		return O(v.Str)
	case Parser:
		return newCLine(v)
	case ALine:
		return v
	case OLine:
		return v
	case ILine:
		panic("assert")
	case []Line:
		panic("assert") // because it should have been ALine
	default:
		panic("assert")
		// panic("impossible type in rule2line: " + reflect.TypeOf(x).String())
	}
}

// name:       The final and correct name for this rule
// rule:       A rule-like object
//                 In coffee it means string, array, object (map) or oline
//                 In this implementation it means Line, among:
//                   SLine (for string), ALine, OLine
// parentRule: The actual parent Rule instance
// attrs:      {cb,...}, extends the result
// opts:       Parse time options

// see line/README.md # internals
func getRule(rank_ *rank, name string, line Line, parentRule Parser, attrs ParseOptions, opts TraceOptions, lazyGrammar *helpers.Lazy[*Grammar]) Parser {
	var answer Parser
	// if name == "decimal_digit" {
	// 	fmt.Printf("getRule name=%s reflect.TypeOf(line).String())=%s attrs=%s\n", name, reflect.TypeOf(line).String(), attrs)
	// }
	switch v := line.(type) {
	case ALine:
		answer = rankFromLines(v.Array, name, GrammarOptions{TraceOptions: opts, LazyGrammar: lazyGrammar})
	case cLine:
		answer = v.Parser
		answer.SetRuleName(name)
	case ILine:
		panic("assert") // ILine is impossible here
	case OLine:
		v.attrs = attrs
		answer = v.toRule(rank_, parentRule, oLineByIndexOrName{name: name}, opts, lazyGrammar)
		answer.SetRuleNameWhenEmpty(name)
		// answer.(gnode).gnode().ParseOptions = attrs
	case sLine:
		// temporarily halt trace when SkipSetup
		traceOptions := opts
		if opts.SkipSetup {
			traceOptions.Loop = false
			traceOptions.Stack = false
		}
		// parse the string
		// a grammar like joeson_handcompiled is needed for that,
		gm := lazyGrammar.Get() // uses Lazy to get the grammar in cache or build it
		ctx := newParseContext(NewCodeStream(v.Str), gm.numrules, traceOptions).setParseOptions(attrs)
		ast := gm.Parse(ctx)
		if IsParseError(ast) {
			panic(ast.(ParseError).String())
		} else {
			answer = ast.(Parser)
		}
		answer.SetRuleName(name)
	default:
		panic("unrecog type " + reflect.TypeOf(line).String())
	}
	rule := answer.gnode()
	if rule.rule != nil && !IsRule(answer) {
		panic("assert")
	}
	rule.rule = answer
	if rule.name != "" && rule.name != name {
		panic("assert")
	}
	// if name == "decimal_digit" {
	// 	fmt.Printf(" name=%s attrs rule.rule.Debug=%s rule.Debug=%s attrs.Debug=%s\n", name, reflect.TypeOf(line).String(), rule.Debug, rule.Debug, attrs.Debug)
	// }
	// rule.SkipCache = attrs.SkipCache
	// rule.SkipLog = attrs.SkipLog
	rule.CbBuilder = attrs.CbBuilder
	// rule.Debug = attrs.Debug
	return answer
}
