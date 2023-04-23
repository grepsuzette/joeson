package joeson

import "github.com/grepsuzette/joeson/helpers"

/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/

type GNode interface {
	GetGNode() *GNodeImpl
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

// Data for a grammar node.
// GNode is the interface, but all implementers hold this impl.
type GNodeImpl struct {
	ParseOptions
	parent    Parser                  // A grammar must be a DAG, which implies 1 Parent at most (root.Parent being nil)
	name      string                  // rule name if IsRule(), empty otherwise. E.g. "AddOp" in `i(named("AddOp", "'+' | '-'"))`
	label     string                  // example "l" when this is `l:list` in `i(named("expr", "l:list | s:string | n:number | operator:operator"), parseExpr),`
	capture   bool                    // helps determining in which way to collect things higher up (see for instance Sequence.calculateType())
	labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript impl.
	captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, ditto.
	rules     map[string]Parser       // Treelike. Grammar collects all rules in its post walk
	rulesK    []string                // because golang maps are unsorted, this helps keeping the insertion order TODO check if really necessary
	id        int                     // rule number in a grammar. They start on 0. See also map grammar.id2Rule.
	rule      Parser                  // what's the rule for the node with this gnode. When Rule == node, it means node is a rule of a grammar (in which case node.IsRule() is true)
	grammar   *Grammar                // the root
	node      Parser                  // node containing this impl. Hackish. Only used by GNode.Captures_ default implementation.
	origin    Origin                  // records the start-end in the source of where this gnode originates from. Unused as of now.
}

func NewGNode() *GNodeImpl {
	gn := &GNodeImpl{
		capture: true,
		rules:   map[string]Parser{},
	}
	// These callbacks can be redefined in Ast objects composing GNode.
	// This helps getting a certain level of flexibiliy.
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

func (gn *GNodeImpl) Include(name string, rule Parser) {
	rule.SetNameWhenEmpty(name)
	gn.rulesK = append(gn.rulesK, name)
	gn.rules[name] = rule
}

func (gn *GNodeImpl) GetGNode() *GNodeImpl { return gn }
func (gn *GNodeImpl) Name() string         { return gn.name }
func (gn *GNodeImpl) Label() string        { return gn.label }
func (gn *GNodeImpl) Capture() bool        { return gn.capture }
func (gn *GNodeImpl) SetName(name string)  { gn.name = name }
func (gn *GNodeImpl) SetNameWhenEmpty(name string) {
	if gn.name == "" {
		gn.name = name
	}
}
func (gn *GNodeImpl) SetLabel(label string) { gn.label = label }
func (gn *GNodeImpl) SetCapture(b bool)     { gn.capture = b }

func (gn *GNodeImpl) Labels() []string                { return gn.labels_.Get() }
func (gn *GNodeImpl) Captures() []Ast                 { return gn.captures_.Get() }
func (gn *GNodeImpl) SetLabels(v []string)            { gn.labels_.Set(v) }
func (gn *GNodeImpl) SetCaptures(v []Ast)             { gn.captures_.Set(v) }
func (gn *GNodeImpl) SetLazyLabels(f func() []string) { gn.labels_ = helpers.NewLazyFromFunc(f) }
func (gn *GNodeImpl) SetLazyCaptures(f func() []Ast)  { gn.captures_ = helpers.NewLazyFromFunc(f) }

