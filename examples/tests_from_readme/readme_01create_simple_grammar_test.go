package readmetest

import (
	"fmt"
	"testing"

	j "github.com/grepsuzette/joeson"
)

func TestReadmeCreatingSimpleGrammar(t *testing.T) {
	aa := [][]j.Line{
		{
			// Invoke GrammarFromLines with its title and a list of rules and you
			// get a dynamically compiled grammar:
			o(named("INPUT", `'hi' _ NAME`)),
			i(named("_", `' '*`)),
			i(named("NAME", `[a-zA-Z]*`), func(it j.Ast) j.Ast {
				return it.(*j.NativeArray).Concat()
			}),
		},
		{
			// When a rule has a parse function, it is called with the production
			// of the joeson parsers for that rule. Let's add a parse function to
			// our entry rule:
			o(named("INPUT", `'hi' _ NAME`), func(it j.Ast) j.Ast {
				fmt.Printf("I am not %s.\n", it.(j.NativeString).String())
				return it
			}),
			i(named("NAME", `[a-zA-Z]*`), func(it j.Ast) j.Ast {
				return it.(*j.NativeArray).Concat()
			}),
			i(named("_", `' '*`)),
		},
		{
			// It's worth noting that "INPUT" is not referenced
			// anywhere in the grammar. The entry rule could therefore
			// have been declared like this instead:
			o(`'hi' _ NAME`, func(it j.Ast) j.Ast {
				fmt.Printf("I am not %s.\n", it.(j.NativeString).String())
				return it
			}),
			i(named("NAME", `[a-zA-Z]*`), func(it j.Ast) j.Ast {
				return it.(*j.NativeArray).Concat()
			}),
			i(named("_", `' '*`)),
		},
	}

	// try these grammars successively
	// those are variants used in "Creating a simple grammar"
	for _, rules := range aa {
		gm := j.GrammarFromLines("my grammar", rules)
		if j.IsParseError(gm.ParseString("hi amigo")) {
			t.Fail()
		}
		if !j.IsParseError(gm.ParseString("bye bye")) {
			t.Fail()
		}
	}
}
