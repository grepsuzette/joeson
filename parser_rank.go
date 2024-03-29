package joeson

import (
	"reflect"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// A rank is created from lines of rule. A grammar contains a rank.
type rank struct {
	*choice
}

// Create a Rank from the provided lines of rules. Use GrammarFromLines().
//
//	`optionalLazyGrammar` is a lazy callback specifying how to create or retrieve a grammar
//	from its cache, should the `lines` contain some string rules (SLine) needing to be compiled.
//	In general leave it nil to have the joeson_handcompiled grammar be used automatically.
func rankFromLines(lines []Line, rankname string, options GrammarOptions) *rank {
	var lazyGm *helpers.Lazy[*Grammar]
	if options.LazyGrammar != nil {
		lazyGm = options.LazyGrammar
	} else {
		lazyGm = helpers.LazyFromFunc[*Grammar](func() *Grammar {
			// Lazy, so this will only get called if the rules have string rules (SLine)
			// and optionalLazyGrammar was left empty. getgetRule()'s `case SLine:` is
			// the only place that needs a grammar.
			return NewJoesonWithOptions(options.TraceOptions)
		})
	}
	rank := newEmptyRank(rankname)
	for _, line := range lines {
		if ol, ok := line.(OLine); ok {
			choice := ol.toRule(rank, rank, oLineNaming{index: rank.Length()}, options.TraceOptions, lazyGm)
			rank.Append(choice)
		} else if il, ok := line.(ILine); ok {
			name, rule := il.toRule(rank, rank, options.TraceOptions, lazyGm)
			rank.Include(name, rule)
		} else {
			panic("Unknown type line, expected 'o' or 'i' line, got '" + line.stringIndent(0) + "' (" + reflect.TypeOf(line).String() + ")")
		}
	}
	return rank
}

func newEmptyRank(name string) *rank {
	x := &rank{newEmptyChoice()}
	x.getRule().name = name
	x.node = x
	return x
}

func (rank *rank) Length() int {
	return len(rank.choice.choices)
}

func (rank *rank) Append(node Parser)      { rank.choice.Append(node) }
func (rank *rank) getRule() *rule          { return rank.choice.getRule() }
func (rank *rank) prepare()                { rank.choice.prepare() }
func (rank *rank) handlesChildLabel() bool { return false }

func (rank *rank) String() string {
	var b strings.Builder
	b.WriteString(Blue("Rank("))
	first := true
	for _, it := range rank.choice.choices {
		if !first {
			b.WriteString(Blue(","))
		}
		b.WriteString(Red(it.getRule().name))
		first = false
	}
	b.WriteString(Blue(")"))
	return b.String()
}

func (rank *rank) forEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	ch := rank.choice.forEachChild(f) // see Choice.ForEachChild, which have the same @defineChildren
	rank.choice = ch.(*choice)
	return rank
}

func (rank *rank) parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		for _, choice := range rank.choice.choices {
			pos := ctx.Code.Pos()
			// Rank inherits from Choice in the original coffee implementation.
			// In coffee, the Parse function of Rank is bound to Rank,
			// In go, no inheritance, we inline the call instead.
			result := choice.parse(ctx)
			if result == nil {
				ctx.Code.SetPos(pos)
			} else {
				return result
			}
		}
		return nil
	}, rank)(ctx)
}
