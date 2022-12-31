package core

import "grepsuzette/joeson/helpers"

/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/
type GNode struct {
	Node      Astnode                   // the node containing this GNode. This is only used by GNode._Captures default implementation
	Name      string                    // "" normally but rule name if IsRule()
	Label     string                    // "" if no label
	Capture   bool                      // usually true, false for instance for Str
	Labels_   *helpers.Lazy0[[]string]  // internal cache for Labels().
	Captures_ *helpers.Lazy0[[]Astnode] // internal cache for Captures().
	Rules     map[string]Astnode        // Treelike. Grammar collects all rules in its post walk
	RulesK    []string                  // because in golang maps are unsorted, the insertion order can be stored in this array
	Id        int                       // Numeric id of a Rule. It is incremented in Grammar. joeson.coffee:604: `node.id = @numRules++`
	Index     int                       // joeson.coffee:1303
	Rule      Astnode                   // An Astnode is a Rule when ast.GetGNode().Rule == ast. Set by Grammar.Postinit's walk into grammar nodes.
	Parent    Astnode                   // Grammar tree should be a DAG implying 1 Parent. Set by Grammar.Postinit's walk into grammar nodes.
	Grammar   Astnode                   // why Astnode and not *ast.Grammar, because core package must not depend on ast package. joeson.coffee:592, joeson.coffee:532.
	_origin   Origin                    // automatically set by prepareResult when a node is being parsed (prepareResult is called by wrap). Unused ATM

	/*
	 `cbBuilder` represents optional callbacks declared within inlined rules.
	 E.g. the func in `o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?",
	 		   func(result Astnode) Astnode { return ast.NewPattern(result) }),`

	 Since this example have labels, `result` will be of type NativeMap (which
	 implements Astnode) with the 3 keys "value", "join" and "@". Otherwise
	 it will be a NativeArray.

	 Second arg `...*ParseContext` is rarely passed in practice,
	 see a rare use in joescript.coffee:660.

	 Third arg `Astnode` is the caller Astnode (see joeson.js:455
	 or joeson.coffee:278) and represents the bounded `this` in javascript.
	*/
	CbBuilder func(nativeMapUsually Astnode, ctx *ParseContext, caller Astnode) Astnode
	SkipCache bool
	SkipLog   bool
	Debug     bool
}

func NewGNode() *GNode {
	gn := &GNode{
		Capture: true,
		Rules:   map[string]Astnode{},
	}
	gn.Labels_ = helpers.NewLazy0[[]string](func() []string {
		if gn.Label != "" {
			return []string{gn.Label}
		} else {
			return []string{}
		}
	})
	gn.Captures_ = helpers.NewLazy0[[]Astnode](func() []Astnode {
		if gn.Capture {
			return []Astnode{gn.Node}
		} else {
			return []Astnode{}
		}
	})
	// gn.Labels_ = NewLazy0[[]string](func() []string {
	// 	if gn.Label != "" {
	// 		return []string{gn.Label}
	// 	} else {
	// 		return []string{}
	// 	}
	// })
	// gn.Captures_ = NewLazy1[[]Astnode, Astnode](func(x Astnode) []Astnode {
	// 	if gn.Capture {
	// 		return []Astnode{x}
	// 	} else {
	// 		return []Astnode{}
	// 	}
	// })
	return gn
}

func (gn *GNode) Include(name string, rule Astnode) {
	if rule.GetGNode().Name == "" {
		rule.GetGNode().Name = name
	}
	gn.RulesK = append(gn.RulesK, name)
	gn.Rules[name] = rule
}

// find a parent in the ancestry chain that satisfies condition
func (gn *GNode) FindParentHaving(fcond func(Astnode) bool) Astnode {
	var x = gn.Parent
	for {
		if fcond(x) {
			return x
		} else {
			x = x.GetGNode().Parent
		}
	}
}
