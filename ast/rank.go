package ast

import (
	"fmt"
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/lambda"

	// "strconv"
	"strings"
)

type Rank struct {
	// *GNode
	*Choice
}

// init: in joeson.coffee is (@name, @choices=[], includes={})
// You really want to look at line.NewRankFromRules() instead
func NewEmptyRank(rankname string) *Rank {
	rank := Rank{ /*NewGNode(), */ NewEmptyChoice()}
	rank.GetGNode().Name = rankname
	return &rank
}

func (rank *Rank) Length() int {
	return len(rank.Choice.choices)
}

func (rank *Rank) Append(node Astnode)     { rank.Choice.Append(node) }
func (rank *Rank) GetGNode() *GNode        { return rank.Choice.GetGNode() }
func (rank *Rank) Prepare()                { rank.Choice.Prepare() }
func (rank *Rank) HandlesChildLabel() bool { return false }
func (rank *Rank) Labels() []string        { return rank.Choice.Labels() }
func (rank *Rank) Captures() []Astnode     { return rank.Choice.Captures() }
func (rank *Rank) ContentString() string {
	var b strings.Builder
	if !IsRule(rank) {
		return "oops, not a rule"
	}
	b.WriteString(LabelOrName(rank))
	b.WriteString(Blue("Rank("))
	a := lambda.Map(rank.Choice.choices, func(x Astnode) string {
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
	// we must first wall through rules, and then only through choices
	//  this is done in Choice.ForEachChild
	ch := rank.Choice.ForEachChild(f)
	rank.Choice = ch.(*Choice)
	return rank
}

// Rank.Parse is special because Rank inherits from Choice
// in original implementation. In coffee, the Parse function
// of Rank is bound to Rank, but the code is that of Choice.Parse.
// To replicate that here easily, we repeat Choice.Parse code instead.
func (rank *Rank) Parse(ctx *ParseContext) Astnode {
	// return rank.Choice.Parse(ctx) // no, it needs to be bound to rank as
	// below:
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
		for i, choice := range rank.Choice.choices {
			fmt.Printf("Rank..Choice n=%d %s\n", i, choice.ContentString())
			pos := ctx.Code.Pos
			result := choice.Parse(ctx)
			if result == nil {
				ctx.Code.Pos = pos
			} else {
				return result
			}
		}
		return nil
	}, rank)(ctx)
}
