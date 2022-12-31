package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
	"grepsuzette/joeson/lambda"
	"strings"
)

type Rank struct {
	*Choice
}

// See also line.NewRankFromLines()
func NewEmptyRank(name string) *Rank {
	rank := &Rank{NewEmptyChoice()}
	rank.GetGNode().Name = name
	rank.GetGNode().Node = rank
	return rank
}

func (rank *Rank) Length() int {
	return len(rank.Choice.choices)
}

func (rank *Rank) Append(node Astnode)     { rank.Choice.Append(node) }
func (rank *Rank) GetGNode() *GNode        { return rank.Choice.GetGNode() }
func (rank *Rank) Prepare()                { rank.Choice.Prepare() }
func (rank *Rank) HandlesChildLabel() bool { return false }

func (rank *Rank) ContentString() string {
	var b strings.Builder
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
	ch := rank.Choice.ForEachChild(f)
	rank.Choice = ch.(*Choice)
	return rank
}

func (rank *Rank) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
		for _, choice := range rank.Choice.choices {
			pos := ctx.Code.Pos
			// Rank inherits from Choice in the original coffee implementation.
			// In coffee, the Parse function of Rank is bound to Rank,
			// In go, no inheritance, we inline the call instead.
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
