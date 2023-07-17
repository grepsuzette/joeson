package joeson

// gnode helps writing parsers.
type gnode interface {
	gnode() *gnodeimpl
	Name() string
	Label() string
	Capture() bool

	SetName(string)
	SetNameWhenEmpty(string)
	SetLabel(string)
	SetCapture(bool)

	// The following functions are cached.
	// They work with an helpers.Lazy and can be optionally set by each Parser:

	Labels() []string
	Captures() []Ast

	SetLabels([]string)
	SetCaptures([]Ast)
	SetLazyLabels(func() []string)
	SetLazyCaptures(func() []Ast)
}
