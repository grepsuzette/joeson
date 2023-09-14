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
	gnode() *rule

	ForEachChild(func(Parser) Parser) Parser // depth-first walk mapper

	GetRuleName() string
	GetRuleLabel() string
	SetRuleLabel(string)
	SetRuleName(string)
	SetRuleNameWhenEmpty(string)

	// HandlesChildLabel() is hardly used at all right now
	// TODO must be settable
	// It prevents collecting labelled children as part
	// of a NativeMap in packrat.go prepareResult()
	HandlesChildLabel() bool

	Capture() bool
	SetCapture(bool)

	// The following functions are cached.
	// They work with an helpers.Lazy and can be optionally set by each Parser:
	Labels() []string   // get label (immediate, cached, or evaluate the lazy)
	Captures() []Ast    // if current gnode capture is true, it will of course be part of Captures()
	SetLabels([]string) // set an immediate value
	SetCaptures([]Ast)
	SetLabelsLazy(func() []string) // set a lazy value evaluated the 1st time it's needed
	SetCapturesLazy(func() []Ast)

	prepare()
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
