package joeson

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/grepsuzette/joeson/helpers"
)

// As in the initial implementation, ParseContext does not know Grammar. It's
// just a context.
type ParseContext struct {
	TraceOptions        // grammar.TraceOptions at the moment this context is created.
	Counter      int    // [For debugging] iteration counter shown with TRACE=stack. Useful with conditional breakpoints
	GrammarName  string // [For debugging] set in grammar.Parse to the value of grammar.Name(). Useful for conditional breakpoints (typically in packrat loopify()) to only break when your final grammar is being used to parse anything. See docs/diffing.md # debugging methodology
	Code         *CodeStream

	numRules     int
	frames       [][]*frame // 2D: [len(code.text) + 1][numRules]
	stack        [1024]*frame
	stackLength  int
	loopStack    []string
	parseOptions ParseOptions // Defined arbitrarily within a rule (setParseOptions()), e.g. in I("INT", "/[0-9]+/", someCb, ParseOptions{SkipLog: false}), and then passed to the ParseContext.
}

// Create a new parse context.
// numRules: grammar numRules at the moment context is created (can be 0 before the very first grammar is created)
func newParseContext(code *CodeStream, numRules int, opts TraceOptions) *ParseContext {
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
		frames:       frames,
		stackLength:  0,
		Counter:      0,
	}
}

// ParseOptions are only set from within the callback within a rule.
func (ctx *ParseContext) setParseOptions(opts ParseOptions) *ParseContext {
	ctx.parseOptions = opts
	return ctx
}

func (ctx *ParseContext) String() string {
	line := ctx.Code.Line()
	codeSgmnt := white(strconv.Itoa(line) + "," + strconv.Itoa(ctx.Code.Col()))
	p := helpers.Escape(ctx.Code.Peek(NewPeek().BeforeChars(5)))
	codeSgmnt += "\t" + boldBlack(helpers.PadRight(helpers.SliceString(p, len(p)-5, len(p)), 5))
	p = helpers.Escape(ctx.Code.Peek(NewPeek().AfterChars(20)))
	codeSgmnt += green(helpers.PadLeft(helpers.SliceString(p, 0, 20), 20))
	if ctx.Code.Pos+20 < len(ctx.Code.text) {
		codeSgmnt += boldBlack(">")
	} else {
		codeSgmnt += boldBlack("]")
	}
	return codeSgmnt + " " + cyan(strings.Join(make([]string, ctx.stackLength), "| "))
}

func (ctx *ParseContext) log(message string, opts TraceOptions) {
	if !ctx.parseOptions.SkipLog {
		line := ctx.Code.Line()
		if opts.FilterLine == -1 || line == opts.FilterLine {
			fmt.Printf("%s %s\n", ctx.String(), message)
		}
	}
}

func (ctx *ParseContext) loopStackPush(name string) { ctx.loopStack = append(ctx.loopStack, name) }
func (ctx *ParseContext) loopStackPop()             { ctx.loopStack = ctx.loopStack[:len(ctx.loopStack)-1] }
func (ctx *ParseContext) stackPeek(skip int) *frame { return ctx.stack[ctx.stackLength-1-skip] }
func (ctx *ParseContext) stackPush(x Parser) {
	ctx.stack[ctx.stackLength] = ctx.getFrame(x)
	ctx.stackLength++
}
func (ctx *ParseContext) stackPop() { ctx.stackLength-- }

func (ctx *ParseContext) getFrame(x Parser) *frame {
	id := x.getgnode().id
	pos := ctx.Code.Pos
	posFrames := ctx.frames[pos]
	frame := posFrames[id]
	if frame != nil {
		return frame
	} else {
		posFrames[id] = newFrame(pos, id)
		return posFrames[id]
	}
}

func (ctx *ParseContext) wipeWith(frame_ *frame, makeStash bool) stash {
	// return *stash, or nil if !makeStash
	// default for makeStash was true in coffee
	if frame_.wipemask == nil {
		panic(errors.New("need frame.wipemask to know what to wipe"))
	}
	var stash_ []*frame
	if makeStash {
		stash_ = make([]*frame, ctx.numRules)
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
		posFrames[i] = nil
		stashCount++
	}
	if stash_ != nil {
		return stash_
	} else {
		return nil
	}
}

func (ctx *ParseContext) restoreWith(stash_ stash) {
	for i, frame := range stash_ {
		if frame != nil {
			ctx.frames[frame.pos][i] = frame
		}
	}
}
