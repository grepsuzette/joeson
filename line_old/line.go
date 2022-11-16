package line

// Line interface is just a way to have []Line really
// These are a system to enter rules of a grammar
// in a code-like fashion (as a tree, rather than linearly).
type Line interface {
	String() string
	StringIndent(nIndent int) string // same as String(), but indenting with `nIdent` levels (for nested rules)
	LineType() string
	IsO() bool
}
