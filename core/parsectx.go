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

// One rule is parsed by a different ParseContext, with an increasing Counter.
// Suppose a rule is defined like so: i(named("LABEL", "'&' | '@' | WORD")),
// its Code will contain "'&' | '@' | WORD" and start at position 0.
type ParseContext struct {
	TraceOptions     // grammar.TraceOptions at the moment this context is created
	ParseOptions     // Defined arbitrarily within a rule, e.g. I("INT", "/[0-9]+/", someCb, ParseOptions{SkipLog: false}), and then passed to the ParseContext.
	Counter      int // the iteration counter that is shown in the stack trace, useful when debugging
	Code         *CodeStream

	numRules    int        // grammar.numRules at the moment context is created. If no grammar (esp., when joeson rules are parsed the first time) pass 0.
	Frames      [][]*frame // 2D. frames[len(code.text) + 1][grammar.numRules]frame. Though it's public only core.grammar should access it.
	stack       [1024]*frame
	stackLength int
	loopStack   []string
}

// numRules: grammar numRules at the moment context is created. If no grammar (esp. when
// joeson rules are parsed the very first time) pass 0.
func NewParseContext(code *CodeStream, numRules int, attrs ParseOptions, opts TraceOptions) *ParseContext {
	// frames is 2d
	// frames[len(code.text) + 1][grammar.numRules]frame
	//                         ^---- +1 is to include EOF
	frames := make([][]*frame, len(code.text)+1)
	for i := range frames {
		frames[i] = make([]*frame, numRules)
	}
	return &ParseContext{
		Code:         code,
		numRules:     numRules,
		TraceOptions: opts,
		ParseOptions: attrs,
		Frames:       frames,
		stackLength:  0,
		Counter:      0,
	}
}

func (ctx *ParseContext) log(message string, opts TraceOptions) {
	if !ctx.SkipLog {
		line := ctx.Code.Line()
		if opts.FilterLine == -1 || line == opts.FilterLine {
			codeSgmnt := White(strconv.Itoa(line) + "," + strconv.Itoa(ctx.Code.Col()))
			p := helpers.Escape(ctx.Code.Peek(NewPeek().BeforeChars(5)))
			codeSgmnt += "\t" + BoldBlack(helpers.PadRight(helpers.SliceString(p, len(p)-5, len(p)), 5))
			p = helpers.Escape(ctx.Code.Peek(NewPeek().AfterChars(20)))
			codeSgmnt += Green(helpers.PadLeft(helpers.SliceString(p, 0, 20), 20))
			if ctx.Code.Pos+20 < len(ctx.Code.text) {
				codeSgmnt += BoldBlack(">")
			} else {
				codeSgmnt += BoldBlack("]")
			}
			fmt.Printf("%s %s%s\n", codeSgmnt, Cyan(strings.Join(make([]string, ctx.stackLength), "| ")), message)
		}
	}
}

func (ctx *ParseContext) loopStackPush(name string) { ctx.loopStack = append(ctx.loopStack, name) }
func (ctx *ParseContext) loopStackPop()             { ctx.loopStack = ctx.loopStack[:len(ctx.loopStack)-1] }
func (ctx *ParseContext) StackPeek(skip int) *frame { return ctx.stack[ctx.stackLength-1-skip] }
func (ctx *ParseContext) StackPush(x Ast) {
	ctx.stack[ctx.stackLength] = ctx.getFrame(x)
	ctx.stackLength++
}
func (ctx *ParseContext) StackPop() { ctx.stackLength-- }

func (ctx *ParseContext) getFrame(x Ast) *frame {
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
		stash_ = make([]*frame, ctx.numRules)
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
		posFrames[i] = nil
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
