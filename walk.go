package joeson

// Depth-first Walk() function below accepts a WalkPrepost structure allowing optional walk callbacks during respectively
// the initial descent (Pre) and the climbing back (Post).
type WalkPrepost struct {
	Pre  func(node Parser, parent Parser) string // called during 🡖  .  "__stop__" to interrupt
	Post func(node Parser, parent Parser) string // called during 🡕 .
}

// Depth first walk of entire tree.
// `ast` is the node on which to start descending recursively.
// `parent` is available for algorithms needing it (just provide the
// father of `ast` or nil).
func Walk(ast Parser, parent Parser, prepost WalkPrepost) Parser {
	if prepost.Pre != nil {
		stop := prepost.Pre(ast, parent)
		if stop == "__stop__" {
			return ast
		}
	}
	ast.ForEachChild(func(child Parser) Parser {
		return Walk(child, ast, prepost)
	})
	if prepost.Post != nil {
		prepost.Post(ast, parent)
	}
	return ast
}

// ForEachChild specialization for []Parser arrays
func ForEachChild_Array(a []Parser, f func(Parser) Parser) []Parser {
	anew := []Parser{}
	for _, child := range a {
		if r := f(child); r != nil {
			anew = append(anew, r)
		} // else removed
	}
	return anew
}

// ForEachChild specialization for Parser's Rules
// working with RulesK will guarantee they are processed in order
func ForEachChildInRules(x Parser, f func(Parser) Parser) map[string]Parser {
	hnew := map[string]Parser{}
	gn := x.gnode()
	if gn.rules == nil {
		return nil
	}
	for _, name := range gn.rulesK {
		if parser, exists := gn.rules[name]; !exists {
			panic("assert")
		} else {
			if r := f(parser); r != nil {
				hnew[name] = r
			}
		}
	}
	return hnew
}
