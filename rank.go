package main

import "regexp"
import "strings"
import "grepsuzette/joeson/lambda"

type Rank struct {
	Choice
}

// note: the init: in joeson.coffee contains more
// but is never used. Don't implement a NewRank()!
func NewEmptyRank(rankname string) Rank {
	rank := Rank{NewEmptyChoice()}
	rank.Choice.GNode.Name = rankname
	return rank
}

func NewRankFromLines(rankname string, lines []iorule) Rank {
	rank := NewEmptyRank(rankname)
	for idx, line := range lines { // idx necessary for unnamed rules, e.g. 4th unnamed rule could be called "PARENTNAME[3]"
		if line.isO { // OLine
			choice := line.toRuleWithIndex(rank, len(rank.choices))
			rank.Choice.choices = append(rank.Choice.choices, choice)
		} else { // ILine
			for name, rule := range line.toRules() {
				rank.Choice.GNode.include(name, rule)
			}
		}
	}
	return rank
}

func (rank Rank) GetGNode() GNode                 { return rank.Choice.GNode }
func (rank Rank) Prepare()                        {}
func (rank Rank) HandlesChildLabel()              { return false }
func (rank Rank) Labels() []string                { return rank.GNode.Labels() }
func (rank Rank) Captures() []astnode             { return rank.GNode.Captures() }
func (rank Rank) Parse(ctx *ParseContext) astnode { return rank.Choice.Parse(ctx) }
func (rank Rank) ContentString() string {
	var b strings.Builder
	b.WriteString(Blue("Rank("))
	a := lambda.Map(rank.Choice.choices, func(nn astnode) {
		return Red(nn.GetGNode().Name)
	})
	b.WriteString(strings.Join(a, Blue(",")))
	b.WriteString(Blue(")"))
	return b.String()
}
