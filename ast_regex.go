package joeson

import (
	"grepsuzette/joeson/helpers"
	"regexp"
	"strings"
)

type Regex struct {
	*GNode
	reStr string
	re    regexp.Regexp
}

func NewRegexFromString(sRegex string) *Regex {
	if compiledRegexp, e := regexp.Compile("(" + sRegex + ")"); e != nil {
		panic("Invalid regex: " + sRegex)
	} else {
		re := &Regex{NewGNode(), sRegex, *compiledRegexp}
		re.GNode.Node = re
		return re
	}
}

func NewRegex(it Ast) *Regex {
	return NewRegexFromString(joinNativeArrayOfNativeString(it))
}
func NewRegexCharClass(it Ast) *Regex {
	return NewRegexFromString("[" + joinNativeArrayOfNativeString(it) + "]")
}

func (re *Regex) GetGNode() *GNode { return re.GNode }
func (re *Regex) ContentString() string {
	// below /g is purely for output conformance to original coffee impl.
	return magenta("/" + re.re.String() + "/g")
}
func (re *Regex) HandlesChildLabel() bool { return false }
func (re *Regex) Prepare()                {}
func (re *Regex) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		if didMatch, sMatch := ctx.Code.MatchRegexp(re.re); !didMatch {
			return nil
		} else {
			return NewNativeString(sMatch)
		}
	}, re)(ctx)
}

func joinstr(a []NativeString, join string) string {
	a2 := helpers.AMap(a, func(ns NativeString) string { return ns.Str })
	return strings.Join(a2, join)
}

func joinNativeArrayOfNativeString(node Ast) string {
	switch node.(type) {
	case *NativeArray:
		var b strings.Builder
		for _, elt := range node.(*NativeArray).Array {
			switch elt.(type) {
			case NativeString:
				b.WriteString(elt.(NativeString).Str)
			default:
				panic("expected native string")
			}
		}
		return b.String()
	default:
		panic("assert")
	}
}
func (re *Regex) ForEachChild(f func(Ast) Ast) Ast {
	// no children defined for Ref, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	re.GetGNode().Rules = ForEachChild_InRules(re, f)
	return re
}
