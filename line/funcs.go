package line

// import . "grepsuzette/joeson/core"
// import "fmt"
import (
	"grepsuzette/joeson/ast"
	"grepsuzette/joeson/helpers"
	"reflect"
	"strings"
)

type Lines []Line

func NewRankFromLines(rankname string, lines []Line, grammar *ast.Grammar) *ast.Rank {
	rank := ast.NewEmptyRank(rankname)
	grammar.SetRankIfEmpty(rank)
	for _, line := range lines {
		if il, ok := line.(ILine); ok {
			name, rule := il.ToRule(grammar, rank)
			rank.GetGNode().Include(name, rule)
		} else if ol, ok := line.(OLine); ok {
			choice := ol.ToRule(grammar, rank, OLineByIndexOrByName{index: helpers.NewNullInt(rank.Length())})
			rank.Append(choice)
			// } else if someAttr. But it won't be useful now
		} else {
			panic("Unknown type line, expected 'o' or 'i' line, got '" + line.String() + "' (" + reflect.TypeOf(line).String() + ")")
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
