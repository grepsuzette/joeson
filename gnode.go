package joeson

import "grepsuzette/joeson/helpers"

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
}

// TODO sink that lower in this file after it's done
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

// A grammar node.
type GNodeImpl struct {
	ParseOptions
	name      string // rule name if IsRule(), empty otherwise
	label     string
	capture   bool
	Labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript
	Captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, ditto.
	Rules     map[string]Parser       // Treelike. Grammar collects all rules in its post walk
	RulesK    []string                // because golang maps are unsorted, this helps keeping the insertion order
	Id        int                     // Numeric id of a Rule. Start on 0.
	Index     int                     // joeson.coffee:1303
	Rule      Parser                  // What's the rule for the node with this gnode. When Rule == node, it means node is a rule of a grammar (in which case node.IsRule() is true)
	Parent    Parser                  // Grammar must be a DAG, which implies 1 Parent at most (root.Parent being nil)
	Grammar   *Grammar
	node      Parser // The node containing this GNode. Only used by GNode.Captures_ default implementation.
	Origin    Origin // Where this node originates from.
}

func NewGNode() *GNodeImpl {
	gn := &GNodeImpl{
		capture: true,
		Rules:   map[string]Parser{},
	}
	// These callbacks can be redefined in Ast objects composing GNode.
	// This helps getting a certain level of flexibiliy.
	gn.Labels_ = helpers.NewLazyFromFunc(func() []string {
		if gn.label != "" {
			return []string{gn.label}
		} else {
			return []string{}
		}
	})
	gn.Captures_ = helpers.NewLazyFromFunc(func() []Ast {
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
	gn.RulesK = append(gn.RulesK, name)
	gn.Rules[name] = rule
}
