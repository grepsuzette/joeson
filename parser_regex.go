package joeson

import (
	"regexp"
)

type regex struct {
	*GNodeImpl
	reStr string
	re    regexp.Regexp
}

func newRegexFromString(sRegex string) *regex {
	if compiledRegexp, e := regexp.Compile("(" + sRegex + ")"); e != nil {
		panic("Invalid regex: " + sRegex)
	} else {
		re := &regex{NewGNode(), sRegex, *compiledRegexp}
		re.GNodeImpl.node = re
		return re
	}
}

func (re *regex) GetGNode() *GNodeImpl { return re.GNodeImpl }
func (re *regex) ContentString() string {
	// below /g is purely for output conformance to original coffee impl.
	return magenta("/" + re.re.String() + "/g")
}
func (re *regex) HandlesChildLabel() bool { return false }
func (re *regex) Prepare()                {}
func (re *regex) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Parser) Ast {
		if didMatch, sMatch := ctx.Code.MatchRegexp(re.re); !didMatch {
			return nil
		} else {
			return NewNativeString(sMatch)
		}
	}, re)(ctx)
}

func (re *regex) ForEachChild(f func(Parser) Parser) Parser {
	// no children defined for Ref, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	re.GetGNode().rules = ForEachChild_InRules(re, f)
	return re
}
