package joeson

// rule partially implements Parser,
// it is used by compositition by all parser_*.go
type rule struct {
	parseOptions
	parent  Parser            // A grammar must be a DAG (root.parent being nil)
	parser  Parser            // what's the Parser to use to parse this rule
	name    string            // rule name, if IsRule(). E.g. "AddOp" in `i(Named("AddOp", `'+' | '-'`))`
	label   string            // rule label, e.g. "t" in `i(named("expr", `t:list | s:string`), parseExpr),`
	capture bool              // see Sequence.calculateType()
	rules   map[string]Parser // key is the rule name.
	rulesK  []string          // golang maps are unsorted, this keeps the insertion order
	id      int               // rule number in a grammar. They start on 0. Use TRACE=grammar to list the rules and their ids. See also map grammar.id2Rule.
	grammar *Grammar          // the grammar itself
	node    Parser            // node containing this impl. Hack. Only used by GNode.Captures_ default implementation.

	labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript impl. See NewGNode() doc below.
	captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, ditto.
}

func newRule() *rule {
	gn := &rule{
		capture: true,
		rules:   map[string]Parser{},
		rulesK:  []string{},
	}

	// labels and captures are lazy. These are objects that can contain
	// a value or a callback whose result will be cached for later.
	//
	// They can be redefined by individual parsers
	// such as Sequence, Not etc.
	//
	// This helps regaining a certain level of flexibility for edge cases
	// when going from the very dynamic javascript + clazz combination to golang.
	gn.labels_ = helpers.LazyFromFunc(func() []string {
		if gn.label != "" {
			return []string{gn.label}
		} else {
			return []string{}
		}
	})
	gn.captures_ = helpers.LazyFromFunc(func() []Ast {
		if gn.capture {
			return []Ast{gn.node}
		} else {
			return []Ast{}
		}
	})
	return gn
}

// for now you must not include rules manually
// after the grammar was initialized
func (gn *rule) Include(name string, parser Parser) {
	parser.SetRuleNameWhenEmpty(name)
	gn.rulesK = append(gn.rulesK, name)
	gn.rules[name] = parser
}

func (gn *rule) GetRuleName() string     { return gn.name }
func (gn *rule) SetRuleName(name string) { gn.name = name }
func (gn *rule) SetRuleNameWhenEmpty(name string) {
	if gn.name == "" {
		gn.name = name
	}
}
func (gn *rule) GetRuleLabel() string      { return gn.label }
func (gn *rule) SetRuleLabel(label string) { gn.label = label }
func (gn *rule) Capture() bool             { return gn.capture }
func (gn *rule) SetCapture(b bool)         { gn.capture = b }

func (gn *rule) Labels() []string                { return gn.labels_.Get() }
func (gn *rule) Captures() []Ast                 { return gn.captures_.Get() }
func (gn *rule) SetLabels(v []string)            { gn.labels_.Set(v) }
func (gn *rule) SetCaptures(v []Ast)             { gn.captures_.Set(v) }
func (gn *rule) SetLabelsLazy(f func() []string) { gn.labels_ = helpers.LazyFromFunc(f) }
func (gn *rule) SetCapturesLazy(f func() []Ast)  { gn.captures_ = helpers.LazyFromFunc(f) }

func IsRule(parser Parser) bool {
	return parser.gnode().parser == parser
}
