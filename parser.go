package joeson

import (
	"github.com/grepsuzette/joeson/helpers"
)

// Parser objects are normally built by the joeson grammar
// and are able in turn to parse a ParseContext, producing Ast nodes.
//
// Parser are partially implemented by a common implementation called
// gnodeimpl.
type Parser interface {
	Ast
	Parse(ctx *ParseContext) Ast
	gnode() *gnodeimpl

	ForEachChild(func(Parser) Parser) Parser // depth-first walk mapper

	GetRuleName() string
	GetRuleLabel() string
	SetRuleLabel(string)
	Capture() bool
	SetRuleName(string)
	SetRuleNameWhenEmpty(string)
	SetCapture(bool)

	// The following functions are cached.
	// They work with an helpers.Lazy and can be optionally set by each Parser:
	Labels() []string
	Captures() []Ast
	SetLabels([]string)
	SetCaptures([]Ast)
	SetLazyLabels(func() []string)
	SetLazyCaptures(func() []Ast)

	prepare()
	handlesChildLabel() bool
}

var (
	_ Parser = &Grammar{}
	_ Parser = &choice{}
	_ Parser = &existential{}
	_ Parser = &lookahead{}
	_ Parser = &not{}
	_ Parser = &pattern{}
	_ Parser = &rank{}
	_ Parser = &regex{}
	_ Parser = &sequence{}
	_ Parser = &str{}
	_ Parser = &cLine{}
	_ Parser = &NativeUndefined{}
)

// a partial common implementation for Parser
type gnodeimpl struct {
	ParseOptions
	parent    Parser                  // A grammar must be a DAG (root.Parent being nil)
	name      string                  // rule name, if IsRule(). E.g. "AddOp" in `i(Named("AddOp", "'+' | '-'"))`
	label     string                  // rule label, e.g. "l" in `l:list` in `i(named("expr", "l:list | s:string"), parseExpr),`
	capture   bool                    // determines in which way to collect things higher up (see for instance Sequence.calculateType())
	labels_   *helpers.Lazy[[]string] // the lazy labels getter, redefinable to simulate GNode behavior in the original coffeescript impl. See NewGNode() doc below.
	captures_ *helpers.Lazy[[]Ast]    // the lazy captures getter, ditto.
	rules     map[string]Parser       // key is the rule name.
	rulesK    []string                // golang maps are unsorted, this keeps the insertion order
	id        int                     // rule number in a grammar. They start on 0. Use TRACE=grammar to list the rules and their ids. See also map grammar.id2Rule.
	rule      Parser                  // what's the Parser to use to parse this gnode
	grammar   *Grammar                // the grammar itself
	node      Parser                  // node containing this impl. Hack. Only used by GNode.Captures_ default implementation.
}

func newGNode() *gnodeimpl {
	gn := &gnodeimpl{
		capture: true,
		rules:   map[string]Parser{},
		rulesK:  []string{},
	}
	// labels_ and captures_ callbacks can be redefined by individual parsers
	// such as Sequence, Not etc.
	//
	// This helped regaining a certain level of flexibility that seemed
	// lost when going from the very dynamic combo javascript + clazz to
	// golang. Technically helpers.Lazy is just a callback whose execution
	// result is cached for ulterior calls.
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

// for now you must not include rules manually
// after the grammar was initialized
func (gn *gnodeimpl) Include(name string, rule Parser) {
	rule.SetRuleNameWhenEmpty(name)
	gn.rulesK = append(gn.rulesK, name)
	gn.rules[name] = rule
}

func (gn *gnodeimpl) GetRuleName() string     { return gn.name }
func (gn *gnodeimpl) SetRuleName(name string) { gn.name = name }
func (gn *gnodeimpl) SetRuleNameWhenEmpty(name string) {
	if gn.name == "" {
		gn.name = name
	}
}
func (gn *gnodeimpl) GetRuleLabel() string      { return gn.label }
func (gn *gnodeimpl) SetRuleLabel(label string) { gn.label = label }
func (gn *gnodeimpl) Capture() bool             { return gn.capture }
func (gn *gnodeimpl) SetCapture(b bool)         { gn.capture = b }

func (gn *gnodeimpl) Labels() []string                { return gn.labels_.Get() }
func (gn *gnodeimpl) Captures() []Ast                 { return gn.captures_.Get() }
func (gn *gnodeimpl) SetLabels(v []string)            { gn.labels_.Set(v) }
func (gn *gnodeimpl) SetCaptures(v []Ast)             { gn.captures_.Set(v) }
func (gn *gnodeimpl) SetLazyLabels(f func() []string) { gn.labels_ = helpers.NewLazyFromFunc(f) }
func (gn *gnodeimpl) SetLazyCaptures(f func() []Ast)  { gn.captures_ = helpers.NewLazyFromFunc(f) }

func IsRule(parser Parser) bool {
	return parser.gnode().rule == parser
}

// Return a prefix consisting of a name or a label when appropriate.
func prefix(parser Parser) string {
	if IsRule(parser) {
		return Red(parser.GetRuleName() + ": ")
	} else if parser.GetRuleLabel() != "" {
		return Cyan(parser.GetRuleLabel() + ":")
	} else {
		return ""
	}
}
