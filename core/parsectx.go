package core

import (
	"errors"
	"fmt"
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/helpers"
	"strconv"
	"strings"
)

type stash struct {
	frames []*frame
	count  int
}

/**
  This is passed as the arg therebelow called `$` with the following methods (file core/funcs.go):
    - core.Stack($)
    - core.Loopify($)
    - core.PrepareResult($)
    - core.Wrap($)
   In the original joeson.coffee, it was named '$' and
   those functions names were prefixed by a dollar ($stack, $loopify etc)
*/
type ParseContext struct {
	grammar GrammarRuleCounter // grammar instance
	Code    *CodeStream
	stack   [1024]*frame
	// Result      Astnode    // joeson.coffee:625, see parseCode() in grammar.go. However a local var seems ok, see grammar.go
	Frames      [][]*frame // frames is 2d, dim is [filelen][grammar.NumRules]
	stackLength int
	counter     int
	SkipLog     bool
	Debug       bool
}

func NewParseContext(code *CodeStream, grammar GrammarRuleCounter, attrs ParseOptions) *ParseContext {
	// frames is 2d
	// frames[len(code.text) + 1][grammar.numRules]frame
	//                         ^---- +1 is to include EOF
	frames := make([][]*frame, len(code.text)+1)
	for i := range frames {
		// not using an explicit grammar.NumRules here, see how it goes
		frames[i] = make([]*frame, grammar.CountRules())
	}
	return &ParseContext{
		Code:        code,
		grammar:     grammar,
		Frames:      frames,
		stackLength: 0,
		counter:     0,
		SkipLog:     attrs.SkipLog,
	}
}

func (ctx *ParseContext) log(message string) {
	if !ctx.SkipLog {
		line := ctx.Code.Line()
		if Trace.FilterLine == -1 || line == Trace.FilterLine {
			codeSgmnt := White(strconv.Itoa(line)) + "," + strconv.Itoa(ctx.Code.Col())
			p := helpers.Escape(ctx.Code.Peek(NewPeek().BeforeChars(5)))
			codeSgmnt += "\t" + Black(helpers.PadRight(helpers.SliceString(p, len(p)-5, len(p)), 5))
			p = helpers.Escape(ctx.Code.Peek(NewPeek().AfterChars(20)))
			codeSgmnt += Green(helpers.PadLeft(helpers.SliceString(p, 0, 20), 20))
			if ctx.Code.Pos+20 < len(ctx.Code.text) {
				codeSgmnt += Black(">")
			} else {
				codeSgmnt += Black("]")
			}
			fmt.Printf("%s %s%s\n", codeSgmnt, Cyan(strings.Join(make([]string, ctx.stackLength), "| ")), message)
		}
	}
}

func (ctx *ParseContext) StackPeek(skip int) *frame {
	return ctx.stack[ctx.stackLength-1-skip]
}
func (ctx *ParseContext) StackPush(x Astnode) {
	ctx.stack[ctx.stackLength] = ctx.getFrame(x)
	ctx.stackLength++
}
func (ctx *ParseContext) StackPop() { ctx.stackLength-- }
func (ctx *ParseContext) getFrame(x Astnode) *frame {
	id := x.GetGNode().Id
	pos := ctx.Code.Pos
	posFrames := ctx.Frames[pos]
	frame := posFrames[id]
	if frame != nil {
		return frame
	} else {
		posFrames[id] = newFrame(pos, id)
		return posFrames[id]
	}
}

func (ctx *ParseContext) wipeWith(frame_ *frame, makeStash bool) *stash {
	// return *stash, or nil if !makeStash
	// default for makeStash was true in coffee
	if TimeStart != nil {
		TimeStart("wipewith")
	}
	if frame_.wipemask == nil {
		panic(errors.New("need frame.wipemask to know what to wipe"))
	}
	var stash_ []*frame
	if makeStash {
		// note: not using numRules below, see how it goes
		stash_ = make([]*frame, ctx.grammar.CountRules())
	} else {
		stash_ = nil
	}
	stashCount := 0
	pos := frame_.pos
	posFrames := ctx.Frames[pos]
	for i, bWipe := range frame_.wipemask {
		if !bWipe {
			continue
		}
		if makeStash {
			stash_[i] = posFrames[i]
		}
		posFrames[i] = nil // TODO make sure we can store nil, else use pointers
		stashCount++
	}
	if TimeEnd != nil {
		TimeEnd("wipewith")
	}
	if stash_ != nil {
		return &stash{frames: stash_, count: stashCount}
	} else {
		return nil
	}
}

func (ctx *ParseContext) restoreWith(stash_ *stash) {
	if TimeStart != nil {
		TimeStart("restorewith")
	}
	stashCount := stash_.count

	for i, frame := range stash_.frames {
		if frame == nil {
			continue
		}
		ctx.Frames[frame.pos][i] = frame
		stashCount--
		if stashCount == 0 {
			break
		}
	}
	if TimeEnd != nil {
		TimeEnd("restorewith")
	}
}
