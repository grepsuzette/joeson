package joeson

/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/

type gnode interface {
	getgnode() *gnodeimpl
	Name() string
	Label() string
	Capture() bool

	SetName(string)
	SetNameWhenEmpty(string)
	SetLabel(string)
	SetCapture(bool)

	// The following funcs are cached and work with a helpers.Lazy that
	//  can be set (optionally) by each Parser type:

	Labels() []string
	Captures() []Ast

	SetLabels([]string)
	SetCaptures([]Ast)
	SetLazyLabels(func() []string)
	SetLazyCaptures(func() []Ast)
}
