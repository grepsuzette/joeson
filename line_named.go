package joeson

import (
	"fmt"
)

// A Key-value pair, where Key is the name.
// This is exclusively used with joeson ILine and OLine to name things.
func Named(name string, lineStringOrAst any) NamedRule {
	switch v := lineStringOrAst.(type) {
	case Line:
		return NamedRule{name, v}
	case string:
		return NamedRule{name, newSLine(v)}
	case Parser:
		return NamedRule{name, newCLine(v)}
	case []Line:
		return NamedRule{name, NewALine(v)}
	default:
		msg := fmt.Sprintf("Invalid argument to Named(%s, %v)\n", name, lineStringOrAst)
		panic(msg)
	}
}

// A Key-value pair, where Key is the name.
// Use Named() instead of building directly.
type NamedRule struct {
	name string
	line Line // O, I or A
}
