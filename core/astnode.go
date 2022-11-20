package core

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
