package main

import "errors"
import "fmt"
import . "grepsuzette/joeson/colors"
import "grepsuzette/joeson/helpers"
import "strconv"
import "strings"

type stash struct {
	frames []*frame
	count  int
}

/**
  This is passed as the arg therebelow called `$` with the following methods:
    _stack($)
  	_loopify($)
  	_prepareResult($)
  	_wrap($)
   In the original joeson.coffee, it was named '$' and
   those functions names were prefixed by a dollar ($stack, $loopify etc)
*/
type ParseContext struct {
	grammar     Grammar // grammar instance
	code        CodeStream
	stack       [1024]*frame
	result      astnode    // joeson.coffee:625, see parseCode() in grammar.go
	frames      [][]*frame // frames is 2d, dim is [filelen][grammar.numRules]
	stackLength int
	counter     int
	skipLog     bool
	opts        Opts // TODO
}

func newParseContext(code CodeStream, grammar Grammar, opts Opts) ParseContext {
	//       /----------2-dimensional-------------\
	// frames[len(code.text) + 1][grammar.numRules]frame
	//                         ^---- +1 is to include EOF
	frames := make([][]*frame, len(code.text)+1)
	for i := range frames {
		frames[i] = make([]*frame, grammar.numRules)
	}
	return ParseContext{
		code:        code,
		grammar:     grammar,
		opts:        opts,
		frames:      frames,
		stackLength: 0,
		counter:     0,
	}
}

func (ctx *ParseContext) log(message string) {
	if !ctx.skipLog {
		line := ctx.code.line()
		if trace.filterLine == -1 || line == trace.filterLine {
			codeSgmnt := White(strconv.Itoa(line)) + "," + strconv.Itoa(ctx.code.col())
			p := helpers.Escape(ctx.code.peek(Peek{beforeChars: helpers.NewNullInt(5)}))
			codeSgmnt += "\t" + Black(helpers.PadLeft(p[len(p)-5:], 5))
			p = helpers.Escape(ctx.code.peek(Peek{afterChars: helpers.NewNullInt(20)}))
			codeSgmnt += "\t" + Green(helpers.PadRight(p[0:20], 20))
			if ctx.code.pos+20 < len(ctx.code.text) {
				codeSgmnt += Black(">")
			} else {
				codeSgmnt += Black("]")
			}
			fmt.Printf("%s %s %s", codeSgmnt, Cyan(strings.Join(make([]string, ctx.stackLength), "| ")), message)
		}
	}
}
func (ctx *ParseContext) logIf(cond bool, message string) {
	if cond {
		ctx.log(message)
	}
}

func (ctx *ParseContext) stackPeek(skip int) frame {
	return *ctx.stack[ctx.stackLength-1-skip]
}
func (ctx *ParseContext) stackPush(node Node) {
	ctx.stackLength++
	ctx.stack[ctx.stackLength] = ctx.getFrame(node)
}
func (ctx *ParseContext) stackPop() { ctx.stackLength -= 1 }
func (ctx *ParseContext) getFrame(node Node) *frame {
	id := node.id // id is an int, it is incremented in Grammar: node.id = @numRules++
	pos := ctx.code.pos
	posFrames := ctx.frames[pos]
	frame := posFrames[id]
	if frame != nil {
		posFrames[id] = newFrame(pos, id)
		return posFrames[id]
	} else {
		return frame
	}
}

func (ctx *ParseContext) wipeWith(frame_ *frame, makeStash bool) *stash {
	// return *stash, or nil if !makeStash
	// default for makeStash was true in coffee
	if timeStart != nil {
		timeStart("wipewith")
	}
	if frame_.wipemask == nil {
		panic(errors.New("Need frame.wipemask to know what to wipe"))
	}
	var stash_ []*frame
	if makeStash {
		stash_ = make([]*frame, ctx.grammar.numRules)
	} else {
		stash_ = nil
	}
	stashCount := 0
	pos := frame_.pos
	posFrames := ctx.frames[pos]
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
	if timeEnd != nil {
		timeEnd("wipewith")
	}
	if stash_ != nil {
		return &stash{frames: stash_, count: stashCount}
	} else {
		return nil
	}
}

func (ctx *ParseContext) restoreWith(stash_ *stash) {
	if timeStart != nil {
		timeStart("restorewith")
	}
	stashCount := stash_.count

	for i, frame := range stash_.frames {
		if frame == nil {
			continue
		}
		ctx.frames[frame.pos][i] = frame
		stashCount--
		if stashCount == 0 {
			break
		}
	}
	if timeEnd != nil {
		timeEnd("restorewith")
	}
}
