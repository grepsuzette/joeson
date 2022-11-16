package line

import "fmt"
import "grepsuzette/joeson/ast"
import . "grepsuzette/joeson/core"

// Line interface is just a way to have []Line really
// These are a system to enter rules of a grammar
// in a code-like fashion (as a tree, rather than linearly).
type Line interface {
	Args() []any
	String() string
	StringIndent(nIndent int) string // same as String(), but indenting with `nIdent` levels (for nested rules)
	LineType() string
	IsO() bool
}

// common functions callable by both ILine nd OLine

// weird but hopefully faithful translation
// TODO it's probably nonsense, get rid of this after testing
func lineInit(args []any) []any {
	if len(args) >= 1 {
		a := make([]any, len(args))
		copy(a, args)
		return a
	} else {
		return []any{}
	}
}

// name:       The final and correct name for this rule
// rulelike:   A rule-like object
//                 In coffee it means str, array, object (map) or oline
//                 In this implementation it means []any, amongst:
//                   String, []Line, Oline
// parentRule: The actual parent Rule instance
// attrs:      {cb,...}, extends the result
// opts:       Parse time options
func getRule(grammar *ast.Grammar, name string, rulelike []any, parentRule Astnode, attrs ParseOptions) Astnode {
	var ast Astnode
	switch v := rulelike.(type) {
	case string:
		var ctx *ParseContext
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("Error in rule %s: %s:\n%v\n", name, v.Str, e)
				// make it fail again for real this time
				grammar.Parse(ctx)
			}
		}()
		// TODO, can surround with halt trace instructions as in coffee impl
		ctx = NewParseContext(NewCodeStream(v), grammar, attrs)
		ast = grammar.Parse(ctx)
	case []Line:
		ast = NewRankFromLines(name, v, grammar)
	case OLine:
		ast = v.ToRule(grammar, parentRule, OLineByIndexOrByName{name: name})
	default:
		panic("unrecog type " + reflect.TypeOf(rulelike).String())
	}
	rule = ast.GetGNode()
	// shouldn't rule be directly a GNode?
	if rule.Rule != nil && !rule.IsRule() {
		panic("fjai3289")
	}
	rule.Rule = rule
	if rule.Name != "" && rule.Name != name {
		panic("fa8332")
	}
	if attrs != nil {
		rule.SkipCache = attrs.SkipCache
		rule.SkipLog = attrs.SkipLog
		rule.CbBuilder = attrs.CbBuilder
		rule.Debug = attrs.Debug
	}
	return ast
}

// returns {rule:rule, attrs:{cb,skipCache,skipLog,...}}
func getArgs(line Line) (rule Astnode, attrs ParseOptions) {
	// [rule, rest...] = @args
	rule := line.Args()[0].(Astnode)
	var rest []any = line.Args()[1:]
	for _, next := range rest {
		switch v := next.(type) {
		case func(Astnode) Astnode:
			attrs.CbBuilder = v
		case ParseOptions:
			if v.CbBuilder != nil {
				attrs.CbBuilder = v.CbBuilder
			}
			attrs.SkipCache = v.SkipCache
			attrs.SkipLog = v.SkipLog
		default:
			panic("fawehf")
		}
	}
	return
}
