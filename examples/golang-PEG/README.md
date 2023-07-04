# Joeson example #3: golang-PEG

// `go test . -v`

Try to implement [golang spec](https://go.dev/ref/spec) with Joeson.

So far it's more of a proof of concept.

Content:
```
010_chars.go
011_chars_test.go
020_tokens.go
021_tokens_test.go
030_constants.go
040_variables.go
050_types.go
060_expressions.go
061_expressions_test.go
(...)
README
main.go
parseerror.go
x.go
```
Each numbered file defines a layer. 
Think of a pyramid of rules built in such a way that each layer can be tested independently from what comes on top.

Accordingly, `010_chars.go` is standalone, but `020_tokens.go` depends upon
it (it uses the `rules_tokens` array defined by `010_chars`.) 

# Methodology

We start with simple definitions like [`string_lit`](https://go.dev/ref/spec#string_lit).

From there on we are taking the pages of the specs
sequentially and reimplement them using a PEG grammar.

# EBNF used by go.dev/ref/spec

```
Syntax      = { Production } .
Production  = production_name "=" [ Expression ] "." .
Expression  = Term { "|" Term } .
Term        = Factor { Factor } .
Factor      = production_name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .

Productions are expressions constructed from terms and the following
operators, in increasing precedence:

|   alternation
()  grouping
[]  option (0 or 1 times)
{}  repetition (0 to n times)
```

∎
