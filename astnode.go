package main

type astnode interface {
	/* Parse() reads from ParseContext and
	 * also can modify it (move the position,
	 * by ctx.code.pos = xx).
	 *
	 * Result is... TODO document it
	 */
	Parse(*ParseContext) astnode

	/* To access the sub structure */
	GetGNode() *GNode

	/* CoLoRfUl representation of an AST node */
	ContentString() string

	/* Called after all its children have been prepared. */
	Prepare()

	/* Almost always returns false except for Sequence and Existential */
	HandlesChildLabel() bool

	// In current joeson impl,
	// all Labels() return GetGNode().Labels()
	// except Existential, GNode, Ref and Sequence.
	Labels() []string

	// In current joeson impl,
	// all Captures() return GetGNode().Captures()
	// except Existential, GNode and Sequence.
	Captures() []astnode
}
