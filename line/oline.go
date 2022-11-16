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

type OLine struct {
	args []any
}

type OLineByIndexOrByName struct {
	name  string
	index helpers.NullInt
}

func OEmpty() OLine { return OLine{[]any{}} }

/*
O() is a variadic function which allows a variety of declarations, for example:
- O("EXPR", Rules(....))       // "EXPR" is a rule name
- O("CHOICE _")                // "CHOICE _" is a rule desc because there is no rules() array
- O("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Astnode) Astnode { return new Choice it})
- O("<name>", func(it Astnode, ctx *ParseContext) Astnode { return <...> }, ParseOptions{ SkipLog: true, SkipCache: false }
- O(S(St("{"), R("_"), L("min",E(R("INT"))), R("_"), St(","), R("_"), L("max",E(R("INT"))), R("_"), St("}")))
The last one is a handcompiled rule with which the joeson grammar is initially defined as in ast/handcompiled
*/
func O(a ...any) OLine { return OLine{lineInit(a)} }

func (ol OLine) Args() []any      { return args }
func (ol OLine) IsO() bool        { return true }
func (ol OLine) LineType() string { return "o" }
func (ol OLine) String() string   { return ol.StringIndent(0) }
func (ol OLine) StringIndent(nIndent int) string {
	return helpers.Indent(nIndent) + ol.LineType() +
		" " +
		strings.Join(
			lambda.Map(
				ol.args,
				func(arg any) string {
					switch v := arg.(type) {
					case Astnode:
						return v.ContentString()
					case string:
						return v
					case []Line:
						return BoldBlue("[\n") + strings.Join(
							lambda.Map(v, func(line Line) string { return line.StringIndent(nIndent + 1) }),
							"\n",
						) + "\n" + helpers.Indent(nIndent) + BoldBlue("]")
					case func(Astnode) Astnode:
						return Yellow("𝘧")
					case func(Astnode, *ParseContext) Astnode:
						return Yellow("func (Astnode, *ParseContext) Astnode")
					default:
						return "?" + reflect.TypeOf(v).String() + "?"
					}
				},
			),
			BoldGreen(", "),
		)
}

// note TODO think parentRule could almost simply be GNode. but anyway
func (ol OLine) ToRule(grammar *ast.Grammar, parentRule Astnode, by OLineByIndexOrByName) Astnode {
	rule, attrs := getArgs(ol)
	// figure out the name for this rule
	if nn, isNativeNamed := rule.(NativeNamed); isNativeNamed {
		// in original coffee impl, it is a map with just 1 key
		// in golang impl, see NativeNamed.go
		by.name = nn.Name
		rule = nn.Value
	} else if by.name == "" && by.index.IsSet() && parentRule != nil {
		by.name = parentRule.GetGNode().Name + "[" + strconv.Itoa(by.index) + "]"
	} else if by.name == "" {
		panic("Name undefined for 'o' line")
	}
	rule = getRule(grammar, by.name, rule, attrs)
	rule.GetGNode().Parent = parentRule
	// TODO is the following commented line really useful? I am not sure yet
	// rule.GetGNode().Index = by.index
	return rule
}
