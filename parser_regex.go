package joeson

import (
	"regexp"
	"strings"
)

// avoid regexes with PEG in general, regexes are greedy and this can
// create ambiguity and buggy grammars. As a special case, character classes are OK.
// Regexes can be used to optimize but again avoid them unless
// you know what you're doing.
type regex struct {
	*Attr
	*rule
	reStr string
	re    regexp.Regexp
}

func newRegexFromString(sRegex string) *regex {
	if compiledRegexp, e := regexp.Compile("(" + sRegex + ")"); e != nil {
		panic("Invalid regex: " + sRegex)
	} else {
		re := &regex{newAttr(), newRule(), sRegex, *compiledRegexp}
		re.rule.node = re
		return re
	}
}

func (re *regex) gnode() *rule { return re.rule }
func (re *regex) String() string {
	return Magenta("/" + strings.NewReplacer("\r", "\\r", "\n", "\\n", "\t", "\\t").Replace(re.re.String()) + "/g")
}
func (re *regex) HandlesChildLabel() bool { return false }
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
	re.rules = ForEachChildInRules(re, f)
	return re
}
