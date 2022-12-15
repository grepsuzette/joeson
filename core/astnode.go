package core

import . "grepsuzette/joeson/colors"

type Astnode interface {
	// Parse() reads from ParseContext, updates context's position,
	// returns nil to indicate parse failure.
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

func Prefix(x Astnode) string {
	if IsRule(x) {
		return Red(x.GetGNode().Name + ": ")
	} else if x.GetGNode().Label != "" {
		return Cyan(x.GetGNode().Label + ":")
	} else {
		return ""
	}
}

func ContentStringWithPrefix(x Astnode) string {
	return Prefix(x) + x.ContentString()
}

func LabelOrName(n Astnode) string {
	if IsRule(n) {
		return Red(n.GetGNode().Name + ": ")
	} else if n.GetGNode().Label != "" {
		return Cyan(n.GetGNode().Label + ":")
	}
	return ""
}
