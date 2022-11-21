package core

/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/
type GNode struct {
	// Node                              // Node solely containing .id int at this point, and so not being really necessary
	Name    string             // "" normally but rule name if IsRule()
	Label   string             // "" if no label <- i have a doubt that maybe should be *string.  Because of joeson.go:832...
	Rules   map[string]Astnode // all levels have rules (like a tree). Grammar will collect all rules in its post walk, increasing the NumRules each time, and affecting gnode Id = current NumRules
	Id      int                // Numeric id of a Rule. It is incremented in Grammar. joeson.coffee:604: `node.id = @numRules++`
	Index   int                // joeson.coffee:1303
	Rule    Astnode            // set by Grammar.Postinit's walk into grammar nodes. An Astnode is a Rule when ast.GetGNode().Rule == ast
	Parent  Astnode            // set by Grammar.Postinit's walk into grammar nodes. Grammar tree should be a DAG implies 1 Parent.
	Grammar Astnode            // set by Grammar.Postinit's walk into grammar nodes. joeson.coffee:592, joeson.coffee:532. Type is Grammar.
	Capture bool               // true by default, it's false for instance for Str
	_origin Origin             // automatically set by prepareResult when a node is being parsed (prepareResult is called by wrap)

	/*
	 `cbBuilder` represents optional callbacks declared within inlined rules.
	 E.g. the func in `o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?",
	 		   func(result Astnode) Astnode { return ast.NewPattern(result) }),`

	 Since this example have labels, `result` will be of type NativeMap (which
	 implements Astnode) with the 3 keys "value", "join" and "@". Otherwise
	 it will be a TODO TBD.

	 Second arg `...*ParseContext` is rarely passed in practice,
	 see a rare use in joescript.coffee:660.

	 Third arg `Astnode` is the caller Astnode (see joeson.js:455
	 or joeson.coffee:278) and represents the bounded `this` in javascript.
	*/
	CbBuilder func(nativeMapUsually Astnode, ctx *ParseContext, caller Astnode) Astnode
	Parse     func(ctx ParseContext) Astnode
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
