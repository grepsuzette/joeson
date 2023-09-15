package joeson

// Lines of the intention grammar
// their rules are SLine, thus they will require the
// handcompiled grammar to be compiled the first time
// (after which this grammar can also parse other grammars).
func IntentionRules() []Line {
	return []Line{
		o(named("EXPR", rules(
			o("CHOICE _"),
			o(named("CHOICE", rules(
				o("_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", func(it Ast) Ast { return newChoice(it) }),
				o(named("SEQUENCE", rules(
					o("UNIT{2,}", func(it Ast) Ast { return newSequence(it) }),
					o(named("UNIT", rules(
						o("_ LABELED"),
						o(named("LABELED", rules(
							o("(label:LABEL ':')? &:(DECORATED|PRIMARY)"),
							o(named("DECORATED", rules(
								o("PRIMARY '?'", func(it Ast) Ast { return newExistential(it) }),
								o("value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", func(it Ast) Ast { return newPattern(it) }),
								o("value:PRIMARY '+' join:(!__ PRIMARY)?", func(it Ast) Ast {
									h := it.(*NativeMap)
									h.Set("min", NewNativeInt(1))
									h.Set("max", NewNativeInt(-1))
									return newPattern(h)
								}),
								o("value:PRIMARY @:RANGE", func(it Ast) Ast { return newPattern(it) }),
								o("'!' PRIMARY", func(it Ast) Ast { return newNot(it) }),
								o("'(?' expr:EXPR ')' | '?' expr:EXPR", func(it Ast) Ast { return newLookahead(it) }),
								i(named("RANGE", "'{' _ min:INT? _ ',' _ max:INT? _ '}'")),
							))),
							o(named("PRIMARY", rules(
								o("WORD '(' EXPR ')'", func(it Ast) Ast { return newRef(it) }),
								o("WORD", func(it Ast) Ast { return newRef(it) }),
								o("'(' inlineLabel:(WORD ': ')? expr:EXPR ')' ( _ '->' _ code:CODE )?", fCode),
								i(named("CODE", "'{' (!'}' (ESC1 | .))* '}'"), fCode),
								o("'\\'' (!'\\'' (ESC1 | .))* '\\''", func(it Ast) Ast { return newStr(string(it.(*NativeArray).Concat())) }),
								o("'/' (!'/' (ESC2 | .))* '/'", func(it Ast) Ast { return newRegexFromString(string(it.(*NativeArray).Concat())) }),
								o("'[' (!']' (ESC2 | .))* ']'", func(it Ast) Ast { return newRegexFromString("[" + string(it.(*NativeArray).Concat()) + "]") }),
							))),
						))),
					))),
				))),
			))),
		))),
		i(named("LABEL", "'&' | '@' | WORD")),
		i(named("WORD", "/[a-zA-Z\\._][a-zA-Z\\._0-9]*/")),
		i(named("INT", "/[0-9]+/"), NativeIntFrom),
		i(named("_PIPE", "_ '|'")),
		i(named("_", "(' ' | '\n')*")),
		i(named("__", "(' ' | '\n')+")),
		i(named(".", "/[\\s\\S]/")),
		i(named("ESC1", "'\\\\' .")),
		i(named("ESC2", "'\\\\' ."), func(chr Ast) Ast { return NewNativeString("\\" + string(chr.(NativeString))) }),
	}
}
