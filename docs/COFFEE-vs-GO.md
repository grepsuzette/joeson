A few words.

This is the first golang program I wrote.
I did my best to learn on the way, but a lot of things are not idiomatic.

I also had never work with Coffeescript before. The lib developped by Jae to have dynamic classes (cardamom, clazz) were also quite an abstraction on top of Javascript. Decision was made not to use fancy stuffs such as reflection, as it would surely have ended badly.

Two months were spent on deciding on a couple of things: deciding whether it was acceptable to generate a grammar using a tool such as mna/pigeon (my first prototype was using it, I also looked at pointlander/peg that Jae had also identified a couple of years ago), reading/watching as many papers and videos as possible on PEG parsers, and more importantly trying to understand what exactly was asked ("porting joeson"), and what kind of answer to give. This was summer 2022. 

In the end, I decided to take the meaning of "porting joeson" in the most literal way.  After a first accepted version, it would probably be possible to improve the code to be more idiomatic, optimized and polished. I felt the other approaches  were too risky or took too much freedom with the question. And that they would miss features in most cases (such as would have been the case if a parser had been written and generated with a PEG parser/generator such as mna/pigeon; it would arguably not be joeson).

A port that would be absolutely literal is impossible, of course.  

`go != Javascript != Coffeescript + cardamom`. 

So let's describe the **main differences in implementation**.

In coffee impl., we would declare rules in this way:
```
rulesLR = [
  o Input: "expr:Expression"
  i Expression: "Expression _ binary_op _ Expression | UnaryExpr"
  i binary_op: "'+'"
  i UnaryExpr: "[0-9]+"
  i '_': "[ \t]*"
]
```

In go, after some aliases for O, I and Named, avoiding the namespace, we would have:

```go
gm := joeson.GrammarFromLines([]joeson.Line{
		o(named("Input", "expr:Expression")),
		i(named("Expression", "Expression _ binary_op _ Expression | UnaryExpr")),
		i(named("binary_op", "'+'")),
		i(named("UnaryExpr", "[0-9]+")),
		i(named("_", "[ \t]*")),
	}
```

The coffeescript implementation only had `Iline` and `Oline`. A few new line types were introduced in go, here they are:
 
* `Cline` : compiled rule, 
* `Sline` : uncompiled string rule, 
* `Aline` : an array of lines (i.e. an array of rules)

The programmer never makes use of those; they only exist in memory. Rules depending on how they are written get transformed as `Sline` and `Cline`. They were introduced because go is more strongly typed than js, so we need to rely on interfaces. A grammar which has not been compiled yet contains some `Sline` (string rule). A "compiled" grammar (compiled meaning parsed) contains only `Cline` (technically `Cline` wraps `Parser`, that Parser being a special type of AST capable of parsing, more about that later). The "handcompiled" joeson grammar is of course also a compiled grammar.

**Named**

In coffee version, `k: v` notation is used to declare *named rules*.
This implementation uses a function called `Named()` instead. 

The following line:

```coffee
o Literal: "BasicLit | CompositeLit | FunctionLit"
```

Can be written as follows:

```go
o(Named("Literal", "BasicLit | CompositeLit | FunctionLit"))
```

The syntax here is `Named(ruleName string, <rule>)` where `<rule>` can be a `string` or an array of lines (`Aline`).

