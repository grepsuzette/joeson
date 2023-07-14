package joeson

import "github.com/grepsuzette/joeson/helpers"

// gnodeimpl holds the data for a grammar node
type gnodeimpl struct {
	ParseOptions
	parent    Parser                  // A grammar must be a DAG (root.Parent being nil)
	name      string                  // rule name if IsRule(). E.g. "AddOp" in `i(Named("AddOp", "'+' | '-'"))`
	label     string                  // e.g. "l" in `l:list` in `i(named("expr", "l:list | s:string"), parseExpr),`
	capture   bool                    // determines in which way to collect things higher up (see for instance Sequence.calculateType())
	labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript impl. See NewGNode() doc below.
	captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, ditto.
	rules     map[string]Parser       // key is the rule name.
	rulesK    []string                // golang maps are unsorted, this keeps the insertion order
	id        int                     // rule number in a grammar. They start on 0. Use TRACE=grammar to list the rules and their ids. See also map grammar.id2Rule.
	rule      Parser                  // what's the Parser to use to parse this gnode
	grammar   *Grammar                // the grammar itself
	node      Parser                  // node containing this impl. Hack. Only used by GNode.Captures_ default implementation.
	origin    Origin                  // records where this gnode originates from. Unused for now.
}

func NewGNode() *gnodeimpl {
	gn := &gnodeimpl{
		capture: true,
		rules:   map[string]Parser{},
		rulesK:  []string{},
	}
	// labels_ and captures_ callbacks can be redefined by individual parsers
	// such as Sequence, Not etc.
	//
	// This helped regaining a certain level of flexibility that seemed
	// lost when going from the very dynamic combo javascript + clazz to
	// golang. Technically helpers.Lazy is just a callback whose execution
	// result is cached for ulterior calls.
	gn.labels_ = helpers.NewLazyFromFunc(func() []string {
		if gn.label != "" {
			return []string{gn.label}
		} else {
			return []string{}
		}
	})
	gn.captures_ = helpers.NewLazyFromFunc(func() []Ast {
		if gn.capture {
			return []Ast{gn.node}
		} else {
			return []Ast{}
		}
	})
	return gn
}

func (gn *gnodeimpl) Include(name string, rule Parser) {
	rule.SetNameWhenEmpty(name)
	gn.rulesK = append(gn.rulesK, name)
	gn.rules[name] = rule
}

func (gn *gnodeimpl) GetGNode() *gnodeimpl { return gn }
func (gn *gnodeimpl) Name() string         { return gn.name }
func (gn *gnodeimpl) Label() string        { return gn.label }
func (gn *gnodeimpl) Capture() bool        { return gn.capture }
func (gn *gnodeimpl) SetName(name string)  { gn.name = name }
func (gn *gnodeimpl) SetNameWhenEmpty(name string) {
	if gn.name == "" {
		gn.name = name
	}
}
func (gn *gnodeimpl) SetLabel(label string) { gn.label = label }
func (gn *gnodeimpl) SetCapture(b bool)     { gn.capture = b }

func (gn *gnodeimpl) Labels() []string                { return gn.labels_.Get() }
func (gn *gnodeimpl) Captures() []Ast                 { return gn.captures_.Get() }
func (gn *gnodeimpl) SetLabels(v []string)            { gn.labels_.Set(v) }
func (gn *gnodeimpl) SetCaptures(v []Ast)             { gn.captures_.Set(v) }
func (gn *gnodeimpl) SetLazyLabels(f func() []string) { gn.labels_ = helpers.NewLazyFromFunc(f) }
func (gn *gnodeimpl) SetLazyCaptures(f func() []Ast)  { gn.captures_ = helpers.NewLazyFromFunc(f) }
