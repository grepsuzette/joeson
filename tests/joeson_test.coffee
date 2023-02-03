require './setup'

# Uses the existing Joeson grammar to parse the grammar RAW_GRAMMAR.
# The resulting object, another parser (which is also a Joeson parser), is
# used to parse its own grammar RAW_GRAMMAR.
# The console output shows some benchmark data.
#
# NOTE: keep this grammar in sync with src/joeson.coffee
# Once we have a compiler, we'll just move this into src/joeson.coffee.

{@trace, setTrace, NODES, GRAMMAR, MACROS, Grammar, HandcompiledRules, Choice, Sequence, Lookahead, Existential, Pattern, Not, Ref, Str, Regex} = require '../src/joeson'
{clazz, colors:{red, blue, cyan, magenta, green, normal, black, white, yellow}} = require('cardamom')
{inspect} = require 'util'
assert = require 'assert'
{pad, escape, toAscii} = require '../lib/helpers'

# environment $TEST must govern which test to run
# environment $TRACE must govern trace options
sTest = process.env.TEST
sTrace = process.env.TRACE
if sTrace?
    aTrace = sTrace.split ","
    localTrace = {
        stack: no,
        loop: no,
        grammar: no,
        filterLine: -1,
        skipSetup: no,
    }
    for opt in aTrace
        switch opt.toLowerCase()
            when "none"
                localTrace.stack = no
                localTrace.loop = no
                localTrace.grammar = no
                localTrace.skipSetup = no
                localTrace.filterLine = -1
            when "stack"
                localTrace.stack = yes
            when "loop"
                localTrace.loop = yes
            when "grammar"
                localTrace.grammar = yes
            when "skipsetup"
                localTrace.skipSetup = yes
            when "all"
                localTrace.stack = yes
                localTrace.loop = yes
                localTrace.grammar = yes
                localTrace.filterLine = -1
                localTrace.skipSetup = no
            when ""
            else
                if opt.indexOf("line=") == 0 || opt.indexOf("filterline=") == 0
                    pos = opt.indexOf("=")
                    if pos <= 0
                        throw "TRACE option: " + yellow(opt) + " requires an =<INT> suffix"
                    else
                        localTrace.filterLine = opt.substr(pos+1)
                else
                    throw "unrecognized TRACE option: " + yellow(opt)
    setTrace localTrace

{o, i, t} = MACROS
QUOTE = "'\\''"
FSLSH = "'/'"
LBRAK = "'['"
RBRAK = "']'"
LCURL = "'{'"
RCURL = "'}'"
RAW_GRAMMAR = [
  o EXPR: [
    o "CHOICE _"
    o "CHOICE": [
      o "_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", (it) -> new Choice it
      o "SEQUENCE": [
        o "UNIT{2,}", (it) -> new Sequence it
        o "UNIT": [
          o "_ LABELED"
          o "LABELED": [
            o "(label:LABEL ':')? &:(DECORATED|PRIMARY)"
            o "DECORATED": [
              o "PRIMARY '?'", (it) -> new Existential it
              o "value:PRIMARY '*' join:(!__ PRIMARY)? @:RANGE?", (it) -> new Pattern it
              o "value:PRIMARY '+' join:(!__ PRIMARY)?", ({value,join}) -> new Pattern value:value, join:join, min:1
              o "value:PRIMARY @:RANGE", (it) -> new Pattern it
              o "'!' PRIMARY", (it) -> new Not it
              o "'(?' expr:EXPR ')' | '?' expr:EXPR", (it) -> new Lookahead it
              i "RANGE": "'{' _ min:INT? _ ',' _ max:INT? _ '}'"
            ]
            o "PRIMARY": [
              o "WORD '(' EXPR ')'", (it) -> new Ref it...
              o "WORD", (it) -> new Ref it
              o "'(' inlineLabel:(WORD ': ')? expr:EXPR ')' ( _ '->' _ code:CODE )?", ({expr, code}) ->
                assert.ok not code?, "code in joeson deprecated"
                return expr
              i CODE: "#{LCURL} (!#{RCURL} (ESC1 | .))* #{RCURL}", (it) -> require('../src/joescript').parse(it.join '')
              o "#{QUOTE} (!#{QUOTE} (ESC1 | .))* #{QUOTE}", (it) -> new Str       it.join ''
              o "#{FSLSH} (!#{FSLSH} (ESC2 | .))* #{FSLSH}", (it) -> new Regex     it.join ''
              o "#{LBRAK} (!#{RBRAK} (ESC2 | .))* #{RBRAK}", (it) -> new Regex "[#{it.join ''}]"
            ]
          ]
        ]
      ]
    ]
  ]
  i LABEL:      "'&' | '@' | WORD"
  i WORD:       "/[a-zA-Z\\._][a-zA-Z\\._0-9]*/"
  i INT:        "/[0-9]+/", (it) -> new Number it
  i _PIPE:      "_ '|'"
  i _:          "(' ' | '\n')*"
  i __:         "(' ' | '\n')+"
  i '.':        "/[\\s\\S]/"
  i ESC1:       "'\\\\' ."
  i ESC2:       "'\\\\' .", (chr) -> '\\'+chr
]

