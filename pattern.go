package main

import "grepsuzette/joeson/lambda"

type Pattern struct {
	GNode
	value astnode // in coffee, declared as `value: {type:GNode}`
	join  astnode // in coffee, declared as `join:  {type:GNode}`
	min   int     // -1 if unspec.
	max   int     // -1 if unspec.
}

func NewPattern(value astnode, join astnode, min int, max int) Pattern {
	gnode := NewGNode()
	if value != nil {
		gnode.capture = value.GetGNode().capture
	}
	return Pattern{NewGNode(), value, join, min, max}
}
func (patt Pattern) GetGNode() GNode { return patt.GNode }
func (patt Pattern) Parse(ctx *ParseContext) astnode {
	return patt.GNode._wrap(func(_, _) astnode {
		// TODO see what to do about the @$wrap
		matches := make([]string, 0)
		pos := ctx.code.pos
		resV := patt.value.Parse(ctx)
		if resV == nil {
			ctx.code.pos = pos
			if patt.min > 0 {
				return nil
			}
			return NewNativeArray([]astnode{})
		}
		matches = append(matches, resV)
		for true {
			pos2 := ctx.code.pos
			resJ := patt.join.Parse(ctx)
			if patt.join != nil {
				// return nil to revert pos
				if resJ == nil {
					ctx.code.pos = pos2
					break
				}
			}
			resV := patt.value.Parse(ctx)
			// return nil to revert pos
			if resV == nil {
				ctx.code.pos = pos2
				break
			}
			matches = append(matches, resV)
			if patt.max > -1 && len(matches) >= patt.max {
				break
			}
		}
		if patt.min > -1 && patt.min > len(matches) {
			ctx.code.pos = pos
			return nil
		}
		a := lambda.Map(matches, func(s) { return NewNativeString(s) })
		return NewNativeArray(a)
	})(patt, ctx)
}

func (patt Pattern) ContentString() string { return "TODO" }
func (patt Pattern) HandlesChildLabel()    { return false }
func (patt Pattern) Labels() []string      { return patt.GNode.Labels() }
func (patt Pattern) Captures() []astnode   { return patt.GNode.Captures() }
func (patt Pattern) Prepare()              {}
