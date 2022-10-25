package line

// import . "grepsuzette/joeson/core"
// import "fmt"
import "grepsuzette/joeson/ast"
import "strings"

type Lines []Line

func NewRankFromLines(rankname string, lines []Line, grammar *ast.Grammar) *ast.Rank {
	rank := ast.NewEmptyRank(rankname)
	grammar.SetRank(rank)
	for i, line := range lines {
		if il, ok := line.(ILine); ok {
			// fmt.Println(il)
			name, rule := il.ToRule(grammar, rank)
			rank.GetGNode().Include(name, rule)
		} else if ol, ok := line.(OLine); ok {
			rank.Append(ol.ToRuleWithIndex(rankname, i, grammar))
		} else {
			panic("expect only o and i lines")
		}
	}
	return rank
}

func NewGrammarFromLines(name string, lines []Line) *ast.Grammar {
	gm := ast.NewEmptyGrammarNamed(name)
	NewRankFromLines(name, lines, gm)
	gm.Postinit()
	return gm
}

func (a Lines) String() string {
	var b strings.Builder
	for _, line := range a {
		b.WriteString(line.String() + "\n")
	}
	return b.String()
}
