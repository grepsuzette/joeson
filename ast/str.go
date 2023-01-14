package ast

import (
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"

	. "grepsuzette/joeson/colors"
)

// Str does not represent a native string
// but rather a `Str` in the joeson grammar.
// see NativeString which is a terminal element (a native string in js)
type Str struct {
	*GNode
	Str string
}

func NewStr(s string) Str {
	str := Str{NewGNode(), s}
	str.GNode.Capture = false
	str.GNode.Node = str
	return str
}
func (str Str) GetGNode() *GNode        { return str.GNode }
func (str Str) Prepare()                {}
func (str Str) HandlesChildLabel() bool { return false }
func (str Str) ContentString() string {
	return Green("'" + helpers.Escape(str.Str) + "'")
}
func (str Str) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		if didMatch, sMatch := ctx.Code.MatchString(str.Str); didMatch {
			// a string is not a terminal element
			// so return NativeString.
			return NewNativeString(sMatch)
		} else {
			return nil
		}
	}, str)(ctx)
}
func (str Str) ForEachChild(f func(Ast) Ast) Ast {
	// no children defined for Str, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	str.GetGNode().Rules = ForEachChild_InRules(str, f)
	return str
}
