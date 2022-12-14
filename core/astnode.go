package core

import . "grepsuzette/joeson/colors"

// Astnode are parsed node (with an Origin in the original CodeStream)
// GNode are the grammar node rule. An Astnode always originates from
// an Origin and a GNode. Some Astnode (Native{Map,Array,String,Undefined})
// come without a grammar, however they still have an origin and therefore
// an almost empty GNode.
type Astnode interface {
	// Parse() reads from ParseContext, updating context's position.
	// Return nil to revert ParseContext.pos whenever parsing failed.
	Parse(ctx *ParseContext) Astnode
	ContentString() string // colorful representation of an AST node
	GetGNode() *GNode      // nodes without a grammar node (Native*) return nil
	Prepare()              // called after children prepared
	HandlesChildLabel() bool
	Labels() []string
	Captures() []Astnode
	ForEachChild(f func(Astnode) Astnode) Astnode
}

func IsRule(x Astnode) bool {
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
