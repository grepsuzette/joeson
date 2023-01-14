package core

import "grepsuzette/joeson/helpers"

type frame struct {
	Result    Ast
	endPos    helpers.NullInt // can be left undefined
	loopStage helpers.NullInt // can be left undefined
	wipemask  []bool          // len = ctx.grammar.numRules
	pos       int
	id        int
	Param     Ast // used in ref.go or joeson.coffee:536
}

func (f frame) toString() string {
	return "N/A frame.toString"
}

func (fr *frame) cacheSet(result Ast, endPos int) {
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
