package core

// Depth-first walk

type WalkPrepost struct {
	Pre  func(node Ast, parent Ast) string // called during 🡖  .  "__stop__" to interrupt
	Post func(node Ast, parent Ast) string // called during 🡕 .
}

// Depth first walk of entire tree.
// `ast` is the node on which to start descending recursively.
// `parent` is available for algorithms needing it (just provide the
// father of `ast` or nil).
func Walk(ast Ast, parent Ast, prepost WalkPrepost) Ast {
	// TODO there can be a big difference. checki
	if prepost.Pre != nil {
		// note: joeson.coffee can return "__stop__" here, meaning to end here (Not implemented yet)
		var stop = prepost.Pre(ast, parent) // don't implement coffee version "__stop__" just yet
		if stop == "__stop__" {
			return ast
		}
	}
	ast.ForEachChild(func(child Ast) Ast {
		return Walk(child, ast, prepost)
	})
	if prepost.Post != nil {
		prepost.Post(ast, parent)
	}
	return ast
}

// -- following are shortcut functions.

// shortcut calling ForEachChild for members being []Ast
func ForEachChild_Array(a []Ast, f func(Ast) Ast) []Ast {
	anew := []Ast{}
	for _, child := range a {
		if r := f(child); r != nil {
			anew = append(anew, r)
		} // else, removed
	}
	return anew
}

// shortcut calling ForEachChild for members being map[string]Ast
// beware, maps are not ordered in golang. Use instead ForEachChild_InRules
// func ForEachChild_MapString(h map[string]Ast, f func(Ast) Ast) map[string]Ast {
// 	hnew := map[string]Ast{}
// 	sortedKeys := []string{}
// 	for k := range h {
// 		sortedKeys = append(sortedKeys, k)
// 	}
// 	sort.Strings(sortedKeys)
// 	for _, k := range sortedKeys {
// 		if r := f(h[k]); r != nil {
// 			hnew[k] = r
// 		} // else, removed
// 	}
// 	return hnew
// }

// where x.GetGNode().Rules and x.GetGNode().RulesK are considered
// this is ordered, the new x.GetGNode().Rules is returned,
func ForEachChild_InRules(x Ast, f func(Ast) Ast) map[string]Ast {
	hnew := map[string]Ast{}
	gn := x.GetGNode()
	if gn == nil || gn.Rules == nil {
		return nil
	}
	for _, k := range gn.RulesK {
		if v, exists := gn.Rules[k]; exists {
			if r := f(v); r != nil {
				hnew[k] = r
			}
		}
	}
	return hnew
}
