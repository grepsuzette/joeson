package line

import (
	. "grepsuzette/joeson/ast"
	. "grepsuzette/joeson/core"
)

// Lines of the intention grammar
// their rules are SLine, thus they will require the
// handcompiled grammar to be compiled the first time
// (after which this grammar can also parse other grammars).
func IntentionRules() []Line {
	return []Line{
		o(named("EXPR", rules(
			o("CHOICE _"),
			o(named("CHOICE", rules(
				o("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Ast) Ast { return NewChoice(it) }),
				o(named("SEQUENCE", rules(
					o("UNIT{2,}", func(it Ast) Ast { return NewSequence(it) }),
					o(named("UNIT", rules(
						o("_ LABELED"),
						o(named("LABELED", rules(
							o("(label:LABEL ':')? &:(DECORATED|PRIMARY)"),
							o(named("DECORATED", rules(
								o("PRIMARY '?'", func(it Ast) Ast { return NewExistential(it) }),
								o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", func(it Ast) Ast { return NewPattern(it) }),
								o("value:PRIMARY '+' join:(!__ PRIMARY)?", func(it Ast) Ast {
									h := it.(NativeMap)
									h.Set("min", NewNativeInt(1))
									h.Set("max", NewNativeInt(-1))
									return NewPattern(h)
								}),
								o("value:PRIMARY @:RANGE", func(it Ast) Ast { return NewPattern(it) }),
								o("'!' PRIMARY", func(it Ast) Ast { return NewNot(it) }),
								o("'(?' expr:EXPR ')' | '?' expr:EXPR", func(it Ast) Ast { return NewLookahead(it) }),
								i(named("RANGE", "'{' _ min:INT? _ ',' _ max:INT? _ '}'")),
							))),
							o(named("PRIMARY", rules(
								o("WORD '(' EXPR ')'", func(it Ast) Ast {
									na := it.(*NativeArray)
									if na.Length() != 4 {
										panic("logic")
									}
									return NewRef(NewNativeArray([]Ast{na.Get(1), na.Get(3)}))
								}),
								o("WORD", func(it Ast) Ast { return NewRef(it) }),
								o("'(' inlineLabel:(WORD ': ')? expr:EXPR ')' ( _ '->' _ code:CODE )?", fCode),
								i(named("CODE", "'{' (!'}' (ESC1 | .))* '}'"), fCode),
								o("'\\'' (!'\\'' (ESC1 | .))* '\\''", func(it Ast) Ast { return NewStr(stringFromNativeArray(it)) }),
								o("'/' (!'/' (ESC2 | .))* '/'", func(it Ast) Ast { return NewRegexFromString(stringFromNativeArray(it)) }),
								o("'[' (!']' (ESC2 | .))* ']'", func(it Ast) Ast { return NewRegexFromString("[" + stringFromNativeArray(it) + "]") }),
							))),
						))),
					))),
				))),
			))),
		))),
		i(named("LABEL", "'&' | '@' | WORD")),
		i(named("WORD", "/[a-zA-Z\\._][a-zA-Z\\._0-9]*/")),
		i(named("INT", "/[0-9]+/"), func(it Ast) Ast { return NewNativeIntFrom(it) }),
		i(named("_PIPE", "_ '|'")),
		i(named("_", "(' ' | '\n')*")),
		i(named("__", "(' ' | '\n')+")),
		i(named(".", "/[\\s\\S]/")),
		i(named("ESC1", "'\\\\' .")),
		i(named("ESC2", "'\\\\' ."), func(chr Ast) Ast { return NewNativeString("\\" + chr.(NativeString).Str) }),
	}
}
