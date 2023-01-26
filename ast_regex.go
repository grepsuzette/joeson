package joeson

import (
	"regexp"
)

type regex struct {
	*GNode
	reStr string
	re    regexp.Regexp
}

func newRegexFromString(sRegex string) *regex {
	if compiledRegexp, e := regexp.Compile("(" + sRegex + ")"); e != nil {
		panic("Invalid regex: " + sRegex)
	} else {
		re := &regex{NewGNode(), sRegex, *compiledRegexp}
		re.GNode.Node = re
		return re
	}
}

func (re *regex) GetGNode() *GNode { return re.GNode }
func (re *regex) ContentString() string {
	// below /g is purely for output conformance to original coffee impl.
	return magenta("/" + re.re.String() + "/g")
}
func (re *regex) HandlesChildLabel() bool { return false }
func (re *regex) Prepare()                {}
func (re *regex) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		if didMatch, sMatch := ctx.Code.MatchRegexp(re.re); !didMatch {
			return nil
		} else {
			return NewNativeString(sMatch)
		}
	}, re)(ctx)
}

func (re *regex) ForEachChild(f func(Ast) Ast) Ast {
	// no children defined for Ref, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	re.GetGNode().Rules = ForEachChild_InRules(re, f)
	return re
}
