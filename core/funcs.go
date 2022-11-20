package core

import . "grepsuzette/joeson/colors"
import "grepsuzette/joeson/helpers"
import "strconv"
import "strings"

// type ParseFunction func(*ParseContext, Astnode) Astnode
type ParseFunction func(*ParseContext) Astnode

// This file is the dirtiest part of this implementation.
// It is a direct port of joeson.coffee, when it works it can be cleaned.
// For now an almost exact mapping with joeson.coffee is probably good to have.

var _loopStack []string

func _loopStackPop() {
	_loopStack = _loopStack[:len(_loopStack)-1]
}

// TODO rename to LabelOrName
func ShowLabelOrNameIfAny(n Astnode) string {
	if n.GetGNode().IsRule() {
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

func Stack(fparse ParseFunction, x Astnode) ParseFunction {
	return func(ctx *ParseContext) Astnode {
		ctx.StackPush(x)
		if TimeStart != nil {
			TimeStart(x.GetGNode().Name)
		}
		// pos := ctx.Code.pos  // TODO original is `pos = $.code.pos` but seems effectless
		result := fparse(ctx)
		if TimeEnd != nil {
			TimeEnd(x.GetGNode().Name)
		}
		ctx.StackPop()
		return result
	}
}

// Loopify requires the 2nd arg `Astnode`
func Loopify(fparse ParseFunction, x Astnode) ParseFunction {
	return func(ctx *ParseContext) Astnode {
		ctx.logIf(Trace.Stack, Blue("*")+" "+x.ContentString()+" "+Black(strconv.Itoa(ctx.counter)))
		if x.GetGNode().SkipCache {
			result := fparse(ctx)
			ctx.logIf(Trace.Stack, Cyan("`->:")+" "+helpers.Escape(result.ContentString())+" "+Black(helpers.TypeOfToString(result)))
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
				ctx.logIf(Trace.Stack, Cyan("`-hit:")+" "+helpers.Escape(frame.Result.ContentString())+" "+Black(helpers.TypeOfToString(frame.Result)))
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
				ctx.logIf(Trace.Stack, Cyan("`-set:")+" "+helpers.Escape(result.ContentString())+" "+Black(helpers.TypeOfToString(result)))
				return result
			case 2: // recursion detected by subroutine above
				if result == nil {
					ctx.logIf(Trace.Stack, Yellow("`--- loop nil --- "))
					frame.loopStage.Set(0)
					// cacheSet(frame, nil) // unneeded (already nil)
					return result
				} else {
					frame.loopStage.Set(3)
					if Trace.Loop && ((Trace.FilterLine < 0) || ctx.Code.Line() == Trace.FilterLine) {
						line := ctx.Code.Line()
						_loopStack = append(_loopStack, x.GetGNode().Name)
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
						s += " - " + strings.Join(_loopStack, ", ")
						s += " - " + Yellow(helpers.Escape(result.ContentString()))
						s += ": " + Blue(helpers.Escape(ctx.Code.Peek(Peek{BeforeChars: helpers.NewNullInt(10), AfterChars: helpers.NewNullInt(10)})))
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
						ctx.logIf(Trace.Stack, Yellow("|`--- loop iteration ---")+frame.toString())
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
						_loopStackPop()
					}
					ctx.wipeWith(frame, false)
					ctx.restoreWith(bestStash)
					ctx.Code.Pos = bestEndPos
					ctx.logIf(Trace.Stack, Yellow("`--- loop done! --- ")+"best result: "+helpers.Escape(bestResult.ContentString()))
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
			ctx.logIf(Trace.Stack, Yellow("`-base: ")+helpers.Escape(frame.Result.ContentString())+" "+Black(helpers.TypeOfToString(frame.Result)))
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

func PrepareResult(fparse ParseFunction, node Astnode) ParseFunction {
	// this attaches _origin to the provided `node`
	//      is called by wrap which is called by Parse
	return func(ctx *ParseContext) Astnode {
		gn := node.GetGNode()
		ctx.counter++
		result := fparse(ctx) // TODO one problem here is
		// in js there is the fn.call this, $
		// which binds the gnode where it was called
		// to fparse.
		if result != nil {
			// handle labels for standalone nodes
			if gn.Label != "" && gn.Parent != nil && !gn.Parent.HandlesChildLabel() {
				result = NewNativeMap(map[string]Astnode{gn.Label: result})
			}
			start := ctx.StackPeek(0).pos
			end := ctx.Code.Pos
			origin := Origin{
				code: ctx.Code.text,
				start: Cursor{
					line: ctx.Code.PosToLine(start),
					col:  ctx.Code.PosToCol(start),
					pos:  start,
				},
				end: Cursor{
					line: ctx.Code.PosToLine(end),
					col:  ctx.Code.PosToCol(end),
					pos:  end,
				},
			}
			if gn.CbBuilder != nil {
				result.GetGNode()._origin = origin
				// TODO please check this, `@cb.call this, result, $`
				// TODO if there are labels, result
				// must be a NativeMap.
				// If there are none, i suppose probably
				// it must be a NativeArray?
				result = gn.CbBuilder(result, ctx)
			}
			result.GetGNode()._origin = origin // set it again
		}
		return result
	}
}

func Wrap(fparse ParseFunction, node Astnode) ParseFunction {
	wrapped1 := Stack(Loopify(PrepareResult(fparse, node), node), node)
	wrapped2 := PrepareResult(fparse, node)
	gn := node.GetGNode()
	// TODO see if it was not oversimplified....
	return func(ctx *ParseContext) Astnode {
		// TODO because it is interface cmp, triple check it,
		// I suspect it will be wrong
		if node == gn.Rule {
			return wrapped1(ctx)
		} else if gn.Label != "" &&
			(gn.Parent != nil && !gn.Parent.HandlesChildLabel()) ||
			gn.CbBuilder != nil {
			return wrapped2(ctx)
		} else {
			return fparse(ctx)
		}
	}
}
