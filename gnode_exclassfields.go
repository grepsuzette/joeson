package main

// GNode class functions and fields of joeson.coffee have been put here

optionKeys := []string{"skipLog", "skipCache", "cb"}

var _loopStack []string // trace stack TODO this comes directly from joeson.coffee but should clean it
func _loopStackPop() {
	_loopStack = _loopStack[:len(_loopStack)-1]
}

type ParseFunc func(*ParseContext, astnode) astnode

// in joeson.coffee those functions were originally declared as
// class method to GNode and had a $ prefix:
// @GNode
//   @$stack = (fn) -> ($) -> astnode
//   @$loopify = (fn) -> ($) -> astnode
//   @$prepareResult = (fn) -> ($) -> astnode
//   @$wrap = (fn) -> astnode
//
// Here they are called _stack, _loopify, _prepareResult and _wrap

func _stack(fn ParseFunc, x astnode) ParseFunc {
	return func(_, _) astnode {
		ctx.stackPush(x) // <- WTF the argument must be the class... 
		if timeStart != nil {
			timeStart(x.GetGNode().name)
		}
		// pos := ctx.code.pos  // TODO original is `pos = $.code.pos` but seems effectless
		result := fn(x, ctx)
		if timeEnd != nil {
			timeEnd(x.GetGNode().name)
		}
		ctx.stackPop() // TODO I don't understand why original has this in `$.stackPop this`. Seems effectless
		return result
	}
}

