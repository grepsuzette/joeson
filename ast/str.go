package ast

import "grepsuzette/joeson/helpers"
import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"

// This does not represent a native string
// but a Str in the joeson grammar.
// see NativeString which is a terminal element
// both are similar otherwise
type Str struct {
	*GNode
	str string
}

func NewStr(s string) Str {
	g := NewGNode()
	g.Capture = false
	return Str{g, s}
}
func (str Str) GetGNode() *GNode        { return str.GNode }
func (str Str) Prepare()                {}
func (str Str) HandlesChildLabel() bool { return false }
func (str Str) Labels() []string        { return str.GNode.Labels() }
func (str Str) Captures() []Astnode     { return str.GNode.Captures() }
func (str Str) ContentString() string {
	return ShowLabelOrNameIfAny(str) + Green("'"+helpers.Escape(str.str)+"'")
}
func (str Str) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
		if didMatch, sMatch := ctx.Code.MatchString(str.str); didMatch {
			// a string is not a terminal element
			// so return NativeString.
			// TODO Very likely one of those classes can be taken away
			return NewNativeString(sMatch)
		} else {
			return nil
		}
	}, str)(ctx)
}
func (str Str) ForEachChild(f func(Astnode) Astnode) Astnode {
	// no children defined in coffee
	return str
}
