package joeson

import (
	"grepsuzette/joeson/helpers"
)

// str does not represent a native string
// but rather a `str` in the joeson grammar.
// see NativeString which is a terminal element (a native string in js)
type str struct {
	*GNodeImpl
	Str string
}

func newStr(s string) str {
	str := str{NewGNode(), s}
	str.GNodeImpl.capture = false
	str.GNodeImpl.node = str
	return str
}
func (s str) GetGNode() *GNodeImpl    { return s.GNodeImpl }
func (s str) Prepare()                {}
func (s str) HandlesChildLabel() bool { return false }
func (s str) ContentString() string {
	return green("'" + helpers.Escape(s.Str) + "'")
}
func (s str) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Parser) Ast {
		if didMatch, sMatch := ctx.Code.MatchString(s.Str); didMatch {
			// a string is not a terminal element
			// so return NativeString.
			return NewNativeString(sMatch)
		} else {
			return nil
		}
	}, s)(ctx)
}
func (s str) ForEachChild(f func(Parser) Parser) Parser {
	// no children defined for Str, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	s.GetGNode().Rules = ForEachChild_InRules(s, f)
	return s
}
