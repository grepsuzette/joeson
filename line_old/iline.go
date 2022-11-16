package line

import (
	"fmt"
	"grepsuzette/joeson/ast"
	"grepsuzette/joeson/helpers"
	"grepsuzette/joeson/lambda"
	"reflect"
	"strconv"
	"strings"

	. "grepsuzette/joeson/core"

	. "grepsuzette/joeson/colors"
)

type ILine struct {
	args []any
}

func IEmpty() ILine { return ILine{[]any{}} }

/*
I() is a variadic function which allows a variety of declarations, for example:
- I("INT", "/[0-9]+/")
- I("INT", "/[0-9]+/", func(it Astnode) Astnode { return new NativeInt(it) })
- I("INT", "/[0-9]+/", func(it Astnode, ctx *ParseContext) Astnode { return <...> })
- I("INT", "/[0-9]+/", func(it Astnode) Astnode { return <...> }, core.ParseOptions{SkipLog: false, SkipCache: true})
- I("RANGE", O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}"))))
  This one is a handcompiled rule with an O which the joeson grammar is initially defined as in ast/handcompiled
- I("LABEL", C(St('&'), St('@'), R("WORD"))),
  That one is a handcompiled rule that doesn't use an O rule.
*/
func I(a ...any) ILine { return ILine{a} }

func (il ILine) IsO() bool        { return false }
func (il ILine) LineType() string { return "i" }
func (il ILine) String() string   { return il.StringIndent(0) }
func (il ILine) StringIndent(nIndent int) string {
	return helpers.Indent(nIndent) + il.LineType() +
		" " +
		strings.Join(lambda.Map(
			il.args,
			func(arg any) string {
				prefix := (reflect.TypeOf(arg).String() + "=")
				switch v := arg.(type) {
				case string:
					// return Cyan("\"") + (v) + Cyan("\"")
					return v
				case NativeString:
					return BoldRed("NativeString=\"") + BoldRed(v.Str) + BoldRed("\"")
				case OLine:
					return v.String()
				case Astnode:
					return v.ContentString()
				case func(Astnode) Astnode:
					return Yellow("𝘧")
				case func(Astnode, *ParseContext) Astnode:
					return Yellow("func (Astnode, *ParseContext) Astnode {...}")
				default:
					return prefix + "??"
					//Magenta("<unhandled>reflect.TypeOf=") + reflect.TypeOf(arg).String()
				}
			},
		), BoldGreen(", "))
}

// Since ILine is constructed with variadic ...any, `getArgs()` helps get those
// args in a structured way.
//
// The 4th return `Astnode` is only non-nil whenever the ILine second argument
// (called `def`) was an Astnode. In that case, the 2nd return `def` is "".
//
// Return:
// - `name` the name of the rule
// - `def` is a definition for the rule called `name` (e.g. name: "OPER" and def: "[+-]")
//     Can be a string,
//     Can be an OLine call, e.g. `I("RANGE", O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}"))))`
// - `attrs`:
//    - `CbBuilder` is either nil or the inline callback in a grammar
//       definition to build nodes
//    - `SkipLog`
//    - `SkipCache`
// - `astnode`: normally nil. Only set when `il.args[1]` is Astnode or OLine.
//       in which case it is a handcompiled rule and can be returned directly (meaning
//       other returns can be ignored)
// parentRule only use is when the 2nd arg (def) is given as an OLine
func (il ILine) getArgs(parentRule Astnode) (name string, def string, attrs ParseOptions, astnode Astnode) {
	// reminder how an ILine can be declared (i.e. args to iterate through):
	// I(
	//   "INT",
	//   "/[0-9]+/",
	//   func(it Astnode, ctx ...*ParseContext) Astnode { return nil },
	//   ParseOptions{SkipLog: false, SkipCache: true}
	// )
	if len(il.args) < 2 {
		panic("Not enough args for Iline: " + fmt.Sprintf("%v", il.args))
	}
	for i, arg := range il.args {
		if i == 0 {
			name = arg.(string)
		} else if i == 1 {
			switch v := arg.(type) {
			case string:
				def = v
			case OLine:
				// for example,
				// v can be result of `o` line below, and `name` can be "RANGE":
				// i "RANGE":
				//   o S(St('{'), R("_"), L("min",E(R("INT"))), R("_"), St(','), R("_"), L("max",E(R("INT"))), R("_"), St('}'))
				def = ""
				astnode = v.ToRuleWithName(parentRule, name)
				return
			case Astnode:
				astnode = v
				def = ""
				return
			default:
				panic("ILine unexpected type for arg 2 (def)")
			}
		} else {
			if f, ok := arg.(func(Astnode) Astnode); ok {
				attrs.CbBuilder = func(z Astnode, _ *ParseContext) Astnode { return f(z) }
				// TODO forgot OLine, see i "RANGE": o S(St('{'), R("_"), L("min",E(R("INT"))), R("_"), St(','), R("_"), L("max",E(R("INT"))), R("_"), St('}'))
			} else if f, ok := arg.(func(Astnode, *ParseContext) Astnode); ok {
				attrs.CbBuilder = f
			} else if passedAttrs, ok := arg.(ParseOptions); ok {
				if passedAttrs.CbBuilder != nil {
					attrs.CbBuilder = passedAttrs.CbBuilder
				}
				attrs.SkipCache = passedAttrs.SkipCache
				attrs.SkipLog = passedAttrs.SkipLog
			} else {
				fmt.Println("Ignoring arg " + strconv.Itoa(i) + ": " + fmt.Sprintf("%v", arg) + " in ILine " + fmt.Sprintf("%v", il.args))
			}
		}
	}
	return
}

// called by line.funcs.go (NewRankFromLines)
// parentRule only usage is when the 2nd arg (def) is given as an OLine
func (il ILine) ToRule(grammar *ast.Grammar, parentRule Astnode) (name string, rule Astnode) {
	// note: current impl. warrants one rule, in contrast to the
	// orig inal joeson.coffee impl. which returned an {key:val}
	name, def, attrs, astnode := il.getArgs(parentRule)
	if astnode != nil {
		if def != "" {
			panic("logic")
		}
		return name, astnode
	}
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Error in ILine named \"%s\" with def \"%s\":\n%v\n", name, def, e)
			grammar.ParseString(def, attrs) // make it fail for real this time
		}
	}()
	// temporarily halt trace
	oldTrace := Trace
	ctx := NewParseContext(NewCodeStream(def), grammar, attrs)
	// ------------------------
	rule = grammar.Parse(ctx)
	// ------------------------
	Trace = oldTrace
	rule.GetGNode().Name = name
	rule.GetGNode().CbBuilder = attrs.CbBuilder
	rule.GetGNode().SkipCache = attrs.SkipCache
	rule.GetGNode().SkipLog = attrs.SkipLog
	rule.GetGNode().Debug = attrs.Debug
	return
}
