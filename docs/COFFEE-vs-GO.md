# This was ported from https://github.com/jaekwon/JoeScript

Originally written in 2013, the joeson parser was in Coffeescript at the time.
This was my first Go program, 
Did my best to learn on the way, but some things may not be idiomatic yet. In particular, the packrat algorithm from JoeScript is implemented as faithfully as possible. Many variable names are reused as is, same for the comments. Even now is possible to make a line per line comparison with the coffee version.

Two months were spent on deciding on a couple of things: deciding whether it was acceptable to generate a grammar using a tool such as mna/pigeon (my first prototype was using it, I also looked at pointlander/peg that Jae had also identified a couple of years ago), reading/watching as many papers and videos as possible on PEG parsers, and more importantly trying to understand what exactly was asked ("porting joeson"), and what kind of answer to give. This was the summer of 2022. 

# What a grammar looks like

In coffee impl., we would declare rules in this way:
```javascript
rulesLR = [
  o Input: "expr:Expression"
  i Expression: "Expression _ binary_op _ UnaryExpr | UnaryExpr"
  i binary_op: "'+'"
  i UnaryExpr: "[0-9]+"
  i '_': "[ \t]*"
]
```
In Go:

import j "github.com/grepsuzette/joeson"

func named(name string, v any) j.NamedRule { return j.Named(name, v) }
func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }
```go
gm := j.GrammarFromLines(
    "<some grammar name>", 
    []joeson.Line{
		o(named("Input", "expr:Expression")),
		i(named("Expression", "Expression _ binary_op _ UnaryExpr | UnaryExpr")),
		i(named("binary_op", "'+'")),
		i(named("UnaryExpr", "[0-9]+")),
		i(named("_", "[ \t]*")),
	},
)
```
The coffeescript implementation only had `Iline` and `Oline`. A few new line types were introduced in go, but they are only used internally:
 
* `cline` : compiled rule, 
* `sline` : uncompiled string rule, 
* `aline` : an array of lines (i.e. an array of rules)

Rules depending on how they are written get transformed as `sline` and `cline`. They were introduced because go is more strongly typed than js, so we chose to rely on interfaces rather than to use `any` or `interface{}`. 

A grammar which has not been compiled yet contains some `sline` (string rule). A "compiled" grammar contains only `cline` (`cline` wrapping `Parser` technically, that Parser being a special type of AST capable of parsing). The "handcompiled" joeson grammar is a compiled grammar capable of compiling other grammars ("bootstrapped").

All of this is of course never seen by the end-user. It's just the internals of it.

# Named

In coffee version, `k: v` notation is used to declare *named rules*.
This implementation uses a function called `Named()` instead. 

The following line:

```coffee
o Literal: "BasicLit | CompositeLit | FunctionLit"
```

Can be written like this:
```go
o(Named("Literal", `BasicLit | CompositeLit | FunctionLit`))
```
The syntax here is `Named(ruleName string, <rule>)` where `<rule>` can be a `string` or an array of lines (`Aline`).

Some rules are not named. Which is also why we decided to have it this way. (It is possible to simplify this, but we decided not to hide this layer for now, as it's not very hard to use and it keeps the code more "simple" - without any front layer). 

# Gnode

In original implementation, everything is a gnode (grammar, rule parsers, ast).

Present implementation parsers merely returns `joeson.Ast`. 
Joeson parsers (joeson primitives, used by rules body) return `joeson.Parser`, which are Ast with a `Parse(ctx *ParseContext) Ast` function. 

