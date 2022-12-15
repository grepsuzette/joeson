package core

/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/
type GNode struct {
	Name    string             // "" normally but rule name if IsRule()
	Label   string             // "" if no label
	Rules   map[string]Astnode // Treelike. Grammar collects all rules in its post walk
	Id      int                // Numeric id of a Rule. It is incremented in Grammar. joeson.coffee:604: `node.id = @numRules++`
	Index   int                // joeson.coffee:1303
	Rule    Astnode            // An Astnode is a Rule when ast.GetGNode().Rule == ast. Set by Grammar.Postinit's walk into grammar nodes.
	Parent  Astnode            // Grammar tree should be a DAG implying 1 Parent. Set by Grammar.Postinit's walk into grammar nodes.
	Grammar Astnode            // joeson.coffee:592, joeson.coffee:532.
	Capture bool               // usually true, false for instance for Str
	_origin Origin             // automatically set by prepareResult when a node is being parsed (prepareResult is called by wrap). Unused ATM

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
	return &GNode{Capture: true, Rules: map[string]Astnode{}}
}

func MeIfCaptureOrEmpty(x Astnode) []Astnode {
	if x.GetGNode().Capture {
		return []Astnode{x}
	} else {
		return []Astnode{}
	}
}

func MyLabelIfDefinedOrEmpty(x Astnode) []string {
	if x.GetGNode().Label != "" {
		return []string{x.GetGNode().Label}
	} else {
		return []string{}
	}
}

func (gn *GNode) Include(name string, rule Astnode) {
	if rule.GetGNode().Name == "" {
		rule.GetGNode().Name = name
	}
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
