{Grammar} = require 'joeson'
{red, blue, cyan, magenta, green, normal, black, white, yellow} = require 'joeson/lib/colors'
{clazz} = require 'cardamom'
assert = require 'assert'
_ = require 'underscore'

# Convenience function for letting you use native strings for development.
isWord = (thing) -> thing instanceof Word or typeof thing is 'string'
toString = (thing) -> assert.ok isWord thing; if typeof thing is 'string' then thing else thing.word

# Lexical scope.
LScope = clazz 'LScope', ->
  @VARIABLE = 'variable'
  @PARAMETER = 'parameter'
  init: (@parent) ->
    @variables = []; @parameters = []
    @children = []
    @parent.children.push this if @parent?
  declares: (name) ->
    name = toString name
    return name in @variables or name in @parameters
  isDeclared: (name) ->
    name = toString name
    return true if name in @variables or name in @parameters
    return true if @parent?.isDeclared(name)
    return false
  willDeclare: (name) ->
    name = toString name
    return true if name in @variables or name in @parameters
    return true if _.any @children, (child)->child.willDeclare(name)
    return false
  addVariable: (name, forceDeclaration=no) ->
    name = toString name
    if forceDeclaration
      @variables.push name unless name in @variables
    else
      # This logic decides which scope 'name' belongs to.
      # See docs/coffeescript_lessons/lexical_scoping
      @variables.push name unless @isDeclared(name)
  addParameter: (name) ->
    name = toString name
    @parameters.push name unless name in @parameters
  makeTempVar: (prefix='temp', isParam=no) ->
    # create a temporary variable that is not used in the inner scope
    idx = 0
    loop
      name = "#{prefix}#{idx++}"
      break unless @willDeclare name
    if isParam then @addParameter name else @addVariable name, yes
    return name

# All AST nodes are instances of Node.
Node = clazz 'Node', ->
  walk: ({pre, post, parent}) ->
    # pre, post: (parent, childNode) -> where childNode in parent.children.
    result = pre parent, @ if pre?
    unless result? and result.recurse is no
      children = this.children
      if children?
        for child in children
          if child instanceof Array
            for item in child when item instanceof Node
              item.walk {pre:pre, post:post, parent:@}
          else if child instanceof Node
            child.walk {pre:pre, post:post, parent:@}
    post parent, @ if post?
  # not all nodes will have their own scopes.
  # pass in parentScope if @parent is not available.
  createOwnScope: (parentScope) ->
    parentScope ||= @parent.scope
    assert.ok parentScope, "Where is my parent scope?"
    @scope = @ownScope = new LScope parentScope
  # by default, nodes belong to their parent's scope.
  setScopes: -> @scope ||= @parent.scope
  parameters: null
  variables: null
  # call on the global node to prepare the entire tree.
  prepare: ->
    assert.ok not @prepared, "Node is already prepared."
    @scope = @ownScope = new LScope() if not @scope?
    @walk pre: (parent, node) ->
      node.parent ||= parent
    @walk pre: (parent, node) ->
      node.setScopes()
      node.scope.addParameter param for param in (node.parameters||[])
      node.scope.addVariable _var for _var in (node.variables||[])
      node.prepared = yes

Word = clazz 'Word', Node, ->
  init: (@word) ->
    switch @word
      when 'undefined'
        @_newOverride = Undefined.undefined
      when 'null'
        @_newOverride = Null.null
  toString: -> @word

Block = clazz 'Block', Node, ->
  init: (lines) ->
    @lines = if lines instanceof Array then lines else [lines]
  children$: get: -> @lines
  toString: ->
    (''+line for line in @lines).join '\n'
  toStringWithIndent: ->
    '\n  '+((''+line+';').replace(/\n/g, '\n  ') for line in @lines).join('\n  ')+'\n'

