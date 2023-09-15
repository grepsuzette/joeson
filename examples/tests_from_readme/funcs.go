package readmetest

import j "github.com/grepsuzette/joeson"

// named() creates a rule with a name
func named(name string, v any) j.NamedRule { return j.Named(name, v) }

// to make i and o rules
func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }
