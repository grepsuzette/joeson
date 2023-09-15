package joeson

import (
	"github.com/grepsuzette/joeson/helpers"
)

// rule partially implements Parser,
// it is used by compositition by all parser_*.go
type rule struct {
	*parseOptions                   // whenever options are exceptionally provided inline for a rule. E.g. i(named("FOO", `'foo'`), Debug{true})
	parent        Parser            // A grammar must be a DAG (root.parent being nil)
	parser        Parser            // what's the Parser to use to parse this rule
	grammar       *Grammar          // the grammar itself
	id            int               // rule number in a grammar. They start at 0. Use TRACE=grammar to list the rules and their ids. See also map grammar.id2Rule.
	name          string            // rule name, e.g. "AddOp" in `i(Named("AddOp", `'+' | '-'`))`
	label         string            // rule label, e.g. "t" in `i(named("expr", `t:list | s:string`), parseExpr),`
	capture       bool              // see Sequence.calculateType()
	rules         map[string]Parser // key is the rule name.
	rulesK        []string          // golang maps are unsorted, this keeps the insertion order
	node          Parser            // node containing this impl. Hack. Only used by GNode.Captures_ default implementation.

	// TODO reevaluate if we need it
	labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript impl.
	captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, same concept. (not used at the moment)
}

func newRule() *rule {
	r := &rule{
		parseOptions: newParseOptions(),
		capture:      true,
		rules:        map[string]Parser{},
		rulesK:       []string{},
	}

	// labels and captures are lazy. These are objects that can contain
	// a value or a callback whose result will be cached for later.
	//
	// They can be redefined by individual parsers
	// such as Sequence, Not etc.
	//
	// This helps regaining a certain level of flexibility for edge cases
	// when going from the very dynamic javascript + clazz combination to golang.
	r.labels_ = helpers.LazyFromFunc(func() []string {
		if r.label != "" {
			return []string{r.label}
		} else {
			return []string{}
		}
	})
	r.captures_ = helpers.LazyFromFunc(func() []Ast {
		if r.capture {
			return []Ast{r.node}
		} else {
			return []Ast{}
		}
	})
	return r
}

// for now you must not include rules manually
// after the grammar was initialized.
func (r *rule) Include(name string, parser Parser) {
	if r.name == "" {
		r.name = name
	}
	r.rulesK = append(r.rulesK, name)
	r.rules[name] = parser
}
