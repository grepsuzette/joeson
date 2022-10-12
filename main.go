package main

import (
	"fmt"
	// "regexp"
	"strconv"
)

type Sequence struct{}
type Heredoc struct{}
type If struct{ n int }

func makeIf(it ctx) If {
	return If{2392}
}

func (ch Choice) String() string   { return "oyeah the choice!" }
func (sq Sequence) String() string { return "oyeah the seq!" }
func (hr Heredoc) String() string  { return "oyeah the heredoc!" }
func (ifs If) String() string      { return "oyheann if" }

type ioCallback func(ctx ctx) astnode
type ctx interface { // parsing context
	fancystuffs() string
	andthings() string
}
type sequenceString string

// type grammar []iorule

type token string

// func tokens(grammar grammar, vrd ...token) {
// 	// ignore cb
// 	regexAll := regexp.MustCompile(
// 		"[ ]*(",
// 		strings.Join(tokens, "|"),
// 		")([^a-zA-Z\\$_0-9]|$)",
// 	)
// 	for _, token := range tokens {
// 		name := "_" + strings.ToUpper(token)
// 		// HACK: TODO temporarily halt trace
// 		rule := grammar.parse("/[ ]*/ &:'"+token+"' !/[a-zA-Z\\$_0-9]/")
// 		rule.rule = rule
// 		rule.skipLog = true
// 		rule.skipCache = true
// 		rule.cb = nil // cb if cb?
// 		regexAll.include(name, rule) // TODO whatisit
// 	}
// 	OLine(regexAll)
// }

type what string

func makeNode(clazzname string /*options=undefined*/) func(it what, ctx ctx) astnode {
	return func(it what, ctx ctx) astnode {
		return makeclazz(clazzname /*, options*/)
	}
}
func makeclazz(clazzname string /*, options*/) astnode {
	switch clazzname {
	case "Choice":
		return Choice{}
	case "Sequence":
		return Sequence{}
	case "Heredoc":
		return Heredoc{}
	default:
		return Choice{}
	}
}

func (a iorules) toPigeon(indent int) {
	for _, rule := range a {
		rule.toPigeon(0)
	}
}
func (rule iorule) toPigeon(indent int) {
	fIndent := func(indent int) {
		for i := 0; i < indent; i++ {
			fmt.Print("  ")
		}
	}
	fIndent(indent)
	fmt.Println(rule.sym)
	for _, subrule := range rule.rules {
		subrule.toPigeon(indent + 1)
	}
}

// func newGrammar(rank rank) grammar {
// TODO if rank is function
// if rank.isArray {
// 	rank = rankFromLines("__grammar__")
// }
// g = grammar{rank{rank}}
// return grammar{rank{}}
// }
func parseRawGrammar(raw iorules) Grammar {
	return Grammar{rank{}, 0}
}

func main() {
	// what if oid=o and
	// depending on its string arg:
	//   /^[A-z_.]+$/ => it is named
	//   otherwise    => unnamed, and a mere alternation choice
	var rulez []iorule = rules(
		oid("EXPR", rules(
			o("CHOICE _"),
			oid("CHOICE", rules(
				ocb("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(res *Result) astnode {
					return Choice{}
				}),
				oid("SEQUENCE", rules(
					ocb("UNIT{2,}", func(ctx ctx) astnode {
						return Sequence{}
					}),
					oid("UNIT", rules(
						o("_ LABELED"),
						oid("LABELED", rules(
							o("(label:LABEL ':')? &:(DECORATED|PRIMARY)"),
							oid("DECORATED", rules(
								ocb("PRIMARY '?'", nil),
								ocb("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", nil),
								ocb("value:PRIMARY '+' join:(!__ PRIMARY)?", nil),
								ocb("value:PRIMARY @:RANGE", func(res Result) astnode {
									// ^ should above give PRIMARY and RANGE,
									// with RANGE have 2 args (min and max)
									// and thus PRIMARY maybe having 2 as weell
									// so as to give 4 (because Pattern needs 4)
									// ??
									return NewPattern(res.ExtractPatternArgs())
								}),
								// o "'!' PRIMARY", (it) -> new Not it
								// o "'(?' expr:EXPR ')' | '?' expr:EXPR", (it) -> new Lookahead it
								i("RANGE", "'{' _ min:INT? _ ',' _ max:INT? _ '}'"),
								// ^- this should give 2 args?
								//    min:INT?  and max:INT?     right?
							)),
							oid("PRIMARY", rules(
								ocb("WORD '(' EXPR ')'", func(res Result) astnode {
									name, param := res.ExtractTheMatchIOrWhateverPossiblyWithParam()
									return NewRef(name, param)
								}),
								ocb("WORD", func(res Result) astnode {
									name, param := res.ExtractTheMatchIOrWhateverPossiblyWithParam()
									return NewRef(name, param) // think param may be nil
								}),
							)),
						)),
					)),
				)),
			)),
			i("LINE", rules(
				ost("HEREDOC", " _ '###' !'#' (!'###' .)* '###' ", func(it ctx) astnode { return Heredoc{} }),
				o("LINEEXPR", rules(
					// left recursive
					ost("POSTIF", " block:LINEEXPR _IF cond:EXPR ", func(it ctx) astnode { return makeNode("If")("", it) }),
					ost("POSTIF2", " block:LINEEXPR _IF cond:EXPR ", func(it ctx) astnode { return makeIf(it) }),
				)),
			),
			),
			i("LABEL", "'&' | '@' | WORD"), // TODO
			i("WORD", "/[a-zA-Z\\._][a-zA-Z\\._0-9]*/"),
			i("INT", "/[0-9]+/", func(res Result) astnode { return NewNativeIntFromString(res.ExtractString()) }),
			// ^ this parses as a Regex... Parse of Regex gives
			// ResultIsString. That string is obtained as
			// ExtractString(). NativeInt is a way to be terminal.
			i("_PIPE", "_ '|'"),
			i("_", "(' ' | '\n')*"),
			i("__", "(' ' | '\n')+"),
			i("'.'", "/[\\s\\S]/"),
			i("ESC1", "'\\\\' ."),
			i("ESC2", "'\\\\' .", func(res Result) astnode { return NewStr('\\' + res.ExtractChar) }),
		)),
	)
	var RAW_GRAMMAR iorules = rulez
	// var PARSED_GRAMMAR Grammar = parseRawGrammar(RAW_GRAMMAR)
	rank := RankFromRules(rules)
	var PARSED_GRAMMAR Grammar = NewGrammar(RAW_GRAMMAR)
	m := map[sequenceString][]iorule{
		"default": rulez,
	}
	fmt.Println(m)
	RAW_GRAMMAR.toPigeon(0)
	fmt.Println(PARSED_GRAMMAR)
}
