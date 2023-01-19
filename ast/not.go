package ast

import (
	. "grepsuzette/joeson/colors"
	. "grepsuzette/joeson/core"
)

type Not struct {
	*GNode
	it Ast
}

func NewNot(it Ast) *Not {
	gn := NewGNode()
	not := &Not{gn, it}
	gn.Capture = false
	gn.Node = not
	return not
}

func (not *Not) GetGNode() *GNode        { return not.GNode }
func (not *Not) Prepare()                {}
func (not *Not) HandlesChildLabel() bool { return false }

func (not *Not) Parse(ctx *ParseContext) Ast {
	return Wrap(func(_ *ParseContext, _ Ast) Ast {
		pos := ctx.Code.Pos
		res := not.it.Parse(ctx)
		ctx.Code.Pos = pos
		if res != nil {
			return nil
		} else {
			return NewNativeUndefined()
		}
	}, not)(ctx)
}

func (not *Not) ContentString() string {
	return Yellow("!") + String(not.it)
}
func (not *Not) ForEachChild(f func(Ast) Ast) Ast {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	not.GetGNode().Rules = ForEachChild_InRules(not, f)
	if not.it != nil {
		not.it = f(not.it)
	}
	return not
}
