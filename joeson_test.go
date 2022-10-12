package main

import "grepsuzette/joeson/helpers"

func S(a []astnode) astnode { newSequence(a) }
func P(value astnode, join astnode, min int, max int) astnode {
	return newPattern(value, join, min, max)	
}
func R(ref string) astnode { return newRef(ref, nil) }

// note: this grammar is the one from joeson.coffee,
// while the one in main.go is the one from  joeson_test.coffee.
// It is inverted during this period where we are working
// on nodes and hand definition.
func test() {
	rulez := rules(
		o("EXPR", rules(
			o_handCompiled(S(R("CHOICE"), R("_"))),
			o("CHOICE", rules(
				o_handCompiled(S(P(R("_PIPE")), P(R("SEQUENCE"),R("_PIPE"),2), P(R("_PIPE"))), 
					func (it_would_be_array_of_gnode) astnode { 
						return newChoice(it_would_be_array_of_gnode) }),
				o("SEQUENCE", rules(
					o_handCompiled(P(R("UNIT"), nil, 2), 
						func (it_would_be_WHAT) astnode {
							return newSequence(it_would_be_WHAT)),
						},
					o("UNIT", rules(

					))
				))
			))
		))
	)

}
