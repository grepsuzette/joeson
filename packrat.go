package joeson

import (
	"github.com/grepsuzette/joeson/helpers"
	"strconv"
)

// refer to the original joeson.coffee

type frame struct {
	Result    Ast
	endPos    helpers.NilableInt // can be left undefined
	loopStage helpers.NilableInt // can be left undefined
	wipemask  []bool             // len = ctx.grammar.numRules
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

// These callback types derive from the direct port from the
// coffeescript impl.  In coffeescript/js, the 2nd
// argument (`Ast`) doesn't exist, instead .bind(this) was used

type parseFun2 func(*ParseContext, Parser) Ast
type parseFun func(*ParseContext) Ast

// debugging variables and callbacks
// to rewrite soon in a more idiomatic way
var TimeStart func(name string) = nil
var TimeEnd func(name string) = nil

// in joeson.coffee those functions were originally declared as
// class method to GNode and had a $ prefix:
// @GNode
//   @$stack = (fn) -> ($) -> Ast
//   @$loopify = (fn) -> ($) -> Ast
//   @$prepareResult = (fn) -> ($) -> Ast
//   @$wrap = (fn) -> Ast
//
// Here they are called stack, loopify, prepareResult and Wrap:
//
// - func stack(fparse parseFun, x Ast) parseFun
// - func loopify(fparse parseFun, x Ast) parseFun
// - func prepareResult(fparse2 parseFun2, caller Ast) parseFun
// - func Wrap(fparse2 parseFun2, node Ast) parseFun
//     notice this line in Wrap:  wrapped1 := stack(loopify(prepareResult(fparse2, node), node), node)

func stack(fparse parseFun, x Parser) parseFun {
	return func(ctx *ParseContext) Ast {
		ctx.stackPush(x)
		if TimeStart != nil {
			TimeStart(x.Name())
		}
		result := fparse(ctx)
		if TimeEnd != nil {
			TimeEnd(x.Name())
		}
		ctx.stackPop()
		return result
	}
}

func loopify(fparse parseFun, x Parser) parseFun {
	return func(ctx *ParseContext) Ast {
		opts := ctx.TraceOptions
		if opts.Stack {
			ctx.log(blue("*")+" "+String(x)+" "+boldBlack(strconv.Itoa(ctx.Counter)), opts)
		}
		if x.GetGNode().SkipCache {
			result := fparse(ctx)
			if opts.Stack {
				ctx.log(cyan("`->:")+" "+helpers.Escape(result.ContentString())+" "+boldBlack(helpers.TypeOfToString(result)), opts)
			}
			return result
		}
		frame := ctx.getFrame(x)
		startPos := ctx.Code.Pos
		if !frame.loopStage.IsSet {
			frame.loopStage.Set(0)
		}
		switch frame.loopStage.Int {
		case 0: // non-recursive (so far)
			// The only time a cache hit will simply return is when loopStage is 0
			if frame.endPos.IsSet {
				if opts.Stack {
					if frame.Result != nil {
						s := ""
						s += helpers.Escape(frame.Result.ContentString())
						s += " "
						s += cyan(helpers.TypeOfToString(frame.Result))
						ctx.log(cyan("`-hit:")+" "+s, opts)
					} else {
						ctx.log(cyan("`-hit:")+" nil", opts)
					}
				}
				ctx.Code.Pos = frame.endPos.Int
				return frame.Result
			}
			frame.loopStage.Set(1)
			frame.cacheSet(nil, -1)
			result := fparse(ctx)
			switch frame.loopStage.Int {
			case 1: // non-recursive (i.e. done)
				frame.loopStage.Set(0)
				frame.cacheSet(result, ctx.Code.Pos)
				if opts.Stack {
					s := cyan("`-set:") + " "
					if result == nil {
						s += "nil"
					} else {
						s += helpers.Escape(result.ContentString())
						s += " "
						s += cyan(helpers.TypeOfToString(result))
					}
					ctx.log(s, opts)
				}
				return result
			case 2: // recursion detected by subroutine above
				if result == nil {
					if opts.Stack {
						ctx.log(yellow("`--- loop nil --- "), opts)
					}
					frame.loopStage.Set(0)
					// cacheSet(frame, nil) // unneeded (already nil)
					return result
				} else {
					frame.loopStage.Set(3)
					if opts.Loop && ((opts.FilterLine < 0) || ctx.Code.Line() == opts.FilterLine) {
						ctx.loopStackPush(x.Name())
						// if false {
						//  line := ctx.Code.Line()
						// 	var paintInColor func(string) string = nil
						// 	switch line % 6 {
						// 	case 0:
						// 		paintInColor = blue
						// 	case 1:
						// 		paintInColor = cyan
						// 	case 2:
						// 		paintInColor = white
						// 	case 3:
						// 		paintInColor = yellow
						// 	case 4:
						// 		paintInColor = red
						// 	case 5:
						// 		paintInColor = magenta
						// 	}
						// 	s := ""
						// 	s += paintInColor("@" + strconv.Itoa(line))
						// 	s += "\t"
						// 	for _, frame := range ctx.stack[0:ctx.stackLength] {
						// 		s += red(strconv.Itoa(frame.id))
						// 	}
						// 	s += " - " + strings.Join(ctx.loopStack, ", ")
						// 	s += " - " + yellow(helpers.Escape(result.ContentString()))
						// 	s += ": " + blue(helpers.Escape(ctx.Code.Peek(NewPeek().BeforeChars(10).AfterChars(10))))
						// 	fmt.Println(s) // also this way in original joeson.coffee
						// }
					}
					if TimeStart != nil {
						TimeStart("loopiteration")
					}
					var bestStash *stash = nil
					var bestEndPos int = 0
					var bestResult Ast = nil
					for result != nil {
						if frame.wipemask == nil {
							panic("where's my wipemask")
						}
						bestStash = ctx.wipeWith(frame, true)
						bestResult = result
						bestEndPos = ctx.Code.Pos
						frame.cacheSet(bestResult, bestEndPos)
						if opts.Stack {
							ctx.log(yellow("|`--- loop iteration ---")+frame.toString(), opts)
						}
						ctx.Code.Pos = startPos
						result = fparse(ctx)
						if ctx.Code.Pos <= bestEndPos {
							break
						}
					}
					if TimeEnd != nil {
						TimeEnd("loopiteration")
					}
					if opts.Loop {
						ctx.loopStackPop()
					}
					ctx.wipeWith(frame, false)
					ctx.restoreWith(bestStash)
					ctx.Code.Pos = bestEndPos
					if opts.Stack {
						ctx.log(yellow("`--- loop done! --- ")+"best result: "+helpers.Escape(bestResult.ContentString()), opts)
					}
					// Step 4: return best result, which will get cached
					frame.loopStage.Set(0)
					return bestResult
				}
			default:
				panic("Unexpected stage " + strconv.Itoa(frame.loopStage.Int))
			}
		case 1, 2, 3:
			if frame.loopStage.Int == 1 {
				frame.loopStage.Set(2) // recursion detected
			}
			if TimeStart != nil {
				TimeStart("wipemask")
			}
			// Step 1: Collect wipemask so we can wipe the frames later.
			if opts.Stack {
				ctx.log(yellow("`-base: ")+helpers.Escape(frame.Result.ContentString())+" "+boldBlack(helpers.TypeOfToString(frame.Result)), opts)
			}
			if frame.wipemask == nil {
				frame.wipemask = make([]bool, ctx.numRules)
				for i := ctx.stackLength - 2; i >= 0; i-- {
					i_frame := ctx.stack[i]
					if i_frame.pos > startPos {
						panic("assert failed: i_frame.pos > startPos")
					}
					if i_frame.pos < startPos || i_frame.id == x.GetGNode().id {
						break
					}
					frame.wipemask[i_frame.id] = true
				}
				if TimeEnd != nil {
					TimeEnd("wipemask")
				}
				// Step 2: Return whatever was cacheSet.
				if frame.endPos.IsSet {
					ctx.Code.Pos = frame.endPos.Int
				}
				return frame.Result
			}
		default:
			panic("Unexpected stage " + strconv.Itoa(frame.loopStage.Int) + " (B)")
		}
		return nil
	}
}

// prepares the following postparsing operations:
// - increment ctx.counter (used for debugging and to prevent infinite recursion)
// - handle labels for standalone nodes
// - set GNode.Origin
// - call GNode.CbBuilder(result, ctx, caller), if CbBuilder != nil
func prepareResult(fparse2 parseFun2, caller Parser) parseFun {
	return func(ctx *ParseContext) Ast {
		ctx.Counter++
		result := fparse2(ctx, caller)
		if result != nil {
			// handle labels for standalone nodes
			gn := caller.GetGNode()
			if gn.label != "" && gn.parent != nil && !gn.parent.HandlesChildLabel() {
				result = NewNativeMap(map[string]Ast{gn.label: result})
			}
			if gn.CbBuilder != nil {
				switch x := result.(type) {
				case Parser:
					start := ctx.stackPeek(0).pos
					end := ctx.Code.Pos
					origin := Origin{
						code:  ctx.Code.text,
						start: start,
						end:   end,
					}
					x.GetGNode().origin = origin
				default:
				}
				result = gn.CbBuilder(result, ctx, caller)
			}
		}
		return result
	}
}

func Wrap(fparse2 parseFun2, node Parser) parseFun {
	wrapped1 := stack(loopify(prepareResult(fparse2, node), node), node)
	wrapped2 := prepareResult(fparse2, node)
	gn := node.GetGNode()
	return func(ctx *ParseContext) Ast {
		if IsRule(node) {
			return wrapped1(ctx)
		} else if gn.label != "" &&
			(gn.parent != nil && !gn.parent.HandlesChildLabel()) ||
			gn.CbBuilder != nil {
			return wrapped2(ctx)
		} else {
			return fparse2(ctx, node)
		}
	}

}
