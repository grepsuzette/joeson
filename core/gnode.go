package core

// import . "grepsuzette/joeson/colors"

// Astnode should implement GetGNode(), fine, but
// Why exactly should GNode satisfy Astnode?
// I don't get it.
// TODO try and break this necessity

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
	Id      int                // joeson.coffee:604: `node.id = @numRules++`, in Grammar.
	Rule    Astnode            // set by Grammar.Postinit's walk into grammar nodes.
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
	*/
	CbBuilder func(Astnode, *ParseContext) Astnode
	SkipCache bool
	SkipLog   bool
	Debug     bool
}

func NewGNode() *GNode {
	return &GNode{Capture: true, Rules: map[string]Astnode{}}
}

func (gn *GNode) Labels() []string {
	if gn.Label != "" {
		return []string{gn.Label}
	} else {
		return []string{}
	}
}

func (gn *GNode) Captures() []Astnode {
	if gn.Capture {
		return []Astnode{gn}
	} else {
		return []Astnode{}
	}
}

func (gn *GNode) Prepare()                        {} // please put nothing in here
func (gn *GNode) GetGNode() *GNode                { return gn }
func (gn *GNode) Parse(ctx *ParseContext) Astnode { panic("GNode.Parse() must not be called") }
func (gn *GNode) ContentString() string           { return "<naked GNode, please redefine>" }
func (gn *GNode) HandlesChildLabel() bool         { return false }

// func (gn *GNode) ToString() string {
// 	return "donot use GNode.ToString(), use core.Show()"
// }

func (gn *GNode) Include(name string, rule Astnode) {
	if rule.GetGNode().Name == "" {
		rule.GetGNode().Name = name
	}
	gn.Rules[name] = rule
}

func (gn *GNode) IsRule() bool {
	return gn.Rule != nil && gn == gn.Rule.GetGNode()
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

// GNode by its own can't have hcildren.
// GNode prolly should not even satisfy Astnode in the 1st place
func (gn *GNode) ForEachChild(f func(Astnode) Astnode) Astnode {
	return gn
}
