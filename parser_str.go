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
	Attr
	*gnodeimpl
	Str string
}

func newStr(s string) str {
	str := str{newAttr(), newGNode(), s}
	str.gnodeimpl.capture = false
	str.gnodeimpl.node = str
	return str
}

// used in newLookahead()
func newStrFromAst(ast Ast) str {
	switch v := ast.(type) {
	case NativeMap:
		// convert to str only when there is a single key
		keys := v.Keys()
		if len(keys) < 1 {
			panic("assert Parser expected, got NativeMap but it's got more than one key so can not convert to Str: " + v.String())
		} else {
			if ast, ok := v.GetExists(keys[0]); !ok {
				panic("should not happen")
			} else {
				switch w := ast.(type) {
				case NativeString:
					return newStr(w.Str)
				case str:
					return w
				case *NativeArray:
					return newStr(w.Concat())
				default:
					panic("Could not create a Parser from NativeMap " + v.String())
				}
			}
		}
	case NativeString:
		return newStr(v.Str)
	case str:
		return v
	default:
		panic("Could not create str from " + v.String())
	}
}

func (s str) gnode() *gnodeimpl       { return s.gnodeimpl }
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

func (s str) ForEachChild(f func(Parser) Parser) Parser {
	// no children defined for Str, but GNode has:
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	s.gnode().rules = ForEachChildInRules(s, f)
	return s
}
