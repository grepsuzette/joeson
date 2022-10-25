package ast

import . "grepsuzette/joeson/core"
import . "grepsuzette/joeson/colors"

type Not struct {
	*GNode
	it Astnode
}

func NewNot(it Astnode) *Not {
	g := NewGNode()
	g.Capture = false
	return &Not{g, it}
}

func (not *Not) GetGNode() *GNode        { return not.GNode }
func (not *Not) Prepare()                {}
func (not *Not) HandlesChildLabel() bool { return false }
func (not *Not) Labels() []string        { return not.GNode.Labels() }
func (not *Not) Captures() []Astnode     { return not.GNode.Captures() }

func (not *Not) Parse(ctx *ParseContext) Astnode {
	return Wrap(func(_ *ParseContext) Astnode {
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
	return ShowLabelOrNameIfAny(not) + Yellow("!") + not.it.ContentString()
}
func (not *Not) ForEachChild(f func(Astnode) Astnode) Astnode {
	// @defineChildren
	//   rules:      {type:{key:undefined,value:{type:GNode}}}
	//   it:         {type:GNode}
	if not.it != nil {
		not.it = f(not.it)
	}
	not.GetGNode().Rules = ForEachChild_MapString(not.GetGNode().Rules, f)
	return not
}