If = clazz 'If', Node, ->
  init: ({@cond, @block, @elseBlock}) ->
    @block = Block @block if @block not instanceof Block
  children$: get: -> [@cond, @block, @elseBlock]
  toString: ->
    if @elseBlock?
      "if(#{@cond}){#{@block}}else{#{@elseBlock}}"
    else
      "if(#{@cond}){#{@block}}"

Unless = ({cond, block, elseBlock}) -> If cond:Not(cond), block:block, elseBlock:elseBlock

For = clazz 'For', Node, ->
  init: ({@label, @block, @own, @keys, @type, @obj, @cond}) ->
  children$: get: -> [@label, @block, @keys, @obj, @cond]
  variables$: get: -> @keys
  toString: -> "for #{@own? and 'own ' or ''}#{@keys.join ','} #{@type} #{@obj} #{@cond? and "when #{@cond} " or ''}{#{@block}}"

While = clazz 'While', Node, ->
  init: ({@label, @cond, @block}) -> @cond ?= true
  children$: get: -> [@label, @cond, @block]
  toString: -> "while(#{@cond}){#{@block}}"

Loop = clazz 'Loop', While, ->
  init: ({@label, @block}) -> @cond = true
  children$: get: -> [@label, @block]

Switch = clazz 'Switch', Node, ->
  init: ({@obj, @cases, @default}) ->
  children$: get: -> [@obj, @cases, @default]
  toString: -> "switch(#{@obj}){#{@cases.join('//')}//else{#{@default}}}"

Try = clazz 'Try', Node, ->
  init: ({@block, @doCatch, @catchVar, @catchBlock, @finally}) ->
  children$: get: -> [@block, @catchVar, @catchBlock, @finally]
  setScopes: ->
    @scope ||= @parent.scope
    if @catchVar? and @catchBlock?
      @catchBlock.createOwnScope(@scope)
      @catchBlock.parameters = [@catchVar]
  toString: -> "try{#{@block}}#{
                @doCatch and "catch(#{@catchVar or ''}){#{@catchBlock}}" or ''}#{
                @finally and "finally{#{@finally}}" or ''}"

Case = clazz 'Case', Node, ->
  init: ({@matches, @block}) ->
  children$: get: -> [@matches, @block]
  toString: -> "when #{@matches.join ','}{#{@block}}"

Operation = clazz 'Operation', Node, ->
  init: ({@left, @not, @op, @right}) ->
  children$: get: -> [@left, @right]
  toString: -> "(#{if @left  then @left+' '  else ''}#{
                   if @not   then 'not '     else ''}#{@op}#{
                   if @right then ' '+@right else ''})"

Not = (it) -> Operation op:'not', right:it

Statement = clazz 'Statement', Node, ->
  init: ({@type, @expr}) ->
  children$: get: -> [@expr]
  toString: -> "#{@type}(#{@expr ? ''});"

Invocation = clazz 'Invocation', Node, ->
  init: ({@func, @params}) ->
    @type = if isWord(@func) and ''+@func is 'new' then 'new' else undefined
  children$: get: -> [@func, @params]
  toString: -> "#{@func}(#{@params.map((p)->"#{p}#{p.splat and '...' or ''}")})"

Assign = clazz 'Assign', Node, ->
  # type: =, +=, -=. *=, /=, ?=, ||= ...
  init: ({@target, type, @op, @value}) ->
    @op = type[...type.length-1] if type?
  children$: get: -> [@target, @value]
  variables$: get: -> if isWord @target then [@target] else null
  toString: -> "#{@target}#{@op or ''}=(#{@value})"

Slice = clazz 'Slice', Node, ->
  init: ({@obj, @range}) ->
  children$: get: -> [@obj, @range]
  toString: -> "#{@obj}[#{@range}]"

Index = clazz 'Index', Node, ->
  init: ({obj, attr, type}) ->
    type ?= if isWord attr then '.' else '['
    if type is '::'
      if attr?
        obj = Index obj:obj, attr:'prototype', type:'.'
      else
        attr = 'prototype'
      type = '.'
    @obj = obj
    @attr = attr
    @type = type
  children$: get: -> [@obj, @attr]
  toString: ->
    close = if @type is '[' then ']' else ''
    "#{@obj}#{@type}#{@attr}#{close}"

