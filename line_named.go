package joeson

import (
	"fmt"
)

// Named helps simplifying writing trees of rules.
//
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
		return NamedRule{name, newSLine(v)}
	case Parser:
		return NamedRule{name, newCLine(v)}
	default:
		msg := fmt.Sprintf("Invalid argument to Named(%s, %v)\n", name, lineStringOrAst)
		panic(msg)
	}
}

// NamedRule is produced by Named().
// Note NamedRule technically satisfies Ast.
// There is no much reason to keep it public.
type NamedRule struct {
	name string
	line Line // O, I or A
}
