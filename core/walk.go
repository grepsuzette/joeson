package core

type WalkPrepost struct {
	Pre  func(Astnode, parent Astnode) string // called during 🡖  .  "__stop__" to interrupt
	Post func(Astnode, parent Astnode) string // called during 🡕 .
}

// depth first walk
// `ast` is the node on which to start descending recursively.
// `parent` is available for algorithms needing it (just provide the
// father of `ast` or nil).
func Walk(ast Astnode, parent Astnode, prepost WalkPrepost) Astnode {
	if prepost.Pre != nil {
		// note: joeson.coffee can return "__stop__" here, meaning to end here (Not implemented yet)
		var stop = prepost.Pre(ast, parent) // don't implement coffee version "__stop__" just yet
		if stop == "__stop__" {
			return ast
		}
	}
	ast.ForEachChild(func(child Astnode) Astnode {
		return Walk(child, ast, prepost)
	})
	if prepost.Post != nil {
		prepost.Post(ast, parent)
	}
	return ast
}

// -- following are shortcut functions. They help prevent blunders --

// shortcut calling ForEachChild for members being []Astnode
func ForEachChild_Array(a []Astnode, f func(Astnode) Astnode) []Astnode {
	anew := []Astnode{}
	for _, child := range a {
		if r := f(child); r != nil {
			anew = append(anew, r)
		} // else, removed
	}
	return anew
}

// shortcut calling ForEachChild for members being map[string]Astnode
func ForEachChild_MapString(h map[string]Astnode, f func(Astnode) Astnode) map[string]Astnode {
	hnew := map[string]Astnode{}
	for k, child := range h {
		if r := f(child); r != nil {
			hnew[k] = r
		} // else, removed
	}
	return hnew
}