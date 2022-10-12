package main

import "errors"
import . "grepsuzette/joeson/colors"
import "reflect"

/*
# Main external access.
# I dunno if Grammar should be a GNode or not. It
# might come in handy when embedding grammars
# in some glue language.
  jae 2012-04-18 in original joeson.coffee
*/
type Grammar struct {
	GNode
	rank     Rank
	NumRules int
	// id2Rule: slow lookup for debugging...
	id2Rule map[int]astnode // node.id = @numRules++; @id2Rule[node.id] = node in joeson.coffee:605
}

// MAIN GRAMMAR PARSE FUNCTION
func (gm Grammar) ParseString(sCode string, opts *Opts) astnode {
	return gm.ParseCode(newCodeStream(sCode))
}

func emptyGrammar() Grammar { return Grammar{NewGNode(), NewEmptyRank(), 0} }
func NewGrammarFromRules(rules []iorule) (Grammar, error) {
	gr := emptyGrammar()
	gr.rank = RankFromLines("__grammar__", rules)
	gr.numRules = 0

	// from joeson.coffee
	// "TODO refactor into translation passes."

	// Merge Choices with just a single choice.
	/* it's an optimization which can wait until after something works :)
		gr.walk(walk{
			pre: func(node astnode, parent astnode, desc, key string, index) {
				if choice, ok := node.(Choice); ok && len(choice.choices) == 1 {
					// Merge label
					if choice.choices[0].GNode.label == "" {
						choice.choices[0].GNode.label = choice.GNode.label
					}
					// Merge included rules
					if len(choice.GNode.rules) > 0 {
						for k, v := range choice.GNode.rules {
							choice.choices[0].GNode.rules[k] = v
						}
					}
					// Replace with grandchild
					if index  > -1 {
						// TODO check hypothesis
						// hypothesis is if index is provided, it is an array
						reflect.ValueOf(parent).FieldByName(key).FieldByIndex(index)
					} else {
	i					parent.GetGNode().SetArbitraryField(key, choice.Choice.choices[0])
						// tricky, needs & and Elem(). Also, field name must be
						// exported (capitalized first letter)
						// https://stackoverflow.com/questions/6395076/using-reflect-how-do-you-set-the-value-of-a-struct-field

						reflect.ValueOf(&parent).Elem().FieldByName(key).se
					}
					parent
					reflect.S
				}
			}
		})
	*/

	// Connect all the nodes and collect dereferences into @rules
	gr.walk(walk{
		pre: func(node astnode, parent astnode) {
			// initialize gnode and gparent to not have to write
			// node.GetGNode() all the time
			gnode := node.GetGNode()
			var gparent GNode = nil
			if parent != nil {
				gparent = parent.GetGNode()
			}
			// done, now the code
			// sanity check: it must have no parent yet if it's not a rule
			//  if it has parent it means there's a cycle
			if gnode.parent != nil && gnode != gnode.rule {
				panic("Grammar tree should be a DAG, nodes should not be referenced more than once.")
			}
			gnode.grammar = gr
			gnode.parent = parent
			// inline rules are special
			if gnode.inlineLabel != nil { // TODO this one is strange
				// I see it in rules, but don't
				// see where it would be set.
				// Suspect this is never called.
				// Easy to test in fact
				gnode.rule = node
				gparent.rule.GNode.include(gnode.inlineLabel, node)
			} else {
				// set node.rule, the root node for this rule
				if gnode.rule == nil && gparent != nil {
					gnode.rule = gparent.rule
				}
			}
		},
		post: func(node astnode, parent astnode) {
			// initialize gnode and gparent to not have to write
			// node.GetGNode() all the time
			gnode := node.GetGNode()
			var gparent GNode = nil
			if parent != nil {
				gparent = parent.GetGNode()
			}
			// done, now the code
			if node == gnode.rule {
				gr.rules[gnode.Name] = node
				gr.NumRules++
				gnode.id = gr.NumRules
				gr.id2Rule[gnode.id] = node
				if trace.loop { // print out id->rulename for convenience
					// TODO  in coffee it's just  console.log( blabla.. node) --v
					fmt.Println(Red(strconv.Itoa(gnode.id)) + ":\t" + node.ContentString())
				}
			}
		},
	})
	// Prepare all the nodes, child first.
	gr.walk(walk{
		post: func(node astnode, parent astnode) {
			return node.Prepare()
		},
	})
	return gr, nil
}

func (gm Grammar) Prepare()           {}
func (gm Grammar) HandlesChildLabel() { return false }
func (gm Grammar) ParseCode(code CodeStream, aOpts ...Option) astnode {
	opts := Opts{}
	// note: opts.returnContext is absent from this implementation
	for _, f := range aOpts {
		opts = f(opts)
	}
	ctx := NewParseContext(code, gm, opts)
	if opts.debug {
		// temporarily enable stack tracing
		oldTrace := clone(trace) // TODO. Also see 10lines below
		trace.stack = true
	}

	// parse
	result := gm.rank.Parse(ctx)
	if ctx.result != nil {
		ctx.result.code = code
	}

	// undo temprary stack tracing
	if opts.debug {
		trace = oldTrace
	}

	// if parse is incomplete, compute error message
	if ctx.code.pos != len(ctx.code.text) {
		// find the maximum parsed entity
		maxAttempt := ctx.code.pos
		maxSuccess := ctx.code.pos
		// (below TODO was originally in joeson.coffee, not sure what is meant)
		// "TODO for x, i in something from index to index by skip"
		for pos := ctx.code.pos; pos < len(ctx.frames); pos++ {
			posFrames := ctx.frames[pos]
			for _, frame := range posFrames {
				if frame != nil {
					maxAttempt = pos
					if frame.result != nil {
						maxSuccess = pos
						break
					}
				}
			}
		}
		sErr = "Error parsing at char:" + maxSuccess + "=(line:" + ctx.code.posToLine(maxSuccess) + ",col:" + ctx.code.posToCol(maxSuccess) + ")."
		sErr += "\nDetails:\n"
		sErr += Green("OK") + "/"
		sErr += Yellow("Parsing") + "/"
		sErr += Red("Suspect") + "/"
		sErr += White("Unknown") + "\n\n"
		sErr += Green(ctx.code.peek(Peek{beforeLines: 2}))
		sErr += Yellow(ctx.code.peek(Peek{afterChars: maxSuccess - ctx.code.pos}))
		ctx.code.pos = maxSuccess
		sErr += Red(ctx.code.peek(Peek{afterChars: maxAttempt - ctx.code.pos})) + "/"
		ctx.code.pos = maxAttempt
		sErr += White(ctx.code.peek(Peek{afterLines: 2})) + "\n"
		panic(errors.New(sErr))
	}

	// joeson.coffee has a opts.returnContext but won't implement it
	return ctx.result
}

func (gm Grammar) Labels() []string    { return gm.GNode.Labels() }
func (gm Grammar) Captures() []astnode { return gm.GNode.Captures() }
func (gm Grammar) ContentString() string {
	// TODO:                 is below ----v ContentString() correct? original is just @rank + magenta('}')
	return Magenta("GRAMMAR{") + gm.rank.ContentString() + Magenta("}")
}
