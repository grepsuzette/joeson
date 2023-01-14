package line

import (
	"grepsuzette/joeson/ast"
	// "grepsuzette/joeson/core"
	"grepsuzette/joeson/helpers"
	"reflect"
	"strings"
)

type Lines []Line

func NewRankFromLines(rankname string, lines []Line, grammar *ast.Grammar) *ast.Rank {
	rank := ast.NewEmptyRank(rankname)
	for _, line := range lines {
		if ol, ok := line.(OLine); ok {
			choice := ol.ToRule(grammar, rank, OLineByIndexOrName{index: helpers.NewNilableInt(rank.Length())})
			rank.Append(choice)
		} else if il, ok := line.(ILine); ok {
			name, rule := il.ToRule(grammar, rank)
			rank.GetGNode().Include(name, rule)
		} else {
			panic("Unknown type line, expected 'o' or 'i' line, got '" + line.StringIndent(0) + "' (" + reflect.TypeOf(line).String() + ")")
		}
	}
	return rank
}

// The returned grammar is a new one, while arg `grammar` is the
//  grammar used to parse this new grammar (usually this would be
//  the one in ast/handcompiled)
func NewGrammarFromLines(name string, lines []Line, grammar *ast.Grammar) *ast.Grammar {
	rank := NewRankFromLines(name, lines, grammar)
	newgm := ast.NewEmptyGrammarNamed(name)
	newgm.SetRankIfEmpty(rank)
	newgm.Postinit()
	return newgm
}

func (a Lines) String() string {
	var b strings.Builder
	for _, line := range a {
		b.WriteString(line.StringIndent(0) + "\n")
	}
	return b.String()
}