func _loopify(fn ParseFunc, node astnode) ParseFunc {
	return func(ctx *ParseContext, gn astnode) astnode {
		ctx.logIf(trace.stack, Blue("*")+" "+gn.toString()+" "+Black(strconv.Itoa(ctx.counter)))
		if gn.skipCache {
			result := fn(gn, ctx)
			ctx.logIf(trace.stack, Cyan("`->:")+" "+helpers.Escape(result.toString())+" "+Black(helpers.TypeOfToString(result)))
			return result
		}
		frame := ctx.getFrame(gn.Node)
		startPos := ctx.code.pos
		if !frame.loopStage.IsSet {
			frame.loopStage.Set(0)
		}
		switch frame.loopStage.Int {
		case 0: // non-recursive (so far)
			// The only time a cache hit will simply return is when loopStage is 0
			if frame.endPos.IsSet {
				ctx.logIf(trace.stack, Cyan("`-hit:")+" "+helpers.Escape(frame.result.toString())+" "+Black(helpers.TypeOfToString(frame.result)))
				ctx.code.pos = frame.endPos.Int
				return frame.result
			}
			frame.loopStage.Set(1)
			frame.cacheSet(nil, -1)
			result := fn(gn, ctx)
			switch frame.loopStage.Int {
			case 1: // non-recursive (i.e. done)
				frame.loopStage.Set(0)
				frame.cacheSet(result, ctx.code.pos)
				ctx.logIf(trace.stack, Cyan("`-set:")+" "+helpers.Escape(result.toString())+" "+Black(helpers.TypeOfToString(result)))
				return result
			case 2: // recursion detected by subroutine above
				if result == nil {
					ctx.logIf(trace.stack, Yellow("`--- loop nil --- "))
					frame.loopStage.Set(0)
					// cacheSet(frame, nil) // unneeded (already nil)
					return result
				} else {
					frame.loopStage.Set(3)
					if trace.loop && ((trace.filterLine < 0) || ctx.code.line() == trace.filterLine) {
						line := ctx.code.line()
						_loopStack = append(_loopStack, gn.name)
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
						s += " - " + Yellow(helpers.Escape(result.toString()))
						s += ": " + Blue(helpers.Escape(ctx.code.peek(Peek{beforeChars: helpers.NewNullInt(10), afterChars: helpers.NewNullInt(10)})))
					}
					if timeStart != nil {
						timeStart("loopiteration")
					}
					var bestStash *stash = nil
					var bestEndPos int = 0
					var bestResult *Result = nil
					for result != nil {
						if frame.wipemask == nil {
							panic("where's my wipemask")
						}
						bestStash = ctx.wipeWith(frame, true)
						bestResult = result
						bestEndPos = ctx.code.pos
						frame.cacheSet(bestResult, bestEndPos)
						ctx.logIf(trace.stack, Yellow("|`--- loop iteration ---")+frame.toString())
						ctx.code.pos = startPos
						result = fn(gn, ctx)
						if ctx.code.pos <= bestEndPos {
							break
						}
					}
					if timeEnd != nil {
						timeEnd("loopiteration")
					}
					if trace.loop {
						_loopStackPop()
					}
					ctx.wipeWith(frame, false)
					ctx.restoreWith(bestStash)
					ctx.code.pos = bestEndPos
					ctx.logIf(trace.stack, Yellow("`--- loop done! --- ")+"best result: "+helpers.Escape(bestResult.toString()))
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
			if timeStart != nil {
				timeStart("wipemask")
			}
			// Step 1: Collect wipemask so we can wipe the frames later.
			ctx.logIf(trace.stack, Yellow("`-base: ")+helpers.Escape(frame.result.toString())+" "+Black(helpers.TypeOfToString(frame.result)))
			if frame.wipemask == nil {
				frame.wipemask = make([]bool, ctx.grammar.numRules)
				for i := ctx.stackLength - 2; i >= 0; i-- {
					i_frame := ctx.stack[i]
					if i_frame.pos > startPos {
						panic("assert failed: i_frame.pos > startPos")
					}
					if i_frame.pos < startPos || i_frame.id == gn.id {
						break
					}
					frame.wipemask[i_frame.id] = true
				}
				if timeEnd != nil {
					timeEnd("wipemask")
				}
				// Step 2: Return whatever was cacheSet.
				if frame.endPos.IsSet {
					ctx.code.pos = frame.endPos.Int
				}
				return frame.result
			}
		default:
			panic("Unexpected stage " + strconv.Itoa(frame.loopStage.Int) + " (B)")
		}
		return nil
	}
}

func _prepareResult(fn ParseFunc, node astnode) ParseFunc {
	// this attaches _origin to the provided `node`
	// it's called by wrap which is called by Parse
	return func(_ *ParseContext, _ astnode) astnode {
		gn := node.GetGNode()
		ctx.counter++
		result := fn(gn, ctx)
		if result != nil {
			// handle labels for standalone nodes
			if gn.label != "" && gn.parent != nil && !gn.parent.HandlesChildLabel() {
				result = NewNativeMap(map[string]{gn.label: result})
			}
			start := ctx.stackPeek(0).pos
			end := ctx.code.pos
			origin := Origin{
				code: ctx.code.text,
				start: SomePos{
					line: ctx.code.posToLine(start),
					col:  ctx.code.posToCol(start),
					pos:  start,
				},
				end: SomePos{
					line: ctx.code.posToLine(end),
					col:  ctx.code.posToCol(end),
					pos:  end,
				},
			}
			if gn.cb != nil {
				if result != nil {
					result.GNode._origin = &origin
				}
				// TODO jae please check this,
				// `@cb.call this, result, $`
				result = gn.cb(gn, result, ctx)
				//              ^ gn or node?? TODO
			}
			if result != nil { // set it again
				result.GNode.origin = &origin
			}
		}
		return result
	}
}

//      v-- ?is it     func (gn GNode)  or is it?     ------------v
func (gn GNode) _wrap(fn ParseFunc/*, node astnode*/) func(*ParseContext) astnode {
	wrapped1 := _stack(_loopify(_prepareResult(fn, node), node), node)
	wrapped2 := _prepareResult(fn, node)
	// TODO how about the parse below? is it correct?
	var parse func(*ParseContext) astnode = nil
	return func(ctx *ParseContext, _ astnode) astnode {
		if &gn == gn.rule {
			gn.parse = func(ctx *ParseContext) astnode {
				return wrapped1(gn, ctx)
			}
		} else if gn.label != "" && gn.parent != nil && !gn.parent.HandlesChildLabel() || gn.cb != nil {
			gn.parse = func(ctx *ParseContext) astnode {
				return wrapped2(gn, ctx)
			}
		} else {
			gn.parse = func(ctx *ParseContext) astnode {
				return fn(gn, ctx)
			}
		}
		parse = gn.parse
		return parse(ctx)
	}
}

