; vim: ft=markdown
```go
o(Named("LABELED", rules(
    o(S(E(S(L("label", R("LABEL")), St(":"))), L("&", C(R("DECORATED"), R("PRIMARY"))))),
o(S(St("("), L("inlineLabel", E(S(R("WORD"), St(": ")))), L("expr", R("EXPR")), St(")"), E(S(R("_"), St("->"), R("_"), L("code", R("CODE"))))), func(it Astnode) Astnode {
//Raw:
o(Named("LABELED", Rules(
	o("(label:LABEL ':')? &:(DECORATED|PRIMARY)"),
o("'(' inlineLabel:(WORD ': ')? expr:EXPR ')' ( _ '->' _ code:CODE )?", func(it Astnode) Astnode {
```
0,0                  CHOICE _            ] | | | | | | | | | * LABELED[0]: (@:(label:LABEL ':')? &:(DECORATED | PRIMARY)) 11
0,8             ICE _                    ] | | | | | | | | | | | | | * PRIMARY[2]: ('(' inlineLabel:(WORD ': ')? expr:EXPR ')' @:(_ '->' _ code:CODE)?) 55 

Weird
=====================================
in js we have below. But in go, we never have the lines from 3 (LABEL: ('&' | '@' | WORD))
 PRE connect nodes, parent:undefined node: undefined GRAMMAR{Rank Rank(EXPR)} type:Grammar
 PRE connect nodes, parent:GRAMMAR{Rank Rank(EXPR)} node: undefined Rank(EXPR) type:Rank
 PRE connect nodes, parent:Rank(EXPR) node: LABEL: ('&' | '@' | WORD) ('&' | '@' | WORD) type:*ast.Choice
 PRE connect nodes, parent:('&' | '@' | WORD) node: undefined '&' type:ast.Str
 PRE connect nodes, parent:('&' | '@' | WORD) node: undefined '@' type:ast.Str
 PRE connect nodes, parent:('&' | '@' | WORD) node: undefined WORD type:*ast.Ref
Loop 0: LABEL: ('&' | '@' | WORD)
 PRE connect nodes, parent:Rank(EXPR) node: WORD: /([a-zA-Z\._][a-zA-Z\._0-9]*)/g /([a-zA-Z\._][a-zA-Z\._0-9]*)/g type:Regex
Loop 1: WORD: /([a-zA-Z\._][a-zA-Z\._0-9]*)/g
 PRE connect nodes, parent:Rank(EXPR) node: INT: /([0-9]+)/g /([0-9]+)/g type:Regex
Loop 2: INT: /([0-9]+)/g
 PRE connect nodes, parent:Rank(EXPR) node: _PIPE: (_ '|') (_ '|') type:*ast.Sequence
 PRE connect nodes, parent:(_ '|') node: undefined _ type:*ast.Ref
 PRE connect nodes, parent:(_ '|') node: undefined '|' type:ast.Str
Loop 3: _PIPE: (_ '|')

While in go
 PRE connect nodes, parent:nil node: nilGRAMMAR{*ast.RankRank(EXPR)} type: *ast.Grammar
 PRE connect nodes, parent:GRAMMAR{*ast.RankRank(EXPR)} node: nilRank(EXPR) type: *ast.Rank
 PRE connect nodes, parent:Rank(EXPR) node: EXPR: Rank(EXPR[0],CHOICE)Rank(EXPR[0],CHOICE) type: *ast.Rank
 PRE connect nodes, parent:Rank(EXPR[0],CHOICE) node: EXPR[0]: (CHOICE _)(CHOICE _) type: *ast.Sequence
 PRE connect nodes, parent:(CHOICE _) node: nilCHOICE type: *ast.Ref
 PRE connect nodes, parent:(CHOICE _) node: nil_ type: *ast.Ref


 DONE rewrote Labels_ and Captures_ using lazy.
      Perhaps also needed for HandlesChildLabel().

FIXED BUG we now only have 24 rules i/o 35.
    first 1 ref missing  is RANGE. 
	  Turns out RANGE is I rule && all I lines account for 11, and 24+11=35
	  ∵ DOING make I rule be part of Grammar.Rules too

