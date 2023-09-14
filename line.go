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

// Functions jointly used by ILine and OLine

// lineInit
// Called by `I(...any)` and `O(...any)`.
// Helps destructuring their arguments into name, content and options.
// Unpacks `Named()`.
// Assigns parse functions onto return parseOptions.
// Collects individual ...ParseOption into returned parseOptions
func lineInit(origArgs []any) (name string, lineContent Line, attrs *parseOptions) {
	attrs = newParseOptions()
	for i, arg := range origArgs {
		if i == 0 {
			// Named rule can appear in first position
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
			// support 3 variants of parse functions
			case func(it Ast) Ast:
				attrs.cb = func(x Ast, _ *ParseContext, _ Ast) Ast {
					return v(x)
				}
			case func(Ast, *ParseContext) Ast:
				attrs.cb = func(x Ast, ctx *ParseContext, _ Ast) Ast {
					return v(x, ctx)
				}
			case func(Ast, *ParseContext, Ast) Ast:
				attrs.cb = v
			case ParseOption:
				// A separate option, such as `Debug{true}`
				attrs = v.apply(attrs)
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
				panic(fmt.Sprintf("Error in grammar: unexpected type received in lineInit: %s\n", reflect.TypeOf(v).String()))
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

// getRule is directly transposed from coffeescript.
// See docs/internals.md
//
// name:       The final and correct name for this rule
// rule:       A rule-like object
//
//	In coffee it means string, array, object (map) or oline
//	In this implementation it means Line, among:
//	  SLine (for string), ALine, OLine
//
// parentRule: The actual parent Rule instance
// attrs:      {cb,...}, extends the result
// opts:       Parse time options
func getRule(
	rank_ *rank,
	name string,
	line Line,
	parentRule Parser,
	attrs *parseOptions,
	opts *TraceOptions,
	lazyGrammar *helpers.Lazy[*Grammar],
) Parser {
	var answer Parser
	// 	fmt.Printf("getRule name=%s reflect.TypeOf(line).String())=%s attrs=%s\n", name, reflect.TypeOf(line).String(), attrs)
	switch v := line.(type) {
	case ALine:
		answer = rankFromLines(v.Array, name, GrammarOptions{TraceOptions: opts, LazyGrammar: lazyGrammar})
	case cLine:
		answer = v.Parser
		answer.getRule().name = name
	case ILine:
		panic("assert") // ILine is impossible here
	case OLine:
		v.parseOptions = attrs
		answer = v.toRule(rank_, parentRule, oLineNaming{name: name}, opts, lazyGrammar)
		if answer.getRule().name == "" {
			answer.getRule().name = name
		}
		// answer.(rule).getRule().ParseOptions = attrs
	case sLine:
		// temporarily halt trace when SkipSetup
		var traceOptions *TraceOptions
		if opts.SkipSetup {
			traceOptions = opts.Copy()
			traceOptions.Loop = false
			traceOptions.Stack = false
		} else {
			traceOptions = opts
		}
		// parse the string. A grammar like joeson_handcompiled is needed for that,
		gm := lazyGrammar.Get() // uses Lazy to get the grammar in cache or build it
		ctx := newParseContext(NewRuneStream(v.Str), gm.numrules, traceOptions)
		// Do not use parseoptions during compilation of an sLine
		// NO: ctx = ctx.setParseOptions(attrs)
		// YES: We instead want to store the option inside the rule;
		//      that rule being the compiled answer.
		ast := gm.parse(ctx)
		if IsParseError(ast) {
			panic(ast.(ParseError).String())
		} else {
			answer = ast.(Parser)
		}
		answer.getRule().name = name
		answer.getRule().parseOptions = attrs
	default:
		panic("unrecog type " + reflect.TypeOf(line).String())
	}
	rule := answer.getRule()
	if rule.parser != nil && !IsRule(answer) {
		panic("assert")
	}
	rule.parser = answer
	if rule.name != "" && rule.name != name {
		panic("assert")
	}
	// fmt.Printf(" name=%s attrs rule.rule.Debug=%s rule.Debug=%s attrs.Debug=%s\n", name, reflect.TypeOf(line).String(), rule.Debug, rule.Debug, attrs.Debug)
	rule.parseOptions.cb = attrs.cb
	return answer
}
