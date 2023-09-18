**Joeson** is a go Packrat parser with left recursion, memoization, and tree-like rules. 

This document further explores Joeson's specificities, installation, testing, operators and error handling.

Joeson dynamically creates parsers without code generation.
It was ported from [https://github.com/jaekwon/JoeScript](https://github.com/jaekwon/JoeScript) for the [Gnolang](https://github.com/gnolang/gno) project.

Joeson differs from typical PEG parsers based on Ford et al. in several ways:

- It uses `|` for *horizontal* alternation instead of `/`, 
- It generates parsers at runtime instead of generating code,
- It introduces a few new operators like the join pattern (`value*join`) and the occurrence operator (`{m,n}`),
* It offers a *vertical* alternation system, enabling **tree-like grammars of rules**,

This last feature makes Joeson more difficult to learn but potentially interesting for large grammars.

## Installation and test

To install, run the following command:
```bash
go get -u github.com/grepsuzette/joeson
```
Make sure to use the `-u` flag to update the library if it's already installed.

To run all tests recursively, use:

```
go test ./...
```

To run a named test with traces enabled, use:
```
TRACE=all go test . --run TestHandcompiled -v
```

To run a specific test from examples/calculator named "test_12" with verbosity, use:
```
go test examples/calculator/calculator_test.go --run test_12 -v
```

![calculator_test1](https://user-images.githubusercontent.com/350354/216583474-4fd47a26-54a1-400a-aba6-96af1b06188f.png)

### Tracing

When developping a grammar, you will usually need to enable some form of tracing.

To run a certain test showing the grammar, use:
```
TRACE=grammar go test examples/calculator/calculator_test.go --run test_12 -v
```

By default, grammars verbosity can be controlled by using the TRACE environment variable. 
Here are the useful values:

| Name       | Behavior                                          |
| ---------- | ------------------------------------------------  |
| stack      | print detailed parsing steps                      |
| grammar    | print grammar information and all rules           |
| all        | enable everything                                 |

Read [docs/verbosity.md](docs/verbosity.md) for more.

> To run a certain test showing the parse trace:
```
TRACE=stack go test examples/calculator/calculator_test.go --run test_12 -v
```
![calculator_test2](https://user-images.githubusercontent.com/350354/216583646-0009d83d-36c2-457b-8cc3-e2aa0012edac.png)

## Header

The **following header declarations** will be used (implicitely) throughout the rest of this document:
```go
import j "github.com/grepsuzette/joeson"

// named() creates a rule with a name
// E.g. o(named("Alpha", `'α'`)) // named rule
//      o(`'α'`)                 // the same rule, anonymous
func named(name string, v any) j.NamedRule { return j.Named(name, v) }

// to make i and o rules
func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }
```
That is pretty much all the functions you are going to need with `j.GrammarFromLines()`. Don't worry about the difference between `i` and `o` for now.

If you encounter any difficulty note all examples appearing here (some
functions are not written to make the samples shorter) you can find them in their entirety in
the `examples/` folder.

## Operators 

The following operators are supported:

* `'s'`: Match string literal *s*
* `[c]`: Match against character class *c*
* `(e)`: Grouping
* `e?`: Match *e* zero or one time
* `e*`: Match *e* zero or more times
* `e+`: Match *e* one or more times
* `e{m,n}`: Match *e* t times with t ∈ [m;n]
* `e1*e2`, `e1+e2`, `e1*e2{m,n}`: See **pattern** below
* `?(e)`: Non-consuming match of *e*
* `!e`: Non-consuming negative match of *e*
* `e1 e2`: Sequence
* `e1 | e2`: Ordered choice
* `/r/`: Match against regular expression *r* (discouraged)

Rule bodies should be placed into backquotes, to differentiate them from rule *names* (double-quoted) and *string literal matches* (single-quoted). For example:
```go
o(named("FourthRule", `'only two guys to a fight'`)),
```
Operators will now be reviewed:

* `α β`: **sequence**. α followed by β. 
    * Return `nil` if failed.
* `α | β | γ`: **choice** or "alternation". α, or β, or γ, or `nil`, in that order.
    * Return the parse result of any element that would match, starting from the left.
    * Alternation may also appear as ranks (see section below).
    * Note: most PEG grammars use `/` instead.
* `'xyz'`: **str** ("string" parser). Matching a string, here "xyz". 
    * The string must match exactly. Case is sensitive. 
    * This parser does *not capture* anything *unless* it is the only one inside the rule body. 
    * To capture, you must either: 
      1. reference such string by their name in their own rule (`i(named("COMMA", `','`))`), 
      2. use a character class. For example, in `α [.]` the dot will be captured.
* `(λ)`: grouping. E.g. `('bar')+`.
* `[λ]`: character class.
    * Captures 1 rune.
    * Uses the character class from the `regexp` package,
    * Following are valid: `[a-zA-Z0-9_-]`, `[xyz]`
    * DON'T use negation: `[^a-z]` TODO forbid it programmatically
    * `[xyz]` is often written `'x' | 'y' | 'z'` (more idiomatic in PEG grammars)
    * `[x]` is often used to capture character 'x', since `'x'` is non-capturing
* `α+`, `α*`, `α{n,m}`, `α*β`, `α+β`, `α*β{n,m}`: **pattern** 
    * `α+` matches one or more `α`
    * `α*` matches zero or more `α`
    * `α{n,m}` matches between n and m times, n or m can be omitted. This is usually not in PEG parsers.
    * `α*β` and all the patterns below are joeson specific, `β` is called the *join*.
       * `'a'*'b'` parses `"", "a", "aaa", "aba", "aabababa"`, but not `"abab"`.
    * `α+β` is similar but does not parse `""`. 
       * `'a'+'b'` parses `"a", "aba", "abaaaba"`. 
       * In other words it requires at least one `α`.
    * `α*β{n,m}` is as shown below:
       * `'a'*'b'{1,2}` parses `"a", "aaaa", "aaabaaa"`. 
       * In other words, the `{}` notation repeats the whole `a*b` pattern. 
* `α?`: **existential** operator returning:
    * production of `α` when not `nil`
    * `NativeUndefined` otherwise
* `?α` or `(?α)`: **lookahead**. 
    * Ford et al. called this "non-consuming match" and used notation `&e`.
    * input is not consumed,
    * Returns `NativeUndefined` if matches `α`
    * Returns `nil` if no match.
    * Better use round brackets:
        * `α ?Ω | β` is grouped as `α (?Ω | β)`, 
        * so be explicit: `α (?Ω) | β`
* `!α` or `(!α)`: **negative lookahead** ("not").
    * Same as lookahead, but with reversed meaning.
    * The same considerations apply (use parens with alternation).
* `/ρ/`: **regexp**.
    * Joeson supports regexp,
    * However Regexp can easily be greedy,
    * Prefer character classes `[xy]` or other operators like `('x'|'y')`.
    * Greediness creates ambiguity, ambiguity is the devil in PEG, or at least a major challenge

## Writing a grammar

Joeson possesses certain specificities that require further explanation.

Writing a grammar, regardless of the parser, normally involves 3 parts:

1. *Rules*,
2. *Parse functions*, which are optional callbacks or function name at the end of a rule, mapping a parse result to a different parse result (signature like `func(Ast) Ast`),
3. *Input data*. Those are strings, code samples or anything which should be parsable by the grammar, and often constituting the test data.

### How joeson works

For now, we will skip the more complex details, such as the
 distinction between `o` and `i `rules. These will be explained in
 the "Alternation: Vertical Ranks" section, which will provide
 clarity on this.

In Joeson, a *rule* is a **string using joeson operators, labels, captures**. 

In the provided Go code examples, we conventionally use **backquotes** to denote the body of a rule.

Consider the rule `'a'+ PARAGRAPH? | 'END'`. It is composed of a bipartite choice that can match either:

1. One or more lowercase "a" characters, optionally followed by a PARAGRAPH element ("PARAGRAPH" is the reference to the rule named "PARAGRAPH"),
2. The letters "E" "N" "D", appearing in that order.

If none of these possibilities match, the rule will return `nil`. In this case, the upstream rule will attempt the next alternation if it has one, and so on recursively, as is standard for PEG parsers.

When entering a rule as a string, it is parsed by Joeson to generate a number of parsers internally (this is called "compiling" the rules of a grammar). 

Although this internal process is not mandatory to understand when using Joeson, it will bring confidence if any sort of magic is dispelled for you. The example `'a'+ PARAGRAPH? | 'END'` will invoke the following parsers: `parser_choice`, `parser_sequence`, `parser_str`, `parser_pattern`, `parser_ref`, `parser_existential`. These are not very different from parser functions you may write yourself. They however produce a special kind of `Ast` called `Parser`. They are the building blocks with which Joeson dynamically compiles parsers for your rules. 

Conceptually, the Joeson parsers mentioned above take **runes** as an input from a `ParseContext` and output an object implementing `Ast`.

An `Ast` is basically any type that has a `String()` method along with a few other convenience methods (in practice, implement `Ast` by embedding an object `j.Attr` in your type and by defining `String()`, see `examples/videogamedb`). Typically, custom **parser functions** have this signature:
```go
func (it Ast) Ast
```
These parse functions map one `Ast` to another (it's a "mapper"). In certain cases where the `ParseContext` is needed (usually, to generate a `ParseError`), the signature can be written as:
```go
func (it Ast, ctx *ParseContext) Ast
```
But enough theory, let's see it in practice!

### Creating a simple grammar

Invoke GrammarFromLines with its title and a list of rules and you get a dynamically compiled grammar:
```go
gm := j.GrammarFromLines("my grammar", []j.Line{
    o(named("INPUT", `'hi' _ NAME`)),
    i(named("_", `' '*`)),
    i(named("NAME", `[a-zA-Z]*`), func(it j.Ast) j.Ast {
        return it.(*j.NativeArray).Concat() 
    }),
})
```
The **first rule** declared **in a grammar**, which we may call the **entry rule** (in this case named "INPUT"), is the rule the *entire text is parsed against*. In this case, it requires the text to start with `'hi'`, followed by one or more space characters, followed by a NAME.

Therefore, "hi amigo" would parse, "bye bye" wouldn't.

When a rule has a parse function, it is called with the production of the joeson parsers for that rule. Let's add a parse function to our entry rule:
```go
gm := j.GrammarFromLines("my grammar", []j.Line{
    o(named("INPUT", `'hi' _ NAME`), func(it j.Ast) j.Ast {
        fmt.Printf("I am not %s.\n", it.(j.NativeString).String()) 
        return it
    }),
    i(named("NAME", `[a-zA-Z]*`), func(it j.Ast) j.Ast {
        return it.(*j.NativeArray).Concat()
    }),
    i(named("_", `' '*`)),
})
```
If we ran `gm.ParseString("hi Valery")` it would therefore print "I am not Valery.".

It's worth noting that "INPUT" is not referenced anywhere in the grammar. The entry rule could therefore have been declared like this instead: 
```go
gm := j.GrammarFromLines("my grammar", []Line{
    o(`'hi' _ NAME`, func(it j.Ast) j.Ast {
        fmt.Printf("I am not %s.\n", it.(j.NativeString).String()) 
        return it
    }),
    i(named("NAME", `[a-zA-Z]*`), func(it j.Ast) j.Ast {
        return it.(j.*NativeArray).Concat()
    }),
    i(named("_", `' '*`)),
})
```
There are still several unanswered questions, such as the purpose of the `i` and `o` rules. 

However, before diving into that, we need to understand the special `Native***` types that we used above.

### Native types

Joeson parsers are required to return an `Ast` object. 

Since `Ast` is an interface, the types commonly produced by joeson
parsers (string, int, array, map, and a special undefined value that is distinct from `nil`) need to be wrapped in the following types:

* `*NativeMap`
* `*NativeArray`
* `NativeInt`
* `NativeString`
* `NativeUndefined`

Additionally, there is a `ParseError` that will be discussed later in this document.

Knowing which type is returned by a certain joeson takes some practice. The following rules of thumb may help you to start:

* Strings in a sequence are not captured unless they are the only thing appearing in a rule (rule `'foo' BAR` will only return the production of rule named "BAR"),
* A rule name that starts with an underscore (e.g. "_", "_COMMA") will not be captured in the sequence it appears in.
* An existential operator (`Element?`) will either be `NativeUndefined` or the production of `Element`.
* If a sequence uses labels, it will return a `NativeMap` (`'a' color:COLOR vehicle:('car' | 'truck')` will return a NativeMap with keys "color" and "vehicle"). 
* If a sequence has multiple (captured) values but no labels, it will return a `NativeArray`.
* If a sequence consists of a single element, it will return the unadorned production of that element.

Now you should have a basic understanding, there remains two important items that require explanation:

1. Vertical ranks
2. How to generate custom parse errors

The first item is quite unique, so you may want to take a short break and have a cup of coffee before proceeding.

### Alternation: Horizontal Choice VS Vertical Ranks

There are two ways to write alternation in joeson: horizontally using "choices" (`α | β | γ`) or vertically using "ranks".

Consider the following grammar:
```go
gm := j.GrammarFromLines("example1", []j.Line{
    o(named("VideoGame", `id:[1-9][0-9]* | '"' title:([^"]*) '"'`), makeVideoGame),
})

Although the grammar itself is simple, the callback can be tedious to write.

As the number of choices increases, the complexity of the callback function also increases. This can make it difficult to expand or maintain the grammar. Ideally, since parse functions are mere mappers (`Ast -> Ast`) they should not be large or complicated.

To simplify this, we can break down the alternation into a "rank".

In Joeson, although you won't instantiate **ranks** directly, they are
 semantically central. In fact every time you define rules within `j.GrammarFromLines(name string, rules []j.Line)`, you are working with ranks.

A **rank** is a list of `o` rules. Each of those rules will try to be matched successively, as in a choice. 

* A **choice** is horizontal: `choice1 | choice2 | choice3`; 
* a **rank** is the same thing, expressed vertically:
```go
o(named("choices", []j.Line{
    o(named("choice1", ...), /* optional callback */),
    o(named("choice2", ...), /* optional callback */),
    o(named("choice3", ...), /* optional callback */),
    /* any number of i rules, 
       they can be appear here but will 
       not be part of the alternation */
}))
```
By using ranks, joeson grammars effectively...:

1. allow per-choice callback (1 callback per line, instead of only one for an horizontal choice);
2. become **trees** of rules (because you can have list of rules at any depth)

Here is how to turn the horizontal choice in "VideoGame" of example 1 into a rank:
```go
gm := j.GrammarFromLines("example2", []j.Line{
    o(named("VideoGame", []j.Line{
        o(`[1-9][0-9]*`, func(it j.Ast) j.Ast { return findVideoGameById(it) }),
        o(`'"' [^"]* '"'`, func(it j.Ast) j.Ast { return findVideoGameByTitle(it) }),
    })),
})
```
When "VideoGame" rule will be used to parse a string:

* If `[1-9][0-9]+` matches, findVideoGameById will be called.
* If not, the second `o` rule in the rank will be matched against, and if so findVideoGameById will be called.
* If not, since there is no more `o` rules in the rank, rule "VideoGame" will return `nil`.

Note we don't need the callback but we can reference the functions directly, since they both take an `Ast` and return an `Ast`:
```go
gm3 := j.GrammarFromLines("example3", []j.Line{
    o(named("VideoGame", []j.Line{
        o(`[1-9][0-9]*`, findVideoGameById),
        o(`'"' [^"]* '"'`, findVideoGameByTitle),
    })),
})
```
We can also choose to name unreferenced rules:
```go
gm4 := j.GrammarFromLines("example4", []j.Line{
    o(named("VideoGame", []j.Line{
        o(named("VideoGameId", `[1-9][0-9]*`), findVideoGameById),
        o(named("VideoGameTitle", `'"' [^"]* '"'`), findVideoGameByTitle),
    })),
})
```
Now, let's introduce the final semantic element in the same example - the difference between o and i rules. 

Consider the example provided below:
```go
gm5 := j.GrammarFromLines("example5", []j.Line{
    o(named("VideoGame", []j.Line{
        o(named("VideoGameId", `[1-9][0-9]*`), findVideoGameById),
        o(named("VideoGameTitle", `'"' [^"]* '"'`), findVideoGameByTitle),
        o(named("VideoGameBestOfYear", `'bestIn:' _ Year`), findBestVideoGameOfYear),
        i(named("Year", `a:('19'|'20'|'21') b:[0-9] c:[0-9]`), func(it j.Ast) j.Ast { 
            return it.(*j.NativeMap).Concat() 
        }),
    })),
    i(named("_", `[ \t]*`)),
})
```
The vertical alternation for rank "VideoGame" here consists of only 3 direct children, namely the o rules:

1. o rule "VideoGameId", 
2. o rule "VideoGameTitle",
3. o rule "VideoGameBestOfYear"

The `i` rules are never part of a rank. 
They can in fact be placed anywhere in the tree.

* A `i` rule **MUST be named** (they are inert if not referenced);
* A `i` rule can also be a rank.

One last example for the road:
```go
type Color struct {
	*j.Attr
	r, g, b int
}

