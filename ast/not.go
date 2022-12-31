package ast

import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"

type Not struct {
	*GNode
	it Astnode
}

func NewNot(it Astnode) *Not {
	gn := NewGNode()
	not := &Not{gn, it}
	gn.Capture = false
	gn.Node = not
	return not
}

func (not *Not) GetGNode() *GNode        { return not.GNode }
func (not *Not) Prepare()                {}
func (not *Not) HandlesChildLabel() bool { return false }
func (not *Not) Labels() []string        { panic("z") }
func (not *Not) Captures() []Astnode     { panic("z") }

func (not *Not) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext, _ Astnode) Astnode {
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
	return Yellow("!") + Prefix(not.it) + not.it.ContentString()
}
func (not *Not) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	not.GetGNode().Rules = ForEachChild_InRules(not, f)
	if not.it != nil {
		not.it = f(not.it)
	}
	return not
}
