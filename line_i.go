package joeson

import (
	"regexp"

	"github.com/grepsuzette/joeson/helpers"
)

// ILine represents an inner rule in a rank
type ILine struct {
	name    string // ILine, as terminal elements, are always named
	content Line
	attrs   ParseOptions
}

/*
They have a name (always), a content, an optional parse callback and an optional ParseOptions object.
- I(Named("INT", "/[0-9]+/")
- I(Named("INT", "/[0-9]+/"), func(it Ast) Ast { return new NativeInt(it) })
- I(Named("INT", "/[0-9]+/"), func(it Ast, ctx *ParseContext) Ast { return <...> })
- I(Named("INT", "/[0-9]+/"), func(it Ast) Ast { return <...> }, ParseOptions{SkipLog: false, SkipCache: true})
   \__ These are typical exemples

- I(Named("LABEL", c(st('&'), st('@'), r("WORD")))),
   \__ This one is a handcompiled rule that therefore doesn't use a string
       (that is not going to be useful outside of this lib)
*/

// I() is a helper to declare terminal lines of rules (aka ILine).
// Since ILine are always named, you are going to always call it like this:
//
//	I(Named("hello", "'hi' | 'hello'"), <optionalCallback>)
//
// It is better to refer to the readme, as it is too flexible to explain here.
func I(a ...any) ILine {
	name, content, attrs := lineInit(a)
	if name == "" {
		panic("ILine must always be named")
	}
	return ILine{name, content, attrs}
}

func (il ILine) lineType() string { return "i" }
func (il ILine) stringIndent(nIndent int) string {
	s := helpers.Indent(nIndent)
	s += il.lineType()
	s += " " + il.name + " "
	re := regexp.MustCompile("^ *o *")
	s += re.ReplaceAllString(il.content.stringIndent(nIndent), "o ")
	if il.attrs.CbBuilder != nil {
		s += " " + Yellow("𝘧")
	}
	return s
}

func (il ILine) toRule(rank_ *rank, parentRule Parser, opts TraceOptions, lazyGrammar *helpers.Lazy[*Grammar]) (name string, rule Parser) {
	return il.name, getRule(rank_, il.name, il.content, parentRule, il.attrs, opts, lazyGrammar)
}
