This port is as literal as possible. It's main documentation is meant to be the original JoeScript code. Once it becomes more established it can become its own thing and evolve in a more golang-idiomatic way.

# `line_*.go` files 

They contain the rule declaration helpers called `O()` 
and `I()` which allow declaring and compiling tree-like grammars.

# Line interface

`Line` is a very simple interface:

```go
type Line interface {
	LineType() string                // i, o, a, s, c
	Content() Line                   // Sline, OLine, ALine, CLine (containing an Ast)...
	StringIndent(nIndent int) string // indent with `nIdent` levels (for nested rules)
}
```

It allows to declare arrays of `Line` that can contain any line type.

# Line types

## OLine and ILine

The only `Line` types which are used to explicitely write a grammar are `OLine`
and `ILine`, respectively through the usage of `O()` and `I()` functions.

- `OLine` expresses a **non-terminal rule**.
- `ILine` expresses a **terminal rule** *in the current depth*.

TODO explain `Named`
TODO explain when to use Named

Why do we write *in the current depth*?
Because in joeson, rules form a *tree* and not just a linear array of rules.

## Transient and compiled line types

The `line` package produces more `Line` types than just `ILine` and `OLine`.
That's because golang is strongly typed.

We'll verbatim the comment for `getRule()`, which comes from a comment in the
coffeescript version (save for the indented lines in the middle):

```go
// name:       The final and correct name for this rule
// rule:       A rule-like object
//                 In coffee it means string, array, object (map) or oline
//                 In this implementation it means Line, amongst:
//                   SLine (for string), ALine, OLine
// parentRule: The actual parent Rule instance
func getRule(..., name string, line Line, parentRule Ast, ...) Ast {
```

To sum it up:

- `ALine` helps embedding a `[]Line`
- `SLine` helps embedding a `Str` (a string rule yet to be compiled). A compiled `SLine` becomes an `OLine`.
- `OLine` becomes a `CLine` when compiled.
- `ILine` embeds either an `OLine` or a `CLine`.
- `CLine` is a compiled rule and embeds an `Ast` node (for which `IsRule()` is true).

# Compilation mechanism (internals)

Notable methods:

* `OLine.toRule(parentRule, accessBy)`
* `ILine.toRules(parentRule)`
* `getRule(name, rule, parentRule, attrs)`

Packrat are top-down parsers. In the original joeson.coffee implementation (which is reused here)
the pump is initialized by a top-level O() toRule call, and then (in pseudo-code):

```
+-----------+
|+----------+-----------------------------------------+                            
|| OLine    v                                         |                            
|| func toRule:                                       |     +--------------------------------------------------+                          
||   if not name and ...                              | .---+>func getRule(name, rule, parentRule, attrs) {    |                          
||      ...                                           |/    |     ...                                          |                          
||   else if not name? and index? ...                 |     |     if rule instanceOf Array                     |                          
||      name = ...                                   /|     |         rule = rule.toRule parentRule name:name  |                          
||   rule =  getRule(name, rule, parentRule, attrs) / |     |     ...             ======                       |                          
||   ...     =======--------------------------------  |     | }                      |                         |  
||                                                    |     +------------------------+-------------------------+  
|+----------------------------------------------------+                              |
|                                                                                    |
+------------------------------------------------------------------------------------+
```

As you can see it's a cross-recursion which ultimately will parse all the
grammar. The same is also true in this implementation. 

A slight difference is the original `GetRule` has been inlined in ILine.ToRule and OLine.ToRule.