Soak = clazz 'Soak', Node, ->
  init: (@obj) ->
  children$: get: -> [@obj]
  toString: -> "(#{@obj})?"

Obj = clazz 'Obj', Node, ->
  init: (@items) ->
  children$: get: -> @items
  toString: -> "{#{if @items? then @items.join ',' else ''}}"

Null = clazz 'Null', Node, ->
  @null = new @(yes)
  init: (construct) ->
    if construct isnt yes
      @_newOverride = Null.null
  value: null
  toString: -> "null"

Undefined = clazz 'Undefined', Node, ->
  @undefined = new @(yes)
  init: (construct) ->
    if construct isnt yes
      @_newOverride = Undefined.undefined
  value: undefined
  toString: -> "undefined"

This = clazz 'This', Node, ->
  init: ->
  toString: -> "@"

Arr = clazz 'Arr', Obj, ->
  children$: get: -> @items
  toString: -> "[#{if @items? then @items.join ',' else ''}]"

Item = clazz 'Item', Node, ->
  init: ({@key, @value}) ->
  children$: get: -> [@key, @value]
  toString: -> @key+(if @value?   then ":(#{@value})"   else '')

Str = clazz 'Str', Node, ->
  init: (@parts) ->
  children$: get: -> @parts
  isStatic: get: -> _.all @parts, (part)->typeof part is 'string'
  toString: ->
    if typeof @parts is 'string'
      '"' + @parts.replace(/"/g, "\\\"") + '"'
    else
      parts = @parts.map (x) ->
        if x instanceof Node
          '#{'+x+'}'
        else
          x.replace /"/g, "\\\""
      '"' + parts.join('') + '"'

Func = clazz 'Func', Node, ->
  init: ({@params, @type, @block}) ->
  children$: get: -> [@params, @block]
  setScopes: ->
    @scope ||= @parent.scope
    if @block?
      if @params?
        @block.parameters = parameters = []
        collect = (thing) ->
          if isWord thing
            parameters.push thing
          else if thing instanceof Obj # Arr is a subclass of Obj
            collect(item.value or item) for item in thing.items
          else if thing instanceof Index
            "pass" # nothing to do for properties
        collect(param) for param in @params
      @block.createOwnScope(@scope)
  toString: -> "(#{if @params then @params.map(
      (p)->"#{p}#{p.splat and '...' or ''}#{p.default and '='+p.default or ''}"
    ).join ',' else ''})#{@type}{#{@block}}"

Range = clazz 'Range', Node, ->
  init: ({@start, @type, @end, @by}) ->
    @by ?= 1
  children$: get: -> [@start, @end, @by]
  toString: -> "Range(#{@start? and "start:#{@start}," or ''}"+
                     "#{@end?   and "end:#{@end},"     or ''}"+
                     "type:'#{@type}', by:#{@by})"

NativeExpression = clazz 'NativeExpression', Node, ->
  init: ({@exprStr}) ->
  toString: -> "`#{@exprStr}`"

Heredoc = clazz 'Heredoc', Node, ->
  init: (@text) ->
  toString: -> "####{@text}###"

Dummy = clazz 'Dummy', Node, ->
  init: (@args) ->
  toString: -> "{#{@args}}"

@NODES = {
  Node, Word, Block, If, For, While, Loop, Switch, Try, Case, Operation,
  Statement, Invocation, Assign, Slice, Index, Soak, Obj, This,
  Null, Undefined,
  Arr, Item, Str, Func, Range, Heredoc, Dummy
}

debugIndent = yes

