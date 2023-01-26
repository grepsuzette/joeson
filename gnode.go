package joeson

import "grepsuzette/joeson/helpers"

/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/
type GNode struct {
	Name string // rule name if IsRule(), empty otherwise
	ParseOptions
	Label     string
	Capture   bool                    // usually true, false for instance for Str
	Labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript
	Captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, ditto.
	Rules     map[string]Ast          // Treelike. Grammar collects all rules in its post walk
	RulesK    []string                // because golang maps are unsorted, this helps keeping the insertion order
	Id        int                     // Numeric id of a Rule. Start on 0.
	Index     int                     // joeson.coffee:1303
	Rule      Ast                     // What's the rule for the node with this gnode. When Rule == node, it means node is a rule of a grammar (in which case node.IsRule() is true)
	Parent    Ast                     // Grammar must be a DAG, which implies 1 Parent at most (root.Parent being nil)
	Grammar   Ast                     // More precisely an *ast.Grammar, but core can not depend upon the ast package. joeson.coffee:592, joeson.coffee:532.
	Node      Ast                     // The node containing this GNode. Only used by GNode.Captures_ default implementation.
	Origin    Origin                  // Where this node originates from.
}

func NewGNode() *GNode {
	gn := &GNode{
		Capture: true,
		Rules:   map[string]Ast{},
	}
	// These callbacks can be redefined in Ast objects composing GNode.
	// This helps getting a certain level of flexibiliy.
	gn.Labels_ = helpers.NewLazyFromFunc[[]string](func() []string {
		if gn.Label != "" {
			return []string{gn.Label}
		} else {
			return []string{}
		}
	})
	gn.Captures_ = helpers.NewLazyFromFunc[[]Ast](func() []Ast {
		if gn.Capture {
			return []Ast{gn.Node}
		} else {
			return []Ast{}
		}
	})
	return gn
}

func (gn *GNode) Include(name string, rule Ast) {
	if rule.GetGNode().Name == "" {
		rule.GetGNode().Name = name
	}
	gn.RulesK = append(gn.RulesK, name)
	gn.Rules[name] = rule
}
