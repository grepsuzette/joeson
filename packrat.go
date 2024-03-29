package joeson

import (
	"strconv"

	"github.com/grepsuzette/joeson/helpers"
)

// packrat parsing algorithm. From original joeson.coffee

type (
	frame struct {
		result    Ast
		endpos    int    // -1 means not set
		loopstage int    // -1 means not set TODO enum
		wipemask  []bool // len = ctx.grammar.numRules
		pos       int
		id        int
		param     Ast // used in ref.go or joeson.coffee:536
	}
	stash []*frame
)

func (f frame) toString() string {
	return "N/A frame.toString"
}

func (fr *frame) cacheSet(result Ast, endpos int) {
	fr.result = result
	if endpos < 0 {
		fr.endpos = -1
	} else {
		fr.endpos = endpos
	}
}

func newFrame(pos int, id int) *frame {
	return &frame{
		result:    nil,
		pos:       pos,
		id:        id,
		loopstage: -1,
		endpos:    -1,
		wipemask:  nil,
		param:     nil,
	}
}

// These callback types derive from the direct port from the
// coffeescript impl.  In coffeescript/js, the 2nd
// argument (`Ast`) doesn't exist, instead .bind(this) was used

type (
	parseFunc2 func(*ParseContext, Parser) Ast
	parseFunc  func(*ParseContext) Ast
)

// in joeson.coffee those functions were originally declared as
// class method to GNode and had a $ prefix:
// @GNode
//   @$stack = (fn) -> ($) -> Ast
//   @$loopify = (fn) -> ($) -> Ast
//   @$prepareResult = (fn) -> ($) -> Ast
//   @$wrap = (fn) -> Ast
//
// Here they are called stack, loopify, prepareResult and wrap:
//
// - func stack(fparse parseFun, x Ast) parseFun
// - func loopify(fparse parseFun, x Ast) parseFun
// - func prepareResult(fparse2 parseFun2, caller Ast) parseFun
// - func wrap(fparse2 parseFun2, node Ast) parseFun
//     notice this line in wrap:  wrapped1 := stack(loopify(prepareResult(fparse2, node), node), node)

func stack(fparse parseFunc, x Parser) parseFunc {
	return func(ctx *ParseContext) Ast {
		ctx.stackPush(x)
		result := fparse(ctx)
		ctx.stackPop()
		return result
	}
}