checkIndent = (ws) ->
  @stack[0].indent ?= '' # set default lazily

  container = @stack[@stack.length-2]
  @log "[In] container (@#{@stack.length-2}:#{container.name}) indent:'#{container.indent}', softline:'#{container.softline}'" if debugIndent
  if container.softline?
    # {
    #   get: -> # begins with a softline
    #   set: ->
    # }
    pIndent = container.softline
  else
    # Get the parent container's indent string
    for i in [@stack.length-3..0] by -1
      if @stack[i].softline? or @stack[i].indent?
        pContainer = @stack[i]
        pIndent = pContainer.softline ? pContainer.indent
        @log "[In] parent pContainer (@#{i}:#{pContainer.name}) indent:'#{pContainer.indent}', softline:'#{pContainer.softline}'" if debugIndent
        break
  # If ws starts with pIndent... valid
  if ws.length > pIndent.length and ws.indexOf(pIndent) is 0
    @log "Setting container.indent to '#{ws}'"
    container.indent = ws
    return container.indent
  null

checkNewline = (ws) ->
  @stack[0].indent ?= '' # set default lazily

  # find the container indent (or softline) on the stack
  for i in [@stack.length-2..0] by -1
    if @stack[i].softline? or @stack[i].indent?
      container = @stack[i]
      break

  containerIndent = container.softline ? container.indent
  isNewline = ws is containerIndent
  @log "[NL] container (@#{i}:#{container.name}) indent:'#{container.indent}', softline:'#{container.softline}', isNewline:'#{isNewline}'" if debugIndent
  return ws if isNewline
  null

# like a newline, but allows additional padding
checkSoftline = (ws) ->
  @stack[0].indent ?= '' # set default lazily

  # find the applicable indent
  container = null
  for i in [@stack.length-2..0] by -1
    if i < @stack.length-2 and @stack[i].softline?
      # a strict ancestor's container's softline acts like an indent.
      # this allows softlines to be shortened only within the same direct container.
      container = @stack[i]
      @log "[SL] (@#{i}:#{container.name}) indent(ignored):'#{container.indent}', **softline**:'#{container.softline}'" if debugIndent
      break
    else if @stack[i].indent?
      container = @stack[i]
      @log "[SL] (@#{i}:#{container.name}) **indent**:'#{container.indent}', softline(ignored):'#{container.softline}'" if debugIndent
      break
  assert.ok container isnt null
  # commit softline ws to container
  if ws.indexOf(container.softline ? container.indent) is 0
    topContainer = @stack[@stack.length-2]
    @log "[SL] Setting topmost container (@#{@stack.length-2}:#{topContainer.name})'s softline to '#{ws}'"
    topContainer.softline = ws
    return ws
  null

checkComma = ({beforeBlanks, beforeWS, afterBlanks, afterWS}) ->
  container = @stack[@stack.length-2]
  container.trailingComma = yes if afterBlanks?.length > 0 # hack for INVOC_IMPL, see _COMMA_NEWLINE
  if afterBlanks.length > 0
    return null if checkSoftline.call(this, afterWS) is null
  else if beforeBlanks.length > 0
    return null if checkSoftline.call(this, beforeWS) is null
  ','

checkCommaNewline = (ws) ->
  @stack[0].indent ?= '' # set default lazily
  container = @stack[@stack.length-2]
  return null if not container.trailingComma
  # Get the parent container's indent string
  for i in [@stack.length-3..0] by -1
    if @stack[i].softline? or @stack[i].indent?
      pContainer = @stack[i]
      pIndent = pContainer.softline ? pContainer.indent
      break
  # If ws starts with pIndent... valid
  if ws.length > pIndent.length and ws.indexOf(pIndent) is 0
    return yes
  null

resetIndent = (ws) ->
  @stack[0].indent ?= '' # set default lazily
  # find any container
  container = @stack[@stack.length-2]
  assert.ok container?
  @log "setting container(=#{container.name}).indent to '#{ws}'"
  container.indent = ws
  return container.indent

