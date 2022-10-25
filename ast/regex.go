package ast

import "regexp"
import "grepsuzette/joeson/lambda"
import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"
import "strings"

type Regex struct {
	*GNode
	reStr string
	re    regexp.Regexp
}

func NewRegexFromString(sRegex string) *Regex {
	if compiledRegexp, e := regexp.Compile("(" + sRegex + ")"); e != nil {
		panic("Invalid regex: " + sRegex)
	} else {
		return &Regex{NewGNode(), sRegex, *compiledRegexp}
	}
}

// wtf TODO check those 2 funcs
// Astnode must be a NativeArray of NativeString
func NewRegex(it Astnode) *Regex {
	return NewRegexFromString(joinNativeArrayOfNativeString(it))
}
func NewRegexCharClass(it Astnode) *Regex {
	return NewRegexFromString("[" + joinNativeArrayOfNativeString(it) + "]")
}

func (re *Regex) GetGNode() *GNode { return re.GNode }
func (re *Regex) ContentString() string {
	return ShowLabelOrNameIfAny(re) + BoldRed("/") + Magenta(re.re.String()) + BoldRed("/")
}
func (re *Regex) HandlesChildLabel() bool { return false }
func (re *Regex) Labels() []string        { return re.GNode.Labels() }
func (re *Regex) Captures() []Astnode     { return re.GNode.Captures() }
func (re *Regex) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
		if didMatch, sMatch := ctx.Code.MatchRegexp(re.re); !didMatch {
			return nil
		} else {
			return NewNativeString(sMatch)
		}
	}, re)(ctx)
}

func joinstr(a []NativeString, join string) string {
	a2 := lambda.Map(a, func(ns NativeString) string { return ns.Str })
	return strings.Join(a2, join)
}

func joinNativeArrayOfNativeString(node Astnode) string {
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
		panic("expected a NativeArray containing NativeString elements")
	}
}
func (re *Regex) ForEachChild(f func(Astnode) Astnode) Astnode {
	// no children defined in coffee
	return re
}
