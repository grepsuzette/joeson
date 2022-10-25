package ast

import "strings"
import "grepsuzette/joeson/lambda"
import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"

type Rank struct {
	*GNode
	*Choice
}

// init: in joeson.coffee is (@name, @choices=[], includes={})
// You really want to look at line.NewRankFromRules() instead
func NewEmptyRank(rankname string) *Rank {
	rank := Rank{NewGNode(), NewEmptyChoice()}
	rank.GNode.Name = rankname
	return &rank
}

func (rank *Rank) Append(node Astnode) { rank.Choice.Append(node) }

func (rank *Rank) GetGNode() *GNode                { return rank.GNode }
func (rank *Rank) Prepare()                        {}
func (rank *Rank) HandlesChildLabel() bool         { return false }
func (rank *Rank) Labels() []string                { return rank.GNode.Labels() }
func (rank *Rank) Captures() []Astnode             { return rank.GNode.Captures() }
func (rank *Rank) Parse(ctx *ParseContext) Astnode { return rank.Choice.Parse(ctx) }
func (rank *Rank) ContentString() string {
	var b strings.Builder
	b.WriteString(ShowLabelOrNameIfAny(rank))
	b.WriteString(Blue("Rank("))
	a := lambda.Map(rank.Choice.a, func(x Astnode) string {
		return Red(x.GetGNode().Name)
	})
	b.WriteString(strings.Join(a, Blue(",")))
	b.WriteString(Blue(")"))
	return b.String()
}
func (rank *Rank) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	rank.GetGNode().Rules = ForEachChild_MapString(rank.GetGNode().Rules, f)
	if rank.Choice != nil {
		r, ok := f(rank.Choice).(*Choice)
		if ok {
			rank.Choice = r
		} else {
			panic("unable to get a *Rank")
		}
	}
	return rank
}
