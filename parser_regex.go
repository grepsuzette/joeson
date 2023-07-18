package joeson

import (
	"regexp"
)

type regex struct {
	Attributes
	*gnodeimpl
	reStr string
	re    regexp.Regexp
}

func newRegexFromString(sRegex string) *regex {
	if compiledRegexp, e := regexp.Compile("(" + sRegex + ")"); e != nil {
		panic("Invalid regex: " + sRegex)
	} else {
		re := &regex{Attributes{}, NewGNode(), sRegex, *compiledRegexp}
		re.gnodeimpl.node = re
		return re
	}
}

func (re *regex) gnode() *gnodeimpl { return re.gnodeimpl }
func (re *regex) String() string {
	// below /g is purely for output conformance to original coffee impl.
	return magenta("/" + re.re.String() + "/g")
}
func (re *regex) handlesChildLabel() bool { return false }
func (re *regex) prepare()                {}
func (re *regex) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
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
	re.gnode().rules = ForEachChildInRules(re, f)
	return re
}
