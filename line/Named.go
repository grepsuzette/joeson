package line

import (
	"fmt"
	"grepsuzette/joeson/core"
)

// NamedRule satisfies Astnode. Otherwise, it just adds a Name to a Line
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
func Named(name string, lineStringOrAstnode any) NamedRule {
	switch v := lineStringOrAstnode.(type) {
	case Line:
		return NamedRule{name, v}
	case string:
		return NamedRule{name, NewSLine(v)}
	case core.Ast:
		return NamedRule{name, NewCLine(v)}
	default:
		msg := fmt.Sprintf("Invalid argument to Named(%s, %v)\n", name, lineStringOrAstnode)
		panic(msg)
	}
}

func (nm NamedRule) Parse(ctx *core.ParseContext) core.Ast           { panic("precompiled") }
func (nm NamedRule) ContentString() string                           { return "--Named--" }
func (nm NamedRule) GetGNode() *core.GNode                           { panic("assert") }
func (nm NamedRule) Prepare()                                        {}
func (nm NamedRule) HandlesChildLabel() bool                         { return false }
func (nm NamedRule) ForEachChild(f func(core.Ast) core.Ast) core.Ast { return nm }
