package joeson

import (
	"grepsuzette/joeson/helpers"
	"reflect"
)

// General options to build a grammar
type GrammarOptions struct {
	// Those are the options governing what is traced or not during the initialization or the parsing
	TraceOptions TraceOptions

	// A Lazy of *Grammar specifying how to create or retrieve a grammar
	//  from its cache, should the `lines` contain some string rules (SLine) needing to be compiled.
	//  In general leave it nil to have the joeson_handcompiled grammar to be used automatically.
	LazyGrammar *helpers.Lazy[*Grammar]
}

// Make a new grammar from the rules in `lines`. See also NewJoeson()
// Options can be omitted or specified like so: `GrammarOptions{}`
func GrammarFromLines(lines []Line, name string, options ...GrammarOptions) *Grammar {
	var opts GrammarOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = GrammarOptions{
			TraceOptions: DefaultTraceOptions(),
			LazyGrammar:  nil,
		}
	}
	rank := rankFromLines(lines, name, opts)
	newgm := NewEmptyGrammarWithOptions(opts.TraceOptions)
	newgm.SetRankIfEmpty(rank)
	// The name is also set afterwards in the coffeescript version
	newgm.GetGNode().Name = name
	newgm.Postinit()
	return newgm
}

// Create a Rank from rule Lines.
//  `optionalLazyGrammar` is a lazy callback specifying how to create or retrieve a grammar
//  from its cache, should the `lines` contain some string rules (SLine) needing to be compiled.
//  In general leave it nil to have the joeson_handcompiled grammar be used automatically.
func rankFromLines(lines []Line, rankname string, options GrammarOptions) *Rank {
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
	rank := NewEmptyRank(rankname)
	for _, line := range lines {
		if ol, ok := line.(OLine); ok {
			choice := ol.toRule(rank, rank, OLineByIndexOrName{
				index: helpers.NewNilableInt(rank.Length()),
			}, options.TraceOptions, lazyGm)
			rank.Append(choice)
		} else if il, ok := line.(ILine); ok {
			name, rule := il.toRule(rank, rank, options.TraceOptions, lazyGm)
			rank.GetGNode().Include(name, rule)
		} else {
			panic("Unknown type line, expected 'o' or 'i' line, got '" + line.StringIndent(0) + "' (" + reflect.TypeOf(line).String() + ")")
		}
	}
	return rank
}
