package line

import (
	"grepsuzette/joeson/ast"
	. "grepsuzette/joeson/core"

	//. "grepsuzette/joeson/colors"
	"reflect"
)

// ALine are simply []Line
// SLine are transitory, and when parsed become OLine
// The OLine, as non-terminal, becomes a CLine (CLine wraps an Astnode) when parsed
// The ILine, terminal, holds either OLine, or CLine (an Astnode, including Str which is an Astnode)

// Line interface is just a way to have []Line really
// These are a system to enter rules of a grammar
// in a code-like fashion (as a tree, rather than linearly).
type Line interface {
	LineType() string // i, o, a, s, c
	Content() Line    // Sline, OLine, ALine, CLine (containing an Astnode)...
	String() string
	StringIndent(nIndent int) string // same as String(), but indenting with `nIdent` levels (for nested rules)
}

// common functions callable by both ILine nd OLine

func lineInit(origArgs []any) (name string, lineContent Line, attrs ParseOptions) {
	for i, arg := range origArgs {
		if i == 0 {
			switch v := arg.(type) {
			case NamedRule:
				name = v.Name
				lineContent = rule2line(v.Line)
			default:
				name = ""
				lineContent = rule2line(v)
			}
		} else {
			switch v := arg.(type) {
			case func(it Astnode) Astnode:
				attrs.CbBuilder = func(x Astnode, _ *ParseContext, _ Astnode) Astnode {
					return v(x)
				}
			case func(_ Astnode, _ *ParseContext) Astnode:
				attrs.CbBuilder = func(x Astnode, ctx *ParseContext, _ Astnode) Astnode {
					return v(x, ctx)
				}
			case func(_ Astnode, _ *ParseContext, _ Astnode) Astnode:
				attrs.CbBuilder = v
			case ParseOptions:
				attrs = v
			default:
				panic("unfaweif293ager8")
			}
		}
	}
	return
}

func rule2line(x any) Line {
	switch v := x.(type) {
	case ALine:
		return v
	case string:
		return NewSLine(v)
	case OLine:
		return v
	case Astnode:
		return NewCLine(v)
	case CLine:
		return v
	case SLine:
		// TODO let's try directly to wrap it in O?
		return O(v.Str)
		// panic("Cannot accept a string directly here for now, please wrap it in an OLine: " + v.Str)
	case ILine:
		panic("impossible")
	case []Line:
		panic("unused me thinks it would be ALine instead")
		return NewALine(v)
	default:
		panic("impossible type in rule2line: " + reflect.TypeOf(x).String())
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
func getRule(grammar *ast.Grammar, name string, line Line, parentRule Astnode, attrs ParseOptions) Astnode {
	var ast Astnode
	// fmt.Println("getRule name=" + name + reflect.TypeOf(line).String())
	switch v := line.(type) {
	case ALine:
		// fmt.Println("getRule ALine name=" + name)
		ast = NewRankFromLines(name, v.Array, grammar)
	case CLine:
		// fmt.Println("getRule CLine name=" + name)
		ast = v.ast
		ast.GetGNode().Name = name
	case ILine:
		panic("ILine is impossible here")
	case OLine:
		ast = v.ToRule(grammar, parentRule, OLineByIndexOrName{name: name})
	case SLine:
		var ctx *ParseContext
		// defer func() {
		// 	if e := recover(); e != nil {
		// 		fmt.Printf("Error in rule %s, that is %s:\n%v\n", Cyan(name), Magenta(v.Str), e)
		// 	}
		// }()
		// TODO, can surround with halt trace instructions as in coffee impl
		ctx = NewParseContext(NewCodeStream(v.Str), grammar, attrs)
		ast = grammar.Parse(ctx)
	default:
		panic("unrecog type " + reflect.TypeOf(line).String())
	}
	rule := ast.GetGNode()
	if rule.Rule != nil && !IsRule(ast) {
		panic("assert")
	}
	rule.Rule = ast
	if rule.Name != "" && rule.Name != name {
		panic("assert")
	}
	// if attrs != nil {
	rule.SkipCache = attrs.SkipCache
	rule.SkipLog = attrs.SkipLog
	rule.CbBuilder = attrs.CbBuilder
	rule.Debug = attrs.Debug
	// }
	return ast
}
