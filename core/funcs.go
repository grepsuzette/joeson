package core

import (
	. "grepsuzette/joeson/colors"
	"grepsuzette/joeson/helpers"
	"strconv"
	"strings"
)

// in coffeescript/js, the 2nd argument (`Astnode`) doesn't exist,
// instead .bind(this) is used
type ParseFunction2 func(*ParseContext, Astnode) Astnode
type ParseFunction func(*ParseContext) Astnode

// This file is the dirtiest part of this implementation.
// It is a direct port of joeson.coffee, when it works it can be cleaned.
// For now an almost exact mapping with joeson.coffee is probably good to have.

func LabelOrName(n Astnode) string {
	if IsRule(n) {
		return Red(n.GetGNode().Name + ": ")
	} else if n.GetGNode().Label != "" {
		return Cyan(n.GetGNode().Label + ":")
	}
	return ""
}

// in joeson.coffee those functions were originally declared as
// class method to GNode and had a $ prefix:
// @GNode
//   @$stack = (fn) -> ($) -> Astnode
//   @$loopify = (fn) -> ($) -> Astnode
//   @$prepareResult = (fn) -> ($) -> Astnode
//   @$wrap = (fn) -> Astnode
//
// Here they are called _stack, _loopify, _prepareResult and _wrap

func stack(fparse ParseFunction, x Astnode) ParseFunction {
	return func(ctx *ParseContext) Astnode {
		ctx.StackPush(x)
		if TimeStart != nil {
			TimeStart(x.GetGNode().Name)
		}
		result := fparse(ctx)
		if TimeEnd != nil {
			TimeEnd(x.GetGNode().Name)
		}
		ctx.StackPop()
		return result
	}
}