@GRAMMAR = Grammar ({o, i, tokens}) -> [
  o                                 "_BLANKLINE* LINES ___", (node) -> node.prepare() unless @options?.rawNodes; node
  i LINES:                          "LINE*_NEWLINE", Block
  i LINE: [
    o HEREDOC:                      "_ '###' !'#' (!'###' .)* '###'", (it) -> Heredoc it.join ''
    o LINEEXPR: [
      # left recursive
      o POSTIF:                     "block:LINEEXPR _IF cond:EXPR", If
      o POSTUNLESS:                 "block:LINEEXPR _UNLESS cond:EXPR", Unless
      o POSTFOR:                    "block:LINEEXPR _FOR own:_OWN? keys:SYMBOL*_COMMA{1,2} type:(_IN|_OF) obj:EXPR (_WHEN cond:EXPR)?", For
      o POSTWHILE:                  "block:LINEEXPR _WHILE cond:EXPR", While
      # rest
      o STMT:                       "type:(_RETURN|_THROW|_BREAK|_CONTINUE) expr:EXPR?", Statement
      o EXPR: [
        o FUNC:                     "params:PARAMS? _ type:('->'|'=>') block:BLOCK?", Func
        i PARAMS:                   "_ '(' (&:PARAM default:(_ '=' LINEEXPR)?)*_COMMA _ ')'"
        i PARAM:                    "&:SYMBOL splat:'...'?
                                    |&:PROPERTY splat:'...'?
                                    |OBJ_EXPL
                                    |ARR_EXPL"
        o RIGHT_RECURSIVE: [
          o INVOC_IMPL:             "func:ASSIGNABLE (? __|OBJ_IMPL_INDENTED) params:(&:EXPR splat:'...'?)+(_COMMA | _COMMA_NEWLINE)", Invocation
          i OBJ_IMPL_INDENTED:      "_INDENT OBJ_IMPL_ITEM+(_COMMA|_NEWLINE)", Obj
          o OBJ_IMPL:               "_INDENT? OBJ_IMPL_ITEM+(_COMMA|_NEWLINE)", Obj
          i OBJ_IMPL_ITEM:          "key:(WORD|STRING) _ ':' value:EXPR", Item
          o ASSIGN:                 "target:ASSIGNABLE _ type:('='|'+='|'-='|'*='|'/='|'?='|'||='|'or='|'and=') value:BLOCKEXPR", Assign
        ]
        o COMPLEX: [
          o IF:                     "_IF cond:EXPR block:BLOCK ((_NEWLINE | _INDENT)? _ELSE elseBlock:BLOCK)?", If
          o UNLESS:                 "_UNLESS cond:EXPR block:BLOCK ((_NEWLINE | _INDENT)? _ELSE elseBlock:BLOCK)?", Unless
          o FOR:                    "_FOR own:_OWN? keys:SYMBOL*_COMMA{1,2} type:(_IN|_OF) obj:EXPR (_WHEN cond:EXPR)? block:BLOCK", For
          o LOOP:                   "_LOOP block:BLOCK", While
          o WHILE:                  "_WHILE cond:EXPR block:BLOCK", While
          o SWITCH:                 "_SWITCH obj:EXPR _INDENT cases:CASE*_NEWLINE default:DEFAULT?", Switch
          i CASE:                   "_WHEN matches:EXPR+_COMMA block:BLOCK", Case
          i DEFAULT:                "_NEWLINE _ELSE BLOCK"
          o TRY:                    "_TRY block:BLOCK
                                     (_NEWLINE? doCatch:_CATCH catchVar:EXPR? catchBlock:BLOCK?)?
                                     (_NEWLINE? _FINALLY finally:BLOCK)?", Try
        ]
        o OP_OPTIMIZATION:          "OP40 _ !(OP00_OP|OP05_OP|OP10_OP|OP20_OP|OP30_OP)"
        o OP00: [
          i OP00_OP:                " '&&' | '||' | '&' | '|' | '^' | _AND | _OR "
          o                         "left:(OP00|OP05) _ op:OP00_OP _SOFTLINE? right:OP05", Operation
          o OP05: [
            i OP05_OP:              " '==' | '!=' | '<=' | '<' | '>=' | '>' | _IS | _ISNT "
            o                       "left:(OP05|OP10) _ op:OP05_OP _SOFTLINE? right:OP10", Operation
            o OP10: [
              i OP10_OP:            " '+' | '-' "
              o                     "left:(OP10|OP20) _ op:OP10_OP _SOFTLINE? right:OP20", Operation
              o OP20: [
                i OP20_OP:          " '*' | '/' | '%' "
                o                   "left:(OP20|OP30) _ op:OP20_OP _SOFTLINE? right:OP30", Operation
                o OP30: [
                  i OP30_OP:        "not:_NOT? op:(_IN|_INSTANCEOF)"
                  o                 "left:(OP30|OP40) _  @:OP30_OP _SOFTLINE? right:OP40", Operation
                  o OP40: [
                    i OP40_OP:      " _NOT | '!' | '~' "
                    o               "_ op:OP40_OP right:OP40", Operation
                    o OP45: [
                      i OP45_OP:    " '?' "
                      o             "left:(OP45|OP50) _ op:OP45_OP _SOFTLINE? right:OP50", Operation
                      o OP50: [
                        i OP50_OP:  " '--' | '++' "
                        o           "left:OPATOM op:OP50_OP", Operation
                        o           "_ op:OP50_OP right:OPATOM", Operation
                        o OPATOM:   "FUNC | RIGHT_RECURSIVE | COMPLEX | ASSIGNABLE"
                      ] # end OP50
                    ] # end OP45
                  ] # end OP40
                ] # end OP30
              ] # end OP20
            ] # end OP10
          ] # end OP05
        ] # end OP00
      ] # end EXPR
    ] # end LINEEXPR
  ] # end LINE

  i ASSIGNABLE: [
    # left recursive
    o SLICE:        "obj:ASSIGNABLE !__ range:RANGE", Slice
    o INDEX0:       "obj:ASSIGNABLE type:'['  attr:LINEEXPR _ ']'", Index
    o INDEX1:       "obj:ASSIGNABLE type:'.'  attr:WORD", Index
    o PROTO:        "obj:ASSIGNABLE type:'::' attr:WORD?", Index
    o INVOC_EXPL:   "func:ASSIGNABLE '(' ___ params:(&:LINEEXPR splat:'...'?)*(_COMMA|_SOFTLINE) ___ ')'", Invocation
    o SOAK:         "ASSIGNABLE '?'", Soak
    # rest
    o TYPEOF: [
      o             "func:_TYPEOF '(' ___ params:LINEEXPR{1,1} ___ ')'", Invocation
      o             "func:_TYPEOF __ params:LINEEXPR{1,1}", Invocation
    ]
    o RANGE:        "_ '[' start:LINEEXPR? _ type:('...'|'..') end:LINEEXPR? _ ']' by:(_BY EXPR)?", Range
    o ARR_EXPL:     "_ '[' _SOFTLINE? (&:LINEEXPR splat:'...'?)*(_COMMA|_SOFTLINE) ___ ']'", Arr
    o OBJ_EXPL:     "_ '{' _SOFTLINE? OBJ_EXPL_ITEM*(_COMMA|_SOFTLINE) ___ '}'", Obj
    i OBJ_EXPL_ITEM: "key:(PROPERTY|WORD|STRING) value:(_ ':' LINEEXPR)?", Item
    o PAREN:        "_ '(' ___ LINEEXPR ___ ')'"
    o PROPERTY:     "_ '@' (WORD|STRING)", (attr) -> Index obj:This(), attr:attr
    o THIS:         "_ '@'", This
    o REGEX:        "_ _FSLASH !__ &:(!_FSLASH !_TERM (ESC2 | .))* _FSLASH <words:1> flags:/[a-zA-Z]*/", Str
    o STRING: [
      o             "_ _TQUOTE  (!_TQUOTE  (ESCSTR | INTERP | .))* _TQUOTE", Str
      o             "_ _TDQUOTE (!_TDQUOTE (ESCSTR | INTERP | .))* _TDQUOTE", Str
      o             "_ _DQUOTE  (!_DQUOTE  (ESCSTR | INTERP | .))* _DQUOTE", Str
      o             "_ _QUOTE   (!_QUOTE   (ESCSTR | .))* _QUOTE",  Str
      i ESCSTR:     "_SLASH .", (it) -> {n:'\n', t:'\t', r:'\r'}[it] or it
      i INTERP:     "'\#{' _BLANKLINE* _RESETINDENT LINEEXPR ___ '}'"
    ]
    o NATIVE:       "_ _BTICK (!_BTICK .)* _BTICK", NativeExpression
    o BOOLEAN:      "_TRUE | _FALSE", (it) -> it is 'true'
    o NUMBER:       "_ <words:1> /-?[0-9]+(\\.[0-9]+)?/", Number
    o SYMBOL:       "_ !_KEYWORD WORD"
  ]

  # BLOCKS:
  i BLOCK: [
    o               "_INDENT LINE*_NEWLINE", Block
    o               "_THEN?  LINE+(_ ';')", Block
  ]
  i BLOCKEXPR:      "_INDENT? EXPR"
  i _INDENT:        "_BLANKLINE+ &:_", checkIndent, skipCache:yes
  i _RESETINDENT:   "_BLANKLINE* &:_", resetIndent, skipCache:yes
  i _NEWLINE: [
    o               "_BLANKLINE+ &:_", checkNewline, skipCache:yes
    o               "_ ';'"
  ], skipCache:yes
  i _SOFTLINE:      "_BLANKLINE+ &:_", checkSoftline, skipCache:yes
  i _COMMA:         "beforeBlanks:_BLANKLINE* beforeWS:_ ','
                      afterBlanks:_BLANKLINE*  afterWS:_", checkComma, skipCache:yes
  i _COMMA_NEWLINE: "_BLANKLINE+ &:_", checkCommaNewline, skipCache:yes

  # TOKENS:
  i WORD:           "_ <words:1> /[a-zA-Z\\$_][a-zA-Z\\$_0-9]*/", Word
  i _KEYWORD:       tokens('if', 'unless', 'else', 'for', 'own', 'in', 'of',
                      'loop', 'while', 'break', 'continue',
                      'switch', 'when', 'return', 'throw', 'then', 'is', 'isnt', 'true', 'false', 'by',
                      'not', 'and', 'or', 'instanceof', 'typeof', 'try', 'catch', 'finally')
  i _BTICK:         "'`'"
  i _QUOTE:         "'\\''"
  i _DQUOTE:        "'\"'"
  i _TQUOTE:        "'\\'\\'\\''"
  i _TDQUOTE:       "'\"\"\"'"
  i _FSLASH:        "'/'"
  i _SLASH:         "'\\\\'"
  i '.':            "<chars:1> /[\\s\\S]/",            skipLog:yes
  i ESC1:           "_SLASH .",                        skipLog:yes
  i ESC2:           "_SLASH .", ((chr) -> '\\'+chr),   skipLog:yes

  # WHITESPACES:
  i _:              "<words:1> /[ ]*/",                skipLog:yes
  i __:             "<words:1> /[ ]+/",                skipLog:yes
  i _TERM:          "_ ('\r\n'|'\n')",                 skipLog:yes
  i _COMMENT:       "_ !HEREDOC '#' (!_TERM .)*",      skipLog:yes
  i _BLANKLINE:     "_ _COMMENT? _TERM",               skipLog:yes
  i ___:            "_BLANKLINE* _",                   skipLog:yes

]
# ENDGRAMMAR
