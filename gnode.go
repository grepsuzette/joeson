package joeson

// gnode helps writing parsers.
type gnode interface {
	gnode() *gnodeimpl
	GetRuleName() string
	GetRuleLabel() string
	SetRuleLabel(string)
	Capture() bool

	SetRuleName(string)
	SetRuleNameWhenEmpty(string)
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

var (
	_ gnode = &Grammar{}
	_ gnode = &choice{}
	_ gnode = &existential{}
	_ gnode = &lookahead{}
	_ gnode = &not{}
	_ gnode = &pattern{}
	_ gnode = &rank{}
	_ gnode = &ref{}
	_ gnode = &regex{}
	_ gnode = &sequence{}
	_ gnode = &str{}
	_ gnode = &NativeUndefined{}
)
