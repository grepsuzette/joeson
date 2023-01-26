package joeson

import (
	"fmt"
)

// NamedRule satisfies Ast. Otherwise, it just adds a Name to a Line
type NamedRule struct {
	Name string
	Line Line // note can be OLine, ILine or ALine (array)
}

// Examples:
// o(Named("LABELED", o(<compiled>), o(Named("DECORATED", o(<compiled>, ..))))
// o(Named("DECORATED", o(<compiled>), o(<compiled>), i(<compiled>)))
// i(Named("RANGE", o(<compiled>))
// i(Named("LABEL", C(St("&"), St("@"), R("WORD")))),
// i(Named("LABEL", "'&' | '@' | WORD")),
//
// For i, it is necessary for len(lines) == 1
//   this is to be enforcedby the caller.
func Named(name string, lineStringOrAst any) NamedRule {
	switch v := lineStringOrAst.(type) {
	case Line:
		return NamedRule{name, v}
	case string:
		return NamedRule{name, NewSLine(v)}
	case Ast:
		return NamedRule{name, NewCLine(v)}
	default:
		msg := fmt.Sprintf("Invalid argument to Named(%s, %v)\n", name, lineStringOrAst)
		panic(msg)
	}
}

func (nm NamedRule) Parse(ctx *ParseContext) Ast      { panic("precompiled") }
func (nm NamedRule) ContentString() string            { return "--Named--" }
func (nm NamedRule) GetGNode() *GNode                 { panic("assert") }
func (nm NamedRule) Prepare()                         {}
func (nm NamedRule) HandlesChildLabel() bool          { return false }
func (nm NamedRule) ForEachChild(f func(Ast) Ast) Ast { return nm }