func (c Color) String() string {
	return fmt.Sprintf("Color is rgb ( %d, %d, %d )", c.r, c.g, c.b)
}

func toArrayInt(it j.Ast) []int {
	r := []int{}
	for _, v := range it.(*j.NativeArray).Array() {
		r = append(r, j.NativeIntFrom(v).(j.NativeInt).Int())
	}
	return r
}

gm := j.GrammarFromLines("color example", []j.Line{
    o(named("Color", []j.Line{
        o("'red'", func(it j.Ast) j.Ast { return Color{j.NewAttr(), 255, 0, 0} }),
        o("'green'", func(it j.Ast) j.Ast { return Color{j.NewAttr(), 0, 255, 0} }),
        o("'blue'", func(it j.Ast) j.Ast { return Color{j.NewAttr(), 0, 0, 255} }),
        o(named("Rgb", `'rgb' _ '(' _ Integer*_COMMA{3,3} _ ')'`), func(it j.Ast) j.Ast {
            a := toArrayInt(it)
            return Color{j.NewAttr(), a[0], a[1], a[2]}
        }),
        o(named("Hsl", `'hsl' _ '(' HslTrio ')'`), func(it j.Ast) j.Ast {
            a := toArrayInt(it)
            b := hslToRgb(a[0], a[1], a[2])
            return Color{j.NewAttr(), b[0], b[1], b[2]}
        }),
        i(named("HslTrio", `_ Integer _ ',' _ Integer '%' _ ',' _ Integer '%' _`)),
    })),
    i(named("Integer", `[1-9][0-9]*`), j.NativeIntFrom),
    i(named("_COMMA", `[ ,]*`)),
    i(named("_", `[ \t]*`)),
})
gm.ParseString("blue")
gm.ParseString("red")
gm.ParseString("green")
gm.ParseString("rgb(127, 49, 255)")
gm.ParseString("hsl(127, 49%, 82%)")
```
## Defining your own AST types

To be able to define your own AST types, your nodes simply need to:

* embed `*j.Attr`
* implement `String() string`

See `examples/videogamedb` for a full example.

## Error handling

* joeson operators return `nil` when a sequence failed to parse. 
* this will pass the hand so an alternation higher up tries another choice,
* `NativeUndefined` is returned by the *existential* operator when the referenced element failed to parse (`Element?`).
* *blabla bla :)*

Knowing this, how to produce a custom and clean parse error?

The best way is to panic with a `ParseError`, and to recover from that panic.

Consider the following rule which parses octal numbers in a string (e.g. "\142" in Go) with parse errors:
```go
i(
    named("octal_byte_value", `'\' [0-7]+`), 
    func (it j.Ast, ctx *j.ParseContext) j.Ast {
        a := it.(*j.NativeArray)
        if a.Length() < 3 {
            panic(ctx.Error("illegal: too many octal digits"))
        } else if a.Length() > 3 {
            panic(ctx.Error("illegal: too much octal digits"))
        }
        var n j.NativeInt = j.NewNativeIntFrom(a.Concat())
        if n.Int() > 255 {
            panic(ctx.Error("illegal: octal value over 255"))
        }
        return n // NativeInt being an Ast
    }
),
```

The reason we recommend to panic with a special `Ast` type `ParseError` is because parsing is a recursive process. 

So you either use a mechanism like monads or manual checking for ParseError at *each level of your rules*; Or you go out-of-band and panic with a `ParseError`, to eventually recover and bring back the `ParseError` into the band at the top-level. 

Something like this does just that:
```go
func parseX(gm j.Grammar, s string) (result j.Ast) {
	defer func() {
		if e := recover(); e != nil {
			if pe, ok := e.(j.ParseError); ok {
				result = pe
			} else {
				panic(e)
			}
		}
	}()
    result = gm.ParseString(s)
	return
}
```
This will either:

* return a `ParseError`,
* return an `Ast`,
* panic, when unexpected types are used.

Joeson was ported from [https://github.com/jaekwon/JoeScript](https://github.com/jaekwon/JoeScript) for the [Gnolang](https://github.com/gnolang/gno) project.

You may find more docs here:

- [docs/testing](docs/testing.md)
- [docs/faq](docs/faq.md)
- [docs/diffing](docs/diffing.md)
- [docs/internals](docs/internals.md)

