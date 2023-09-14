package joeson

import (
	"github.com/grepsuzette/joeson/helpers"
)

// str is a parser that must match a string exactly.
// It does not capture by default.
// ```
// p := newStr("foo").
// p.Parse("fool") -> NewNativeString("foo"),
// p.Parse("fbar") -> nil.
// ```
type str struct {
	*Attr
	*rule
	Str string
}

func newStr(s string) str {
	str := str{newAttr(), newRule(), s}
	str.rule.capture = false
	str.rule.node = str
	return str
}

func (s str) getRule() *rule          { return s.rule }
func (s str) prepare()                {}
func (s str) handlesChildLabel() bool { return false }
func (s str) String() string {
	return Green("'" + helpers.Escape(s.Str) + "'")
}

func (s str) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		if didMatch, sMatch := ctx.Code.MatchString(s.Str); didMatch {
			return NewNativeString(sMatch)
		} else {
			return nil
		}
	}, s)(ctx)
}

func (s str) forEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	s.rules = ForEachChildInRules(s, f)
	return s
}