func loopify(fparse parseFunc, x Parser) parseFunc {
	return func(ctx *ParseContext) Ast {
		opts := ctx.TraceOptions
		showLog := opts.Stack || x.getRule().parseOptions.debug
		if showLog {
			ctx.log(Blue("*")+" "+String(x)+" "+BoldBlack(strconv.Itoa(ctx.Counter)), opts)
		}
		if x.getRule().skipCache {
			result := fparse(ctx)
			if showLog {
				ctx.log(Cyan("`->:")+" "+helpers.Escape(result.String())+" "+BoldBlack(helpers.TypeOfToString(result)), opts)
			}
			return result
		}
		frame := ctx.getFrame(x)
		startPos := ctx.Code.Pos()
		if frame.loopstage < 0 {
			frame.loopstage = 0
		}
		switch frame.loopstage {
		case 0: // non-recursive (so far)
			// The only time a cache hit will simply return is when loopStage is 0
			if frame.endpos >= 0 {
				if showLog {
					if frame.result != nil {
						s := ""
						s += helpers.Escape(frame.result.String())
						s += " "
						s += Cyan(helpers.TypeOfToString(frame.result))
						ctx.log(Cyan("`-hit:")+" "+s, opts)
					} else {
						ctx.log(Cyan("`-hit:")+" nil", opts)
					}
				}
				ctx.Code.SetPos(frame.endpos)
				return frame.result
			}
			frame.loopstage = 1
			frame.cacheSet(nil, -1)
			result := fparse(ctx)
			switch frame.loopstage {
			case 1: // non-recursive (i.e. done)
				frame.loopstage = 0
				frame.cacheSet(result, ctx.Code.Pos())
				if showLog {
					s := Cyan("`-set:") + " "
					if result == nil {
						s += "nil"
					} else {
						s += helpers.Escape(result.String())
						s += " "
						s += Cyan(helpers.TypeOfToString(result))
					}
					ctx.log(s, opts)
				}
				return result
			case 2: // recursion detected by subroutine above
				if result == nil {
					if showLog {
						ctx.log(Yellow("`--- loop nil --- "), opts)
					}
					frame.loopstage = 0
					// cacheSet(frame, nil) // unneeded (already nil)
					return result
				} else {
					frame.loopstage = 3
					if opts.Loop && ((opts.FilterLine < 0) || ctx.Code.Line() == opts.FilterLine) {
						ctx.loopStackPush(x.getRule().name)
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
						// 	s += " - " + yellow(helpers.Escape(result.String()))
						// 	s += ": " + blue(helpers.Escape(ctx.Code.Peek(NewPeek().BeforeChars(10).AfterChars(10))))
						// 	fmt.Println(s) // also this way in original joeson.coffee
						// }
					}
					var bestStash stash = nil
					var bestEndPos int = 0
					var bestResult Ast = nil
					for result != nil {
						if frame.wipemask == nil {
							panic("where's my wipemask")
						}
						bestStash = ctx.wipeWith(frame, true)
						bestResult = result
						bestEndPos = ctx.Code.Pos()
						frame.cacheSet(bestResult, bestEndPos)
						if showLog {
							ctx.log(Yellow("|`--- loop iteration ---")+frame.toString(), opts)
						}
						ctx.Code.SetPos(startPos)
						result = fparse(ctx)
						if ctx.Code.Pos() <= bestEndPos {
							break
						}
					}
					if opts.Loop {
						ctx.loopStackPop()
					}
					ctx.wipeWith(frame, false)
					ctx.restoreWith(bestStash)
					ctx.Code.SetPos(bestEndPos)
					if showLog {
						ctx.log(Yellow("`--- loop done! --- ")+"best result: "+helpers.Escape(bestResult.String()), opts)
					}
					// Step 4: return best result, which will get cached
					frame.loopstage = 0
					return bestResult
				}
			default:
				panic("Unexpected stage " + strconv.Itoa(frame.loopstage))
			}
		case 1, 2, 3:
			if frame.loopstage == 1 {
				frame.loopstage = 2 // recursion detected
				// ctx.log("left Recursion detected", opts)
			}
			// Step 1: Collect wipemask so we can wipe the frames later.
			if showLog {
				cs := "nil"
				if frame.result != nil {
					cs = frame.result.String()
				}
				ctx.log(Yellow("`-base: ")+
					helpers.Escape(cs)+" "+BoldBlack(helpers.TypeOfToString(frame.result)), opts)
			}
			if frame.wipemask == nil {
				frame.wipemask = make([]bool, ctx.numRules)
			}
			for i := ctx.stackLength - 2; i >= 0; i-- {
				i_frame := ctx.stack[i]
				if i_frame.pos > startPos {
					panic("assert failed: i_frame.pos > startPos")
				}
				if i_frame.pos < startPos || i_frame.id == x.getRule().id {
					break
				}
				frame.wipemask[i_frame.id] = true
			}
			// Step 2: Return whatever was cacheSet.
			if frame.endpos >= 0 {
				ctx.Code.SetPos(frame.endpos)
			}
			return frame.result
		default:
			panic("Unexpected stage " + strconv.Itoa(frame.loopstage) + " (B)")
		}
		return nil
	}
}

// prepares the following postparsing operations:
// - increment ctx.counter (used for debugging and to prevent infinite recursion)
// - handle labels for standalone nodes
// - call SetLine
// - call rule.cb(result, ctx, caller), if cb != nil
func prepareResult(fparse2 parseFunc2, caller Parser) parseFunc {
	return func(ctx *ParseContext) Ast {
		ctx.Counter++
		result := fparse2(ctx, caller)
		if result != nil {
			// handle labels for standalone nodes
			rule := caller.getRule()
			if rule.label != "" && rule.parent != nil && !rule.parent.handlesChildLabel() {
				result = NewNativeMap(map[string]Ast{rule.label: result})
			}
			if rule.cb != nil {
				result = rule.cb(result, ctx, caller)
				if result == nil {
					return nil
				}
			}
			// set origin
			// This only is for debugging.
			// TODO double check original coffee implementation,
			//      not sure if entirely correct.
			result.SetOrigin(Origin{
				Code:     ctx.Code.Code(),
				Start:    ctx.stackPeek(0).pos,
				End:      ctx.Code.Pos(),
				Line:     ctx.Code.Line(),
				RuleName: caller.getRule().parser.getRule().name,
			})
		}
		return result
	}
}

func wrap(fparse2 parseFunc2, node Parser) parseFunc {
	// wrapped1 := stack(loopify(prepareResult(fparse2, node), node), node)
	// wrapped2 := prepareResult(fparse2, node)
	rule := node.getRule()
	return func(ctx *ParseContext) Ast {
		if IsRule(node) {
			// wrapped1
			return stack(loopify(prepareResult(fparse2, node), node), node)(ctx)
		} else if rule.label != "" &&
			(rule.parent != nil && !rule.parent.handlesChildLabel()) ||
			rule.cb != nil {
			// wrapped2
			return prepareResult(fparse2, node)(ctx)
		} else {
			return fparse2(ctx, node)
		}
	}
}
