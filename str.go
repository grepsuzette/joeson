package main

import "grepsuzette/joeson/helpers"

// This does not represent a native string
// but a Str in the joeson grammar.
// see NativeString which is a terminal element
// both are similar otherwise
type Str struct {
	GNode
	str string
}

func NewStr(s string) Str {
	g := NewGNode()
	g.capture = false
	return Str{g, s}
}
func (str Str) GetGNode() GNode       { return str.GNode }
func (str Str) Prepare()              {}
func (str Str) HandlesChildLabel()    { return false }
func (str Str) Labels() []string      { return str.GNode.Labels() }
func (str Str) Captures() []astnode   { return str.GNode.Captures() }
func (str Str) ContentString() string { return Green("'" + helpers.Escape(str.str) + "'") }
func (str Str) Parse(ctx *ParseContext) astnode {
	return _wrap(func(ctx, _) astnode {
		if didMatch, sMatch := ctx.code.MatchString(str.str); didMatch {
			// a string is not a terminal element
			// so return NativeString.
			// Very likely one of those classes can be taken away
			// but at first it's like this
			return NewNativeString(sMatch)
		}
	})(ctx, str)
}
