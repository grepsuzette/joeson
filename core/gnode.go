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
	Name      string                  // "" normally but rule name if IsRule()
	Label     string                  // "" if no label
	Capture   bool                    // usually true, false for instance for Str
	Labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate original GNode behavior in coffee
	Captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, idem.
	Rules     map[string]Ast          // Treelike. Grammar collects all rules in its post walk
	RulesK    []string                // because in golang maps are unsorted, the insertion order can be stored in this array
	Id        int                     // Numeric id of a Rule. It is incremented in Grammar. joeson.coffee:604: `node.id = @numRules++`
	Index     int                     // joeson.coffee:1303
	Rule      Ast                     // An Ast is a Rule when ast.GetGNode().Rule == ast. Set by Grammar.Postinit's walk into grammar nodes.
	Parent    Ast                     // Grammar tree should be a DAG implying 1 Parent. Set by Grammar.Postinit's walk into grammar nodes.
	Grammar   Ast                     // why Ast and not *ast.Grammar, because core package must not depend on ast package. joeson.coffee:592, joeson.coffee:532.
	Node      Ast                     // Only used by GNode.Captures_ default implementation. The node containing this GNode.
	_origin   Origin                  // automatically set by prepareResult when a node is being parsed (prepareResult is called by wrap). Unused ATM

	/*
	 `cbBuilder` represents optional callbacks declared within inlined rules.
	 E.g. the func in `o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?",
	 		   func(result Ast) Ast { return ast.NewPattern(result) }),`

	 Since this example have labels, `result` will be of type NativeMap (which
	 implements Ast) with the 3 keys "value", "join" and "@". Otherwise
	 it will be a NativeArray.

	 Second arg `...*ParseContext` is rarely passed in practice,
	 see a rare use in joescript.coffee:660.

	 Third arg `Ast` is the caller Ast (see joeson.js:455
	 or joeson.coffee:278) and represents the bounded `this` in javascript.
	*/
	CbBuilder func(nativeMapUsually Ast, ctx *ParseContext, caller Ast) Ast
	SkipCache bool
	SkipLog   bool
	Debug     bool
}

func NewGNode() *GNode {
	gn := &GNode{
		Capture: true,
		Rules:   map[string]Ast{},
	}
	gn.Labels_ = helpers.NewLazy[[]string](func() []string {
		if gn.Label != "" {
			return []string{gn.Label}
		} else {
			return []string{}
		}
	})
	gn.Captures_ = helpers.NewLazy[[]Ast](func() []Ast {
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

// find a parent in the ancestry chain that satisfies condition
func (gn *GNode) FindParentHaving(fcond func(Ast) bool) Ast {
	var x = gn.Parent
	for {
		if fcond(x) {
			return x
		} else {
			x = x.GetGNode().Parent
		}
	}
}
