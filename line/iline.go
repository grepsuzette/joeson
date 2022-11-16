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
func I(a ...any) ILine { return ILine{lineInit(a)} }

func (il ILine) Args() []any      { return il.args }
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

// note TODO think parentRule could almost simply be GNode. but anyway
func (il ILine) ToRules(grammar *ast.Grammar, parentRule Astnode) NativeMap {
	rule, attrs := getArgs(il)
	rules := NewEmptyNativeMap()
	// (in coffee) for an ILine, rule is an object of {"NAME":rule}
	// (in golang) we just use NativeNamed, which is a {Name string, Value Astnode}
	switch v := rule.(type) {
	case NativeNamed:
		rules.Set(v.Name, getRule(grammar, v.Name, v.Value, parentRule, attrs))
	case NativeMap:
		panic("Assume unrreachable for now")
		for _, k := range v.Keys() {
			rules.Set(k, getRule(grammar, k, v.Get(k), parentRule, attrs))
		}
	default:
		panic("wjiefwe")
	}
	return rules
}
