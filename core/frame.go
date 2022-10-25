package core

import "grepsuzette/joeson/helpers"

// in joeson there was this comment:
// # { pos:{ (node.id):{id,result,pos,endPos,stage,...(same object as in stack)}... } }
type frame struct {
	Result    Astnode
	endPos    helpers.NullInt // can be left undefined
	loopStage helpers.NullInt // can be left undefined
	wipemask  []bool          // len = ctx.grammar.numRules
	//TODO delete, its not in the original work. subframes map[string]frame
	pos   int
	id    int
	Param Astnode // used in ref.go or joeson.coffee:536
}

func (f frame) toString() string {
	return "TODO frame.toString"
}
func (fr *frame) cacheSet(result Astnode, endPos int) {
	fr.Result = result
	if endPos < 0 {
		fr.endPos.Unset()
	} else {
		fr.endPos.Set(endPos)
	}
}
func newFrame(pos int, id int) *frame {
	return &frame{
		Result:   nil,
		pos:      pos,
		id:       id,
		wipemask: nil,
		Param:    nil,
	}
}
