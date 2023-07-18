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
		lazyGm = helpers.NewLazyFromFunc[*Grammar](func() *Grammar {
			// Lazy, so this will only get called if the rules have string rules (SLine)
			// and optionalLazyGrammar was left empty. getRule()'s `case SLine:` is
			// the only place that needs a grammar.
			return NewJoesonWithOptions(options.TraceOptions)
		})
	}
	ranke := newEmptyRank(rankname)
	for _, line := range lines {
		if ol, ok := line.(OLine); ok {
			choice := ol.toRule(ranke, ranke, oLineByIndexOrName{
				index: helpers.NewNilableInt(ranke.Length()),
			}, options.TraceOptions, lazyGm)
			ranke.Append(choice)
		} else if il, ok := line.(ILine); ok {
			name, rule := il.toRule(ranke, ranke, options.TraceOptions, lazyGm)
			ranke.gnode().Include(name, rule)
		} else {
			panic("Unknown type line, expected 'o' or 'i' line, got '" + line.stringIndent(0) + "' (" + reflect.TypeOf(line).String() + ")")
		}
	}
	return ranke
}

func newEmptyRank(name string) *rank {
	x := &rank{newEmptyChoice()}
	x.SetName(name)
	x.gnode().node = x
	return x
}

// ranke is used below to differentiate the var from the type. It means nothing special.

func (rank *rank) Length() int {
	return len(rank.choice.choices)
}

func (rank *rank) Append(node Parser)      { rank.choice.Append(node) }
func (rank *rank) gnode() *gnodeimpl       { return rank.choice.gnode() }
func (rank *rank) prepare()                { rank.choice.prepare() }
func (rank *rank) handlesChildLabel() bool { return false }

func (rank *rank) ContentString() string {
	var b strings.Builder
	b.WriteString(blue("Rank("))
	a := helpers.AMap(rank.choice.choices, func(x Parser) string {
		return red(x.Name())
	})
	b.WriteString(strings.Join(a, blue(",")))
	b.WriteString(blue(")"))
	return b.String()
}

func (rank *rank) ForEachChild(f func(Parser) Parser) Parser {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   choices:    {type:[type:GNode]}
	ch := rank.choice.ForEachChild(f) // see Choice.ForEachChild, which have the same @defineChildren
	rank.choice = ch.(*choice)
	return rank
}

func (rank *rank) Parse(ctx *ParseContext) Ast {
	return wrap(func(_ *ParseContext, _ Parser) Ast {
		for _, choice := range rank.choice.choices {
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