Some rules are not Named. Which is also why we decided to have it this way. (It is possible to simplify this, but we decided not to hide this layer for now, as it's not very hard to use and it keeps the code more "simple" - without any front layer). 

**gnode**

In original implementation, everything is a gnode. Grammar, rule parsers, ast. 

In present implementation, `gnode` is used by the `Parser` objects, but this is somewhat hidden for the end-user, indeed `gnode` is not exported. 

Instead we use `Ast` and `Parser` interfaces. 

* `Ast` 
  - the product of a grammar parsing a string (`func Parse(ctx *ParseContext) Ast`). 
  - it is an **interface** defining String()string (`type Ast interface { ContentString() string }`). 

* `Parser` 
  - a special `Ast` produced by a joeson grammar (`Ref, Rank, Sequence, LookAhead, Choice, Not, Pattern, Regex`...). 
  - It allows to `Parse()` some ParseContext and returns an `Ast`.
  - It implements `gnode` (meaning it has `Name()`, `Label()` etc). 

* `Grammar`
  - implements `gnode` too


**"Native" types**

They are used for the returns. 

Since every parser must return an `Ast`, we can not directly return an int, string, map, or array. Instead we have wrappers that also implement `String()`. They are `NativeMap`, `NativeArray`, `NativeInt`, `NativeString` and are all `Ast`.

How does it work?

Consider the following simple rules, randomly extracted from examples/calculator:

```go
o(named("Input", "expr:Expression")),
i(named("MulOp", "'*' | '/'")),
i(named("Integer", "/^-?[0-9]+/"), func(it joeson.Ast) joeson.Ast { return joeson.NewNativeIntFrom(it) }),
```

* The first rule `Input` contains a **label** the label is "expr". 
  - A `NativeMap` will be automatically returned by the implicit `Parser`.
  - it will have a single entry: { "expr": Ast } where `Ast` is whatever is returned by the rule for `Expression`   (determined by the callback for the rule, or by the `Parser` that the parent grammar compiled).

* The second rule `MulOp`
  * `i(named("MulOp", "'*' | '/'"))`
  * defines no label 
  * is an alternation of two strings, therefore it will return a NativeString

* The third rule `Integer`
  * `i(named("Integer", "/^-?[0-9]+/"), func(it joeson.Ast) joeson.Ast { return joeson.NewNativeIntFrom(it) })`

The third rule `Integer` (compiled as a `Regex --|> Parser --|> Ast` because of the `//`), has a **callback** defined. This callback allows it to work on the `it Ast` that is returned by the `Regex` parser. It's a `NativeString` (that's what `Regex` returns) and it creates a `NativeInt` from `it`.

The behavior described above is pretty standard in all PEG parsers. The rules we have examined are perhaps a bit simplistic. Even though it is difficult to explain in plain english, it is not very hard to write rules and callbacks with some practice. Please see examples/ for more practical examples.

**Performances**

Premature optimization is the root of all evil, however some optimization is in order. 

The time it takes for a few things I tested seemed comparable to the implementation in coffeescript, a go implementation should be faster. I believe packrat.go in particular must be slow, as it uses many embedded callbacks (`ParseFunc` and `ParseFunc2`), as was the case in coffeescript. The other slow part, that was
already somewhat measured is parser_sequence.go (parseAsSingle and parseAsObject are both called 100k times
to parse the intention grammar, while this is the same number of times as the original implementation, those
must be considered critical parts).

Another thing to likely do to optimize is to use pointer receivers.

**Error handling**.

`Parse(ctx) Ast` must return `nil` when the current parser failed to recognize anything. Please note it is also possible to return an `AstError` when current context *is* the correct one and you somehow want the parsing to produce an error.

For example, parsing "0o9" for an octal parser may return `NewAstError(parser, "0o9" is invalid octal")`, thus denying further parsing attempts and giving a precise error message. 

`ParseError` is simply:

```go
type ParseError struct {
	ctx         *ParseContext
	ErrorString string
}

func (pe ParseError) String() string {
	return "ERROR " + pe.ErrorString + " " + pe.ctx.String()
}

func NewParseError(ctx *ParseContext, s string) ParseError {
	return ParseError{ctx, s}
}

func IsParseError(ast Ast) bool {
	switch ast.(type) {
	case ParseError:
		return true
	default:
	}
	return false
}
```

Note the upper level code (in grammar.go) uses some panic right now. 
This is during the necessary time for the code to mature a bit.

The Parse methods also take a very simple `string` or `CodeStream` argument right now.
In other words, the high level methods are not very polished yet and subject to change.

The grammars however should not evolve. This part is to be considered stable.