func loopify(fparse ParseFunction, x Astnode) ParseFunction {
	return func(ctx *ParseContext) Astnode {
		log := func(s string) {}
		if Trace.Stack {
			log = func(s string) { ctx.log(s) }
		}
		log(Blue("*") + " " + x.ContentString() + " " + BoldBlack(strconv.Itoa(ctx.counter)))
		if x.GetGNode().SkipCache {
			result := fparse(ctx)
			log(Cyan("`->:") + " " + helpers.Escape(result.ContentString()) + " " + BoldBlack(helpers.TypeOfToString(result)))
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
				if frame.Result != nil {
					log(Cyan("`-hit: ") + helpers.Escape(frame.Result.ContentString()) + " " + BoldBlack(helpers.TypeOfToString(frame.Result)))
				} else {
					log(Cyan("`-hit: nil"))
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
				s := Cyan("`-set: ")
				if result == nil {
					s += "nil"
				} else {
					s += helpers.Escape(result.ContentString())
					s += " "
					s += BoldBlack(helpers.TypeOfToString(result))
				}
				log(s)
				return result
			case 2: // recursion detected by subroutine above
				if result == nil {
					log(Yellow("`--- loop nil --- "))
					frame.loopStage.Set(0)
					// cacheSet(frame, nil) // unneeded (already nil)
					return result
				} else {
					frame.loopStage.Set(3)
					if Trace.Loop && ((Trace.FilterLine < 0) || ctx.Code.Line() == Trace.FilterLine) {
						line := ctx.Code.Line()
						// _loopStack = append(_loopStack, x.GetGNode().Name)
						ctx.loopStackPush(x.GetGNode().Name)
						var paintInColor func(string) string = nil
						switch line % 6 {
						case 0:
							paintInColor = Blue
						case 1:
							paintInColor = Cyan
						case 2:
							paintInColor = White
						case 3:
							paintInColor = Yellow
						case 4:
							paintInColor = Red
						case 5:
							paintInColor = Magenta
						}
						s := ""
						s += paintInColor("@" + strconv.Itoa(line))
						s += "\t"
						for _, frame := range ctx.stack[0:ctx.stackLength] {
							s += Red(strconv.Itoa(frame.id))
						}
						// s += " - " + strings.Join(_loopStack, ", ")
						s += " - " + strings.Join(ctx.loopStack, ", ")
						s += " - " + Yellow(helpers.Escape(result.ContentString()))
						s += ": " + Blue(helpers.Escape(ctx.Code.Peek(NewPeek().BeforeChars(10).AfterChars(10))))
					}
					if TimeStart != nil {
						TimeStart("loopiteration")
					}
					var bestStash *stash = nil
					var bestEndPos int = 0
					var bestResult Astnode = nil
					for result != nil {
						if frame.wipemask == nil {
							panic("where's my wipemask")
						}
						bestStash = ctx.wipeWith(frame, true)
						bestResult = result
						bestEndPos = ctx.Code.Pos
						frame.cacheSet(bestResult, bestEndPos)
						log(Yellow("|`--- loop iteration ---") + frame.toString())
						ctx.Code.Pos = startPos
						result = fparse(ctx)
						if ctx.Code.Pos <= bestEndPos {
							break
						}
					}
					if TimeEnd != nil {
						TimeEnd("loopiteration")
					}
					if Trace.Loop {
						ctx.loopStackPop()
					}
					ctx.wipeWith(frame, false)
					ctx.restoreWith(bestStash)
					ctx.Code.Pos = bestEndPos
					log(Yellow("`--- loop done! --- ") + "best result: " + helpers.Escape(bestResult.ContentString()))
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
			log(Yellow("`-base: ") + helpers.Escape(frame.Result.ContentString()) + " " + BoldBlack(helpers.TypeOfToString(frame.Result)))
			if frame.wipemask == nil {
				frame.wipemask = make([]bool, ctx.grammar.CountRules())
				for i := ctx.stackLength - 2; i >= 0; i-- {
					i_frame := ctx.stack[i]
					if i_frame.pos > startPos {
						panic("assert failed: i_frame.pos > startPos")
					}
					if i_frame.pos < startPos || i_frame.id == x.GetGNode().Id {
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
// - set GNode._origin
// - call GNode.CbBuilder(result, ctx, caller), if CbBuilder != nil
func prepareResult(fparse2 ParseFunction2, caller Astnode) ParseFunction {
	return func(ctx *ParseContext) Astnode {
		ctx.counter++
		result := fparse2(ctx, caller) // .call() is used in js
		if result != nil {
			// handle labels for standalone nodes
			gn := caller.GetGNode()
			if gn.Label != "" && gn.Parent != nil && !gn.Parent.HandlesChildLabel() {
				result = NewNativeMap(map[string]Astnode{gn.Label: result})
			}
			start := ctx.StackPeek(0).pos
			end := ctx.Code.Pos
			origin := Origin{
				code: ctx.Code.text,
				start: Cursor{
					line: ctx.Code.PosToLine(start),
					col:  ctx.Code.PosToLine(start), // TODO there's likely a bug in original here, how about PosToCol, but never mind
					pos:  start,
				},
				end: Cursor{
					line: ctx.Code.PosToLine(end),
					col:  ctx.Code.PosToLine(end), // TODO there's likely a bug in original here, but never mind
					pos:  end,
				},
			}
			if gn.CbBuilder != nil {
				// Native* don't have a GNode.
				// It also doesn't make sense to store an origin
				// for them anyway
				if result.GetGNode() != nil {
					result.GetGNode()._origin = origin
				}
				// in js, it is bounded to this (`caller`)
				result = gn.CbBuilder(result, ctx, caller)
			}
			// TODO remove this next line after proven it works without
			if result.GetGNode() != nil {
				result.GetGNode()._origin = origin // set it again. though original impl. has this, it is most definitely not required here
			}
		}
		return result
	}
}

func Wrap(fparse2 ParseFunction2, node Astnode) ParseFunction {
	wrapped1 := stack(loopify(prepareResult(fparse2, node), node), node)
	wrapped2 := prepareResult(fparse2, node)
	gn := node.GetGNode()
	return func(ctx *ParseContext) Astnode {
		if IsRule(node) {
			return wrapped1(ctx)
		} else if gn.Label != "" &&
			(gn.Parent != nil && !gn.Parent.HandlesChildLabel()) ||
			gn.CbBuilder != nil {
			return wrapped2(ctx)
		} else {
			return fparse2(ctx, node)
		}
	}

}
