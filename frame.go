package main

import (
	"grepsuzette/joeson/helpers"
)

// in joeson there was this comment:
// # { pos:{ (node.id):{id,result,pos,endPos,stage,...(same object as in stack)}... } }
type frame struct {
	result    astnode
	endPos    helpers.NullInt // can be left undefined
	loopStage helpers.NullInt // can be left undefined
	wipemask  []bool          // len = ctx.grammar.numRules
	//TODO delete, its not in the original work. subframes map[string]frame
	pos   int
	id    int
	param any // used in ref.go or joeson.coffee:536
}

func (f frame) toString() string {
	return "TODO frame.toString"
}
func (fr *frame) cacheSet(result *Result, endPos int) {
	fr.result = result
	if endPos < 0 {
		fr.endPos.Unset()
	} else {
		fr.endPos.Set(endPos)
	}
}
func newFrame(pos int, id int) *frame {
	return &frame{
		result:   nil,
		pos:      pos,
		id:       id,
		wipemask: nil,
		param:    nil,
	}
}
