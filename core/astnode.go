package core

import . "grepsuzette/joeson/colors"

// Astnode are parsed node (with an Origin in the original CodeStream)
// GNode are the grammar node rule. An Astnode always originates from
// an Origin and a GNode. Some Astnode (Native{Map,Array,String,Undefined})
// come without a grammar, however they still have an origin and therefore
// an almost empty GNode.
type Astnode interface {

	/* Parse() reads from ParseContext and
	 * also can modify it (move the position,
	 * by ctx.Code.pos = xx).
	 *
	 * Result is... TODO document it
	 */
	// TODO in Existential joeson.coffee has `return result != null ? result : void 0;` (void 0 being translation of `undefined`)
	// is there an important difference when returning nil and undefined?
	// If yes .. it's a big pb to address now
	// Parse() can return nil to revert ParseContext.pos (joeson.coffee:491)
	Parse(ctx *ParseContext) Astnode

	ContentString() string // CoLoRfUl representation of an AST node
	GetGNode() *GNode      // Nodes without a grammar node return nil here (this is the case for Native*.go types)

	Prepare()                // Called after all its children have been prepared. Empty impl. is ok, GNode.Prepare() itself is empty.
	HandlesChildLabel() bool // In current joeson impl, always returns false except for Sequence and Existential
	Labels() []string        // In current joeson impl, all Labels() return GetGNode().Labels() except Existential, GNode, Ref and Sequence.
	Captures() []Astnode     // In current joeson impl, all Captures() return GetGNode().Captures() except Existential, GNode and Sequence.

	ForEachChild(f func(Astnode) Astnode) Astnode
}

func IsRule(x Astnode) bool {
	// more contorsion needed than in orig. impl. because might be a CLine
	// return x.GetGNode().Rule.GetGNode(). != nil && x.GetGNode().Rule.Id == x.GetGNode().Id
	return x.GetGNode().Rule == x
}

// this is the port of GNode.toString().
// It calls x.ContentString() but adds a prefix with the label or name.
func ContentStringWithPrefix(x Astnode) string {
	return Prefix(x) + x.ContentString()
}

func Prefix(x Astnode) string {
	if IsRule(x) {
		return Red(x.GetGNode().Name + ": ")
	} else if x.GetGNode().Label != "" {
		return Cyan(x.GetGNode().Label + ":")
	} else {
		return ""
	}

}
