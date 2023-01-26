package joeson

import (
	"grepsuzette/joeson/helpers"
	"reflect"
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
			default:
				// fmt.Printf("%s", reflect.TypeOf(v).String())
				panic("assert")
			}
		}
	}
	return
}

// sanitize arbitrary content into a Line
func rule2line(x any) Line {
	switch v := x.(type) {
	case string:
		return newSLine(v)
	case sLine:
		return O(v.Str)
	case Ast:
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
//                 In this implementation it means Line, amongst:
//                   SLine (for string), ALine, OLine
// parentRule: The actual parent Rule instance
// attrs:      {cb,...}, extends the result
// opts:       Parse time options

// see line/README.md # internals
func getRule(rank_ *rank, name string, line Line, parentRule Ast, attrs ParseOptions, opts TraceOptions, lazyGrammar *helpers.Lazy[*Grammar]) Ast {
	var retAst Ast
	// fmt.Println("getRule name=" + name + " eflect.TypeOf(line).String()):" + reflect.TypeOf(line).String())
	switch v := line.(type) {
	case ALine:
		retAst = rankFromLines(v.Array, name, GrammarOptions{TraceOptions: opts, LazyGrammar: lazyGrammar})
	case cLine:
		retAst = v.Ast
		retAst.GetGNode().Name = name
	case ILine:
		panic("assert") // ILine is impossible here
	case OLine:
		retAst = v.toRule(rank_, parentRule, oLineByIndexOrName{name: name}, opts, lazyGrammar)
		if retAst.GetGNode().Name == "" {
			retAst.GetGNode().Name = name
		}
	case sLine:
		// HACK: temporarily halt trace when SkipSetup
		var skipSetup bool = opts.SkipSetup
		var oldTrace TraceOptions
		if skipSetup {
			oldTrace = opts
			opts = Mute()
		}
		// parse the string
		// a grammar like joeson_handcompiled is needed for that,
		gm := lazyGrammar.Get() // uses Lazy to get the grammar in cache or build it
		if x, error := gm.parseOrFail(
			newParseContext(NewCodeStream(v.Str), gm.NumRules, attrs, opts),
		); error == nil {
			retAst = x
		} else {
			panic(error)
		}
		retAst.GetGNode().Name = name
		if skipSetup {
			opts = oldTrace
		}
	default:
		panic("unrecog type " + reflect.TypeOf(line).String())
	}
	rule := retAst.GetGNode()
	if rule.Rule != nil && !IsRule(retAst) {
		panic("assert")
	}
	rule.Rule = retAst
	if rule.Name != "" && rule.Name != name {
		panic("assert")
	}
	rule.SkipCache = attrs.SkipCache
	rule.SkipLog = attrs.SkipLog
	rule.CbBuilder = attrs.CbBuilder
	rule.Debug = attrs.Debug
	return retAst
}
