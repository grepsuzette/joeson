package joeson

import (
	"github.com/grepsuzette/joeson/helpers"
)

// str is a simple Parser that tries to parse the string it was built with.
// Let p := newStr("foo").
// p.Parse("fool") will parse as NewNativeString("foo"),
// p.Parse("fbar") will fail.
//
// See also NativeString which is a terminal node, not a Parser.
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

// Given a rule such as o(named("interpreted_string_lit", "'\"' ( !(?'"+`"`+"') (unicode_value | byte_value) )* '\"'"))
// the part `(?")` is a lookahead.
// `newLookahead(Ast)` receives a NativeMap{expr:`"`} in that case, whereas it
// should be built with an object satisfying Parser.
// Hence the call `la := &lookahead{gn, newStrFromAst(ast)}`, which builds the `str` parser for `"`.
func newStrFromAst(ast Ast) str {
	switch v := ast.(type) {
	case NativeMap:
		// try to convert to Str iif it has only one key
		keys := v.Keys()
		if len(keys) < 1 {
			panic("assert Parser expected, got NativeMap but it's got more than one key so can not convert to Str: " + v.ContentString())
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
					return newStr(stringFromNativeArray(v))
				default:
					panic("Could not create a Parser from NativeMap " + v.ContentString())
				}
			}
		}
	case NativeString:
		return newStr(v.Str)
	case str:
		return v
	default:
		panic("Could not create str from " + v.ContentString())
	}
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
	s.GetGNode().rules = ForEachChild_InRules(s, f)
	return s
}