trace = @trace
gmIntention = -> Grammar RAW_GRAMMAR
allFuncs =
    ParseIntention: ->
        # this allows tracing and diffing,
        # it does not do more than compiling the RAW_GRAMMAR
        PARSED_GRAMMAR = Grammar RAW_GRAMMAR
    DebugLabel: ->
        debuglabel = Grammar [
            o In: "l:Br"
            i Br: "'Toy' | 'BZ'"
        ]
        x = debuglabel.parse "Toy"
        if !x.l?
            throw "assert label l"
        else if x.l != "Toy"
            throw "parse error"
        else
            console.log "Test DebugLabel is successful"
    ManyTimes: ->
        start = new Date()
        nbIter = 100
        PARSED_GRAMMAR = Grammar RAW_GRAMMAR
        frecurse = (rule, indent=0, name=undefined) ->
          if rule instanceof Array
            console.log "#{Array(indent*2+1).join ' '}#{red name+":"}" if name
            frecurse r, indent+1 for r in rule
          else if rule instanceof MACROS.o
            [rule, callback] = rule.args
            frecurse rule, indent, name
          else if rule instanceof MACROS.i
            for name, value of rule.args[0]
              frecurse value, indent, name
          else if typeof rule is 'string'
            {result, code} = PARSED_GRAMMAR.parse rule, {debug:no, returnContext:yes}
            console.log "#{Array(indent*2+1).join ' ' \
                        }#{if name? then red(pad({left:(10-indent*2)}, name+':')) else '' \
                        }#{if result? then yellow result else red result} #{white code.peek afterChars:10}"
          else
            for name, r of rule
              frecurse r, indent, name
        for t in [0..nbIter-1]
            frecurse RAW_GRAMMAR
        console.log "Duration for #{nbIter} iterations: #{new Date() - start}ms"
    Squareroot: ->
        gm = Grammar [
            o sqr: "w:word '(' n:int ')'"
            i word: "[a-z]{1,}"
            i "int": "/-?[0-9]{1,}/", (it) -> new Number it
        ]
        x = gm.parse "squareroot(-1)"
        if !x.w? || x.w.join("") != "squareroot"
            throw "expected w label to have value [s,q,u,a,r,e,r,o,o,t], not " + x.w
        else if 0+x.n != -1
            throw "expected n label to be -1, not " + x.n
        else
            console.log "Test Squareroot is successful"
    Choice2: ->
        gm = Grammar [
            o CHOICE: [
                o "_PIPE* SEQUENCE*_PIPE{2,} _PIPE*", (it) ->
                    console.log it
                    return it
                o "SEQUENCE": "WORD _ '_'"
            ]
            i "_PIPE": "_ '|'"
            i WORD: "[A-Z]{1,}"
            i "_": "(' ' | '\n')*"
        ]
        x = gm.parse "CHOICE _"

if sTest?
    if allFuncs[sTest]?
        console.log "------------ Test#{sTest} ----------------------"
        allFuncs[sTest]()
    else if sTest == "-h" # -h to list tests (one per line)
        console.log Object.keys(allFuncs).join("\n")
    else if sTest == "-H" # -H to list tests (comma-separated)
        console.log Object.keys(allFuncs).join(",")
    else
        console.log "#{sTest} is an unrecognized test, available: " + \
            Object.keys(allFuncs).join(",")
else
    console.log "Please set TEST= environment variable, available: " + \
        Object.keys(allFuncs).join(",")


