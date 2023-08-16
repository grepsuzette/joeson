package joeson

import (
	"fmt"
)

// Ast is the result type of a Parse() operation by a grammar.
//
// The return can be:
//
//   - `nil` if a parser failed (in a way where backtracking will happen),
//
//   - `ParseError` if a parsing must definitely fail (ParseError is an Ast).
//
//   - anything else you want that implements Ast. For that, you will
//     need a rule with a callback.
//
//     For joeson grammar without callbacks as you can see in examples/calculator, the
//     parser_*.go generate Ast nodes that are of types NativeArray, NativeInt,
//     NativeMap, NativeString.
//
//     To generate more specific Ast types, you may take a look at examples/lisp.
//
// Note: Parsers such as sequence, choice, not, pattern are also Ast,
// they are produced when parsing a valid joeson grammar; and they in turn help
// parsing that grammar.
type (
	Ast interface {
		String() string // text representation of this ast.
		GetLine() int
		SetLine(int)
		GetOrigin() Origin
		SetOrigin(o Origin)
		HasAttribute(key interface{}) bool
		GetAttribute(key interface{}) interface{}
		SetAttribute(key interface{}, value interface{})
	}
	Origin struct {
		Code     string
		Line     int
		Start    int
		End      int
		RuleName string
	}
)

func (o Origin) String() string {
	return fmt.Sprintf("Origin=(%d,%d,'rule=%s')", o.Start, o.End, o.RuleName)
}

var (
	_ Ast = &Grammar{}
	_ Ast = &choice{}
	_ Ast = &existential{}
	_ Ast = &lookahead{}
	_ Ast = &not{}
	_ Ast = &pattern{}
	_ Ast = &rank{}
	_ Ast = &regex{}
	_ Ast = &sequence{}
	_ Ast = &str{}
	_ Ast = &NativeArray{}
	_ Ast = NewNativeInt(-1)
	_ Ast = &NativeMap{}
	_ Ast = NewNativeString("")
	_ Ast = &NativeUndefined{}
	_ Ast = &ParseError{}
)

// prefix(x) + x.String(x)
func String(ast Ast) string {
	if x, isParser := ast.(Parser); isParser {
		return prefix(x) + x.String()
	} else {
		return ast.String()
	}
}
