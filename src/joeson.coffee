###
JoeSon Parser
Jae Kwon 2012

Unfortunately, there is an issue with with the implementation of Joeson where "transient" cached values
like those derived from a loopify iteration, that do not get wiped out of the cache between iterations.
What we want is a way to "taint" the cache results, and wipe all the tainted results...
We could alternatively wipe out all cache items for a given position, but this proved to be
significantly slower.

Either figure out a way to taint results and clear tainted values, in a performant fasion
while maintaining readability of the joeson code, or
just admit that the current implementation is imperfect, and limit grammar usage.

- June, 2012

###

# to configure the trace for the tests, use $TRACE envvar, see joeson_test.coffee
@trace = trace =
  #filterLine: 299
  stack:      no
  loop:       no
  skipSetup:  no
  grammar:    no

{clazz, colors:{red, blue, cyan, magenta, green, normal, black, white, yellow}} = require('cardamom')

@setTrace = setTrace = (traceOptions) -> trace = traceOptions
boldblack = (s) -> black(s, yes)
bold = (s) -> "[1m#{s}[0m"
{inspect} = require 'util'
assert = require 'assert'
{CodeStream} = require './codestream'
Node = require('./node').createNodeClazz('GrammarNode')
{pad, escape:_escape} = require '../lib/helpers'
escape = (text, options) -> _escape ''+text, options

if no
  counter = 0
  total = {}
  timeStack = []
  timeStart = (name) ->
    now = new Date()
    lastStackItem = timeStack[timeStack.length-1]
    if lastStackItem
      lastStackItem.accum += now - lastStackItem.start
    timeStack.push name:name, start:now, accum:0
  timeEnd = (name) ->
    now = new Date()
    {name, start, accum} = timeStack.pop()
    #assert.equal name, name
    t = total[name]
    elapsed = now - start + accum
    if t?
      t.time += elapsed
      t.count += 1
      t.avg = t.time / t.count
    else
      total[name] = name:name, time:elapsed, count:1, avg:NaN
    lastStackItem = timeStack[timeStack.length-1]
    if lastStackItem
      lastStackItem.start = now

    if (++counter)%100000 is 0
      values = Object.values total
      timeTotals = 0
      console.log values.sortBy((x) -> x.time).map((x) -> timeTotals += x.time; {name:x.name, time:x.time, c:x.count})
      console.log "\n#{timeTotals}\n\n"

_loopStack = [] # trace stack

newFrame = (pos, id) -> {result:undefined, pos:pos, endPos:undefined, id:id, loopStage:undefined, wipemask:undefined, param:undefined}

cacheSet = (frame, result, endPos) ->
  frame.result = result
  frame.endPos = endPos

# aka '$' in parse functions
@ParseContext = ParseContext = clazz 'ParseContext', ->

  # code:       CodeStream instance
  # grammar:    Grammar instance
  # opts:       Parse-time options, accessible from grammar callback functions
  init: ({@code, @grammar, @opts}={}) ->
    @stack = new Array(1024)
    @stackLength = 0
    # { pos:{ (node.id):{id,result,pos,endPos,stage,...(same object as in stack)}... } }
    @frames = (new Array(@grammar.numRules) for i in [0...@code.text.length+1]) # include EOF
    @counter = 0

  log: (message) ->
    unless @skipLog
      line = @code.line
      return if trace.filterLine? and trace.filterLine>-1 and line isnt trace.filterLine
      codeSgmnt = "#{ white ''+line+','+@code.col \
                }\t#{ boldblack pad right:5, (p=escape(@code.peek beforeChars:5))[p.length-5...] \
                  }#{ green pad left:20, (p=escape(@code.peek afterChars:20))[0...20] \
                  }#{ if @code.pos+20 < @code.text.length
                        boldblack '>'
                      else
                        boldblack ']'}"
      console.log "#{codeSgmnt} #{cyan Array(@stackLength).join '| '}#{message}"

  stackPeek: (skip=0) -> @stack[@stackLength-1-skip]
  stackPush: (node) -> @stack[@stackLength++] = @getFrame(node)
  stackPop: (node) -> --@stackLength

  getFrame: (node) ->
    id = node.id
    pos = @code.pos
    posFrames = @frames[pos]
    if not (frame=posFrames[id])?
      return posFrames[id] = newFrame pos, id
    else return frame

  wipeWith: (frame, makeStash=yes) ->
    timeStart? 'wipewith'
    assert.ok frame.wipemask?, "Need frame.wipemask to know what to wipe"
    stash = new Array(@grammar.numRules) if makeStash
    stashCount = 0
    pos = frame.pos
    posFrames = @frames[pos]
    for wipe, i in frame.wipemask when wipe
      stash[i] = posFrames[i] if makeStash
      posFrames[i] = undefined
      stashCount++
    stash?.count = stashCount
    timeEnd? 'wipewith'
    return stash

  restoreWith: (stash) ->
    timeStart? 'restorewith'
    stashCount = stash.count
    for frame, i in stash when frame
      @frames[frame.pos][i] = frame
      stashCount--
      break if stashCount is 0
    timeEnd? 'restorewith'
    return

# typeof replacement
toType = (obj) ->
    ({}).toString.call(obj).match(/\s([a-zA-Z]+)/)[1].toLowerCase()

# so that it doesn't log "object" for everything
showtype = (result) ->
    if result?
        # if true # show as go types, for diff_go_vs_coffee
            # feel free to edit if notice discrepancies
            fields = Object.keys(result)
            # when result is like { expr: xx, _origin: ..}, just have result = xx
            result = result["expr"] if ("expr" in fields and fields.length == 1) or ("expr" in fields and "_origin" in fields and fields.length == 2)
            return "<NativeUndefined>" if result is undefined
            return "*joeson.NativeArray" if Array.isArray(result)
            return "joeson.Grammar" if result.rank?
            return "*joeson.ref" if result.ref?
            return "joeson.str" if result.str?
            return "*joeson.pattern" if result.value?
            return "*joeson.regex" if result.reStr?
            return "*joeson.sequence" if result.sequence?
            return "*joeson.rank" if result?.contentString?().indexOf(blue("Rank(")) == 0
            return "*joeson.choice" if result.choices?
            return "Exis|Not" if result.it?
            return "joeson.lookahead" if result.expr?
            return "joeson.NativeString" if toType(result) == "string"
            return "joeson.NativeInt" if toType(result) is "number"
            # return toType result
            if typeof result is "object"
                s = ""
                first = true
                for key in Object.keys(result).sort() # sort keys (we sort them in golang too)
                    s += ", " if !first
                    s += key + ":" + if result[key] == undefined then "<NativeUndefined>" else result[key]
                    first = false
                return "NativeMap{" + s + "} " + cyan "joeson.NativeMap"
                # return "object keys:{" + Object.keys(result) + "}"
            else
                return "Unknown"
        # else # show original js types
        #     return "array" if Array.isArray(result)
        #     return "Grammar" if result.rank?
        #     return "Ref" if result.ref?
        #     return "Str" if result.str?
        #     return "Regex" if result.reStr?
        #     return "Sequence" if result.sequence?
        #     return "Rank" if result?.contentString?().indexOf(blue("Rank(")) == 0
        #     return "Choic" if result.choices?
        #     return "Exis|Not" if result.it?
        #     return "Lookahead" if result.expr?
        #     return "STRING!" if toType(result) == "string"
        #     return "number" if toType(result) == "number"
        #     if typeof result is "object"
        #         # return "object keys:{" + Object.keys(result) + "}"
        #         s = ""
        #         for key, v in result
        #             s = s + key + ","
        #         return "object keys:{" + s + "}"
        #         # return "object keys:" + Object.keys(result) + "."+ toType(result)
        #     else
        #         return "Unknown"

# used in $.log to facilitate conformance of outputs of the coffee and golang impl.
showcontent = (result) ->
    if result == null
        "nil"
    else
        fields = Object.keys(result)
        # when result is like { expr: xx, _origin: ..}, just have result = xx
        result = result["expr"] if ("expr" in fields and fields.length == 1) or ("expr" in fields and "_origin" in fields and fields.length == 2)
        if Array.isArray(result) # so output is similar to golang
            return "[" + result + "] " + cyan showtype result
        else if (result + "") == "[object Object]" # ditto
            return (cyan showtype result).slice(5, -4) # strip outside ansi sequence for cyan
        else # but, in general
            return result + " " + cyan showtype result

###
  In addition to the attributes defined by subclasses,
    the following attributes exist for all nodes.
  node.rule = The topmost node of a rule.
  node.rule = rule # sometimes true.
  node.name = name of the rule, if this is @rule.
###
@GNode = GNode = clazz 'GNode', Node, ->

  @optionKeys = ['skipLog', 'skipCache', 'cb']

  @$stack = (fn) -> ($) ->
    $.stackPush this
    timeStart? @name
    pos = $.code.pos
    result = fn.call this, $
    timeEnd? @name
    $.stackPop this
    return result

  @$loopify = (fn) -> ($) ->
    # STACK TRACE
    # $.log "#{blue '*'} #{this} #{showtype this} #{boldblack $.counter}" if trace.stack
    $.log "#{blue '*'} #{this} #{boldblack $.counter}" if trace.stack

    if @skipCache
      result = fn.call this, $
      $.log "#{cyan "`->:"} #{escape result} #{boldblack typeof result}" if trace.stack
      return result

    frame = $.getFrame this
    pos = startPos = $.code.pos
    key = "#{@name}@#{pos}"

    switch frame.loopStage
      when 0, undefined # non-recursive (so far)

        # The only time a cache hit will simply return is when loopStage is 0
        if frame.endPos?
          # $.log "#{cyan "`-hit:"} #{if frame.result == null then "nil" else (frame.result + " " + cyan(showtype(frame.result)))}" if trace.stack
          $.log "#{cyan "`-hit:"} #{showcontent frame.result}" if trace.stack
          $.code.pos = frame.endPos
          return frame.result

        frame.loopStage = 1
        cacheSet frame, null
        result = fn.call this, $

        switch frame.loopStage
          when 1 # non-recursive (done)
            frame.loopStage = 0
            cacheSet frame, result, $.code.pos
            # $.log "#{cyan "`-set:"} #{if result == null then "nil" else (result + " " + cyan showtype result)}" if trace.stack
            $.log "#{cyan "`-set:"} #{showcontent result}" if trace.stack
            return result

          when 2 # recursion detected by subroutine above
            if result is null
              $.log "#{yellow "`--- loop null ---"} " if trace.stack
              frame.loopStage = 0
              #cacheSet frame, null # already null
              return result
            else
              frame.loopStage = 3
              if trace.loop and (not trace.filterLine? or
                                 $.code.line is trace.filterLine)
                line = $.code.line
                _loopStack.push(@name)
                # console.log  "#{ (switch line%6
                #                     when 0 then blue
                #                     when 1 then cyan
                #                     when 2 then white
                #                     when 3 then yellow
                #                     when 4 then red
                #                     when 5 then magenta)('@'+line) \
                #            }\t#{ red (frame.id for frame in $.stack[...$.stackLength]) \
                #           } - #{ _loopStack \
                #           } - #{ yellow escape ''+result \
                #            }: #{ blue escape $.code.peek beforeChars:10, afterChars:10 }"
              timeStart? 'loopiteration'
              while result isnt null
                assert.ok frame.wipemask?, "where's my wipemask"
                bestStash = $.wipeWith frame, yes
                bestResult = result
                bestEndPos = $.code.pos
                cacheSet frame, bestResult, bestEndPos
                $.log "#{yellow "|`--- loop iteration ---"} #{frame}" if trace.stack
                $.code.pos = startPos
                result = fn.call this, $
                break unless $.code.pos > bestEndPos
              timeEnd? 'loopiteration'

              if trace.loop
                _loopStack.pop()

              $.wipeWith frame, no
              $.restoreWith bestStash
              $.code.pos = bestEndPos
              $.log "#{yellow "`--- loop done! ---"} best result: #{escape bestResult}" if trace.stack
              # Step 4: return best result, which will get cached
              frame.loopStage = 0
              return bestResult

          else throw new Error "Unexpected stage #{stages[pos]}"

      when 1,2,3
        if frame.loopStage is 1
          frame.loopStage = 2 # recursion detected

        timeStart? 'wipemask'
        # Step 1: Collect wipemask so we can wipe the frames later.
        $.log "#{yellow "`-base:"} #{escape frame.result} #{boldblack typeof frame.result}" if trace.stack
        frame.wipemask ?= new Array($.grammar.numRules)
        for i in [$.stackLength-2..0] by -1
          i_frame = $.stack[i]
          assert.ok i_frame.pos <= startPos
          break if i_frame.pos < startPos
          break if i_frame.id is @id
          frame.wipemask[i_frame.id] = yes
        timeEnd? 'wipemask'

        # Step 2: Return whatever was cacheSet.
        $.code.pos = frame.endPos if frame.endPos?
        return frame.result

      else throw new Error "Unexpected stage #{stage} (B)"

  @$prepareResult = (fn) -> ($) ->
    $.counter++
    result = fn.call this, $
    if result isnt null
      # handle labels for standalone nodes
      if @label? and not @parent?.handlesChildLabel
        # syntax proposal:
        # result = ( it <- (it={})[@label] = result )
        result = ( (it={})[@label] = result; it )
      start = $.stackPeek().pos
      end = $.code.pos
      _origin =
        code:   $.code.text
        start:  line:$.code.posToLine(start), col:$.code.posToLine(start), pos: start
        end:    line:$.code.posToLine(end),   col:$.code.posToLine(end),   pos: end
      if @cb?
        if result instanceof Object
          Object.defineProperty result, '_origin', value:_origin, enumerable:no, writable:yes
        result = @cb.call this, result, $
      if result instanceof Object # set it again
        Object.defineProperty result, '_origin', value:_origin, enumerable:no, writable:yes
    return result

  @$wrap = (fn) ->
    # optimizations...
    wrapped1 = @$stack @$loopify @$prepareResult fn
    wrapped2 = @$prepareResult fn
    ($) ->
      if this is @rule
        @parse = parse = wrapped1.bind(this)
      else if @label? and not @parent?.handlesChildLabel or @cb?
        @parse = parse = wrapped2.bind(this)
      else
        @parse = parse = fn.bind(this)
      return parse($)

  # TODO this adds thousands of iterations when 
  #  grammar is parsed, check it's necessary
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}

  capture:   yes
  labels$:   get: -> @_labels ?= (if @label then [@label] else [])
  captures$: get: -> @_captures ?= (if @capture then [this] else [])

  # called after all its children have been prepared.
  # don't put logic in here, too easy to forget to call super.
  prepare: ->

  toString: ->
    "#{ if this is @rule
          red(@name+': ')
        else if @label?
          cyan(@label+':')
        else ''
    }#{ @contentString() }"

  include: (name, rule) ->
    @rules ?= {}
    #assert.ok name?, "Rule needs a name: #{rule}"
    #assert.ok rule instanceof GNode, "Invalid rule with name #{name}: #{rule} (#{rule.constructor.name})"
    #assert.ok not @rules[name]?, "Duplicate name #{name}"
    rule.name = name if not rule.name?
    @rules[name] = rule

  # find a parent in the ancestry chain that satisfies condition
  findParent: (condition) ->
    parent = @parent
    loop
      return parent if condition parent
      parent = parent.parent

@Yes = Yes = clazz 'Yes', GNode, ->
  parse: @$wrap ($) -> yes

@Choice = Choice = clazz 'Choice', GNode, ->
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    choices:    {type:[type:GNode]}
  init: (@choices=[]) ->
  prepare: ->
    @capture = @choices.every (choice) -> choice.capture
  parse: @$wrap ($) ->
    for choice in @choices
      pos = $.code.pos
      result = choice.parse $
      $.code.pos = pos if result is null
      if result isnt null
        return result
    return null
  contentString: -> blue("(")+(@choices.join blue(' | '))+blue(")")

@Rank = Rank = clazz 'Rank', Choice, ->
  @fromLines = (name, lines) ->
    rank = Rank name
    for line, idx in lines
      if line instanceof OLine
        choice = line.toRule rank, index:rank.choices.length
        rank.choices.push choice
      else if line instanceof ILine
        for own name, rule of line.toRules()
          rank.include name, rule
      else if line instanceof Object and idx is lines.length-1
        assert.ok GNode.optionKeys.intersect(Object.keys(line)).length > 0,
          "Invalid options? #{line.constructor.name}"
        Object.merge rank, line
      else
        throw new Error "Unknown line type, expected 'o' or 'i' line, got '#{line}' (#{typeof line})"
    rank
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    choices:    {type:[type:GNode]}
  init: (@name, @choices=[], includes={}) ->
    @rules = {}
    for choice, i in @choices
      @choices.push choice
    for name, rule of includes
      @include name, rule
    # return XXX why is it faster w/o a return statement??
  contentString: -> blue("Rank(")+(@choices.map((c)->red(c.name)).join blue(','))+blue(")")

@Sequence = Sequence = clazz 'Sequence', GNode, ->
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    sequence:   {type:[type:GNode]}
  handlesChildLabel: yes
  init:       (@sequence) ->
  labels$:    get: -> @_labels ?= ((child.labels for child in @sequence).flatten())
  captures$:  get: -> @_captures ?= (child.captures for child in @sequence).flatten()
  type$:      get: ->
    @_type?=(
      if @labels.length is 0
        if @captures.length > 1 then 'array' else 'single'
      else
        'object'
    )
  parse: @$wrap ($) ->
    switch @type
      when 'array'
        results = []
        for child in @sequence
          res = child.parse $
          if res is null
            return null
          results.push res if child.capture
        return results
      when 'single'
        result = undefined
        for child in @sequence
          res = child.parse $
          if res is null
            return null
          result = res if child.capture
        return result
      when 'object'
        results = undefined
        # results[label] = undefined for label in @labels
        for child in @sequence
          res = child.parse $
          if res is null
            return null
          if child.label is '&'
            results = if results? then Object.merge res, results else res
          else if child.label is '@'
            results = if results? then Object.merge results, res else res
          else if child.label?
            (results?={})[child.label] = res
        return results
      else
        throw new Error "Unexpected type #{@type}"
    throw new Error

  contentString: ->
    labeledStrs = for node in @sequence
      ''+node
    blue("(")+(labeledStrs.join ' ')+blue(")")

@Lookahead = Lookahead = clazz 'Lookahead', GNode, ->
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    expr:       {type:GNode}
  capture: no
  init: ({@expr}) ->
  parse: @$wrap ($) ->
    pos = $.code.pos
    result = @expr.parse $
    $.code.pos = pos
    result
  contentString: -> "#{blue "(?"}#{@expr}#{blue ")"}"

@Existential = Existential = clazz 'Existential', GNode, ->
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    it:         {type:GNode}
  handlesChildLabel$: get: -> @parent?.handlesChildLabel
  init: (@it) ->
  prepare: ->
    # console.log "Existential.prepare " + @contentString() + " parent?parent?:" + @.parent?.parent?.contentString()
    labels   = if @label? and @label not in ['@','&'] then [@label] else @it.labels
    @label   ?= '@' if labels.length > 0
    captures  = @it.captures
    @capture  = captures?.length > 0
    # some strangeness in overwritting getter funcs...
    # they don't become available right away. wtf?
    @labels   = labels
    @captures = captures
  parse: @$wrap ($) ->
    pos = $.code.pos
    result = @it.parse $
    $.code.pos = pos if result is null
    return result ? undefined
  contentString: -> '' + @it + blue("?")

@Pattern = Pattern = clazz 'Pattern', GNode, ->
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    value:      {type:GNode}
    join:       {type:GNode}
  init: ({@value, @join, @min, @max}) ->
    @capture = @value.capture
  parse: @$wrap ($) ->
    matches = []
    pos = $.code.pos
    resV = @value.parse $
    if resV is null
      $.code.pos = pos
      if @min? and @min > 0
        return null
      return []
    matches.push resV
    loop
      pos2 = $.code.pos
      if @join?
        resJ = @join.parse $
        # return null to revert pos
        if resJ is null
          $.code.pos = pos2
          break
      resV = @value.parse $
      # return null to revert pos
      if resV is null
        $.code.pos = pos2
        break
      matches.push resV
      if @max? and matches.length >= @max
        break
    if @min? and @min > matches.length
      $.code.pos = pos
      return null
    return matches
  contentString: ->
    "#{@value}#{cyan "*"}#{@join||''}#{cyan if @min? or @max? then "{#{@min||''},#{@max||''}}" else ''}"

@Not = Not = clazz 'Not', GNode, ->
  @defineChildren
    rules:      {type:{key:undefined,value:{type:GNode}}}
    it:         {type:GNode}
  capture: no
  init: (@it) ->
  parse: @$wrap ($) ->
    pos = $.code.pos
    res = @it.parse $
    $.code.pos = pos
    if res isnt null
      return null
    else
      return undefined
  contentString: -> "#{yellow '!'}#{@it}"

@Ref = Ref = clazz 'Ref', GNode, ->
  # note: @ref because @name is reserved.
  init: (@ref, @param) ->
    @capture = no if @ref[0] is '_'
  labels$: get: ->
    @_labels ?=
      if @label is '@' then @grammar.rules[@ref].labels
      else if @label   then [@label]
      else                  []
  parse: @$wrap ($) ->
    node = @grammar.rules[@ref]
    throw Error "Unknown reference #{@ref}" if not node?
    $.stackPeek().param = @param
    return node.parse $
  contentString: -> red(@ref)

@Str = Str = clazz 'Str', GNode, ->
  capture: no
  init: (@str) ->
  parse: @$wrap ($) -> $.code.match string:@str
  contentString: -> green("'#{escape @str}'")

@Regex = Regex = clazz 'Regex', GNode, ->
  init: (@reStr) ->
    if typeof @reStr isnt 'string'
      throw Error "Regex node expected a string but got: #{@reStr}"
    @re = RegExp '('+@reStr+')', 'g' # TODO document why http://blog.stevenlevithan.com/archives/fixing-javascript-regexp
  parse: @$wrap ($) -> $.code.match regex:@re
  contentString: -> magenta((''+@re).replace("\t", "\\t"))

# Main external access.
# I dunno if Grammar should be a GNode or not. It
# might come in handy when embedding grammars
# in some glue language.
@Grammar = Grammar = clazz 'Grammar', GNode, ->

  @defineChildren rank: {type:Rank}

  init: (rank) ->
    rank = rank(MACROS) if typeof rank is 'function'
    @rank = Rank.fromLines "__grammar__", rank if rank instanceof Array
    @rules = {}
    @numRules = 0
    @id2Rule = {} # slow lookup for debugging...

    # TODO refactor into translation passes.
    # Merge Choices with just a single choice.
    # @walk
    pre: ({child:node, parent, desc, key, index}) =>
        if node instanceof Choice and node.choices.length is 1
          # Merge label
          node.choices[0].label ?= node.label
          # Merge included rules
          Object.merge (node.choices[0].rules?={}), node.rules if node.rules?
          # Replace with grandchild
          if index?
            parent[key][index] = node.choices[0]
          else
            parent[key] = node.choices[0]

    # Connect all the nodes and collect dereferences into @rules
    @walk
      pre: ({child:node, parent}) =>
        # sanity check
        if node.parent? and node isnt node.rule
          throw Error 'Grammar tree should be a DAG, nodes should not be referenced more than once.'
        node.grammar = this
        node.parent = parent
        # inline rules are special
        if node.inlineLabel?
          throw "assert this is never called"
          node.rule = node
          parent.rule.include node.inlineLabel, node
        # set node.rule, the root node for this rule
        else
          node.rule ||= parent?.rule
      post: ({child:node, parent}) =>
        if node is node.rule
          @rules[node.name] = node
          node.id = @numRules++
          @id2Rule[node.id] = node
          if trace.loop # print out id->rulename for convenience
            console.log "Loop #{red node.id}:\t#{node}"

    # Prepare all the nodes, child first.
    @walk post: ({child:node, parent}) ->
        node.prepare()

    if trace.grammar
        @printRules()

  # MAIN GRAMMAR PARSE FUNCTION
  parse$: (code, opts = {}) ->
    opts.returnContext ?= no
    assert.ok code, "Parser wants code"
    code = CodeStream code if code not instanceof CodeStream
    $ = ParseContext {code, grammar:this, opts}

    # temporarily enable stack tracing
    if opts?.debug
      oldTrace = Object.clone trace
      trace.stack = yes

    # parse
    $.result = @rank.parse $
    $.result?.code = code

    # undo temprary stack tracing
    if opts?.debug
      trace = oldTrace

    # if parse is incomplete, compute error message
    if $.code.pos isnt $.code.text.length

      # find the maximum parsed entity
      maxAttempt = $.code.pos
      maxSuccess = $.code.pos
      # TODO for x, i in something from index to index by skip
      for pos in [$.code.pos...$.frames.length-1]
        posFrames = $.frames[pos]
        for frame, id in posFrames
          if frame
            maxAttempt = pos
            if frame.result isnt null
              maxSuccess = pos
              break
      parseError = new Error "Error parsing at char:#{maxSuccess}=(line:#{$.code.posToLine(maxSuccess)},col:#{$.code.posToCol(maxSuccess)})."
      parseError.details =
        "#{green 'OK'}/#{yellow 'Parsing'}/#{red 'Suspect'}/#{white 'Unknown'})\n\n#{
            green  $.code.peek beforeLines:2
        }#{ yellow $.code.peek afterChars:(maxSuccess-$.code.pos)
        }#{ $.code.pos = maxSuccess; red $.code.peek afterChars:(maxAttempt-$.code.pos)
        }#{ $.code.pos = maxAttempt; white $.code.peek afterLines:2}\n"
      throw parseError

    # return the resulting parsed nodes, or the whole context if specified.
    if opts.returnContext
      return $
    else
      return $.result

  contentString: -> magenta('GRAMMAR{') + @rank+magenta('}')

  printRules: ->
    console.log "+--------------- Grammar.Debug() ----------------------------------"
    console.log "| name         : " + bold(@.name)
    console.log "| contentString: " + @.contentString()
    console.log "| rules        : " + @.numRules
    console.log "| "
    console.log "| ",
        pad(left:14, "key/name"),
        pad(left:3, "id"),
        pad(left:20, "type"),
        pad(left:3, "cap"),
        pad(left:7, "label"),
        pad(left:21, "labels()"),
        pad(left:16, "parent.name"),
        pad(left:30, "contentString"),
    console.log "|   -------------------------------------------------------------------------------------"
    for k,v of @rules
        console.log "|  ",
            pad(left:14, k),
            pad(left:3, v.id),
            pad(left:20, showtype(v)),
            pad(left:3,
                if v.capture == true
                    "y"
                else "n"
            ),
            pad(left:7, if v.label?
                v.label
            else ""),
            pad(left:21, v.labels.join ","),
            pad(left:16, v.parent?.name),
            pad(left:30, v.contentString())
    console.log "| "


Line = clazz 'Line', ->
  init: (@args...) ->
  # name:       The final and correct name for this rule
  # rule:       A rule-like object
  # parentRule: The actual parent Rule instance
  # attrs:      {cb,...}, extends the result
  # opts:       Parse time options
  getRule: (name, rule, parentRule, attrs) ->
    if typeof rule is 'string'
      try
        # HACK: temporarily halt trace
        oldTrace = trace
        trace = stack:no, loop:no if trace.skipSetup
        rule = GRAMMAR.parse rule
        trace = oldTrace
      catch err
        console.log "Error in rule #{name}: #{rule}"
        console.log err.stack
        # TODO force debug output
        GRAMMAR.parse rule
    else if rule instanceof Array
      rule = Rank.fromLines name, rule
    else if rule instanceof OLine
      rule = rule.toRule parentRule, name:name
    assert.ok not rule.rule? or rule.rule is rule
    rule.rule = rule
    assert.ok not rule.name? or rule.name is name
    rule.name = name
    Object.merge rule, attrs if attrs?
    rule

  # returns {rule:rule, attrs:{cb,skipCache,skipLog,...}}
  getArgs: ->
    [rule, rest...] = @args
    _a_ = rule:rule, attrs:{}
    for own key, value of rule
      if key in GNode.optionKeys
        _a_.attrs[key] = value
        delete rule[key]
    for next in rest
      if next instanceof Function
        _a_.attrs.cb = next
      else
        Object.merge _a_.attrs, next
    _a_
  toString: ->
    "#{@type} #{@args.join ','}"

ILine = clazz 'ILine', Line, ->
  type: 'i'
  toRules: (parentRule) ->
    {rule, attrs} = @getArgs()
    rules = {}
    # for an ILine, rule is an object of {"NAME":rule}
    for own name, _rule of rule
      rules[name] = @getRule name, _rule, parentRule, attrs
    rules

OLine = clazz 'OLine', Line, ->
  type: 'o'
  toRule: (parentRule, {index,name}) ->
    {rule, attrs} = @getArgs()
    # figure out the name for this rule
    if not name and
      typeof rule isnt 'string' and
      rule not instanceof Array and
      rule not instanceof GNode
        # NAME: rule
        assert.ok Object.keys(rule).length is 1, "Named rule should only have one key-value pair"
        name = Object.keys(rule)[0]
        rule = rule[name]
    else if not name? and index? and parentRule?
      name = parentRule.name + "[#{index}]"
    else if not name?
      throw new Error "Name undefined for 'o' line"
    rule = @getRule name, rule, parentRule, attrs
    rule.parent = parentRule
    rule.index = index
    rule

@MACROS = MACROS =
  # Any rule node, possibly part of a Rank node
  o: OLine
  # Include line... Not included in the Rank order
  i: ILine
  # Helper for declaring tokens
  tokens: (tokens...) ->
    cb = tokens.pop() if typeof tokens[tokens.length-1] is 'function'
    regexAll = Regex("[ ]*(#{tokens.join('|')})([^a-zA-Z\\$_0-9]|$)")
    for token in tokens
      name = '_'+token.toUpperCase()
      # HACK: temporarily halt trace
      oldTrace = trace
      trace = stack:no
      rule = GRAMMAR.parse "/[ ]*/ &:'#{token}' !/[a-zA-Z\\$_0-9]/"
      trace = oldTrace
      rule.rule = rule
      rule.skipLog = yes
      rule.skipCache = yes
      rule.cb = cb if cb?
      regexAll.include name, rule
    OLine regexAll
  # Helper for clazz construction in callbacks
  make: (clazz, options=undefined) -> (it, $) -> new clazz it, options

C  = -> Choice (x for x in arguments)
E  = -> Existential arguments...
L  = (label, node) -> node.label = label; node
La = -> Lookahead arguments...
N  = -> Not arguments...
P  = (value, join, min, max) -> Pattern value:value, join:join, min:min, max:max
R  = -> Ref arguments...
Re = -> Regex arguments...
S  = -> Sequence (x for x in arguments)
St = -> Str arguments...
{o, i, tokens}  = MACROS

# Don't worry, this is just the intermediate hand-compiled form of the grammar you can actually understand,
# located currently in tests/joeson.coffee. Look at that instead, and keep the two in sync until they get merged later.
@HandcompiledRules = HandcompiledRules = [
  o EXPR: [
    o S(R("CHOICE"), R("_"))
    o "CHOICE": [
      o S(P(R("_PIPE")), P(R("SEQUENCE"),R("_PIPE"),2), P(R("_PIPE"))), (it) -> new Choice it
      o "SEQUENCE": [
        o P(R("UNIT"),null,2), (it) -> new Sequence it
        o "UNIT": [
          o S(R("_"), R("LABELED"))
          o "LABELED": [
            o S(E(S(L("label",R("LABEL")), St(':'))), L('&',C(R("DECORATED"),R("PRIMARY"))))
            o "DECORATED": [
              o S(R("PRIMARY"), St('?')), (it) -> new Existential it
              o S(L("value",R("PRIMARY")), St('*'), L("join",E(S(N(R("__")), R("PRIMARY")))), L("@",E(R("RANGE")))), (it) -> new Pattern it
              o S(L("value",R("PRIMARY")), St('+'), L("join",E(S(N(R("__")), R("PRIMARY"))))), ({value,join}) -> new Pattern value:value, join:join, min:1
              o S(L("value",R("PRIMARY")), L("@",R("RANGE"))), (it) -> new Pattern it
              o S(St('!'), R("PRIMARY")), (it) -> new Not it
              o C(S(St('(?'), L("expr",R("EXPR")), St(')')), S(St('?'), L("expr",R("EXPR")))), (it) -> new Lookahead it
              i "RANGE": o S(St('{'), R("_"), L("min",E(R("INT"))), R("_"), St(','), R("_"), L("max",E(R("INT"))), R("_"), St('}'))
            ]
            o "PRIMARY": [
              o S(R("WORD"), St('('), R("EXPR"), St(')')), (it) -> new Ref it...
              o R("WORD"), (it) -> new Ref it
              o S(St('('), L("inlineLabel",E(S(R('WORD'), St(': ')))), L("expr",R("EXPR")), St(')'), E(S(R('_'), St('->'), R('_'), L("code",R("CODE"))))), ({expr, code}) ->
                assert.ok not code?, "code in joeson deprecated"
                return expr
              i "CODE": o S(St("{"), P(S(N(St("}")), C(R("ESC1"), R(".")))), St("}")), (it) -> require('./joescript').parse(it.join '')
              o S(St("'"), P(S(N(St("'")), C(R("ESC1"), R(".")))), St("'")), (it) -> new Str       it.join ''
              o S(St("/"), P(S(N(St("/")), C(R("ESC2"), R(".")))), St("/")), (it) -> new Regex     it.join ''
              o S(St("["), P(S(N(St("]")), C(R("ESC2"), R(".")))), St("]")), (it) -> new Regex "[#{it.join ''}]"
            ]
          ]
        ]
      ]
    ]
  ]
  i LABEL:    C(St('&'), St('@'), R("WORD"))
  i WORD:     Re("[a-zA-Z\\._][a-zA-Z\\._0-9]*")
  i INT:      Re("[0-9]+"), (it) -> new Number it
  i _PIPE:    S(R("_"), St('|'))
  i _:        P(C(St(' '), St('\n')))
  i __:       P(C(St(' '), St('\n')), null, 1)
  i '.':      Re("[\\s\\S]")
  i ESC1:     S(St('\\'), R("."))
  i ESC2:     S(St('\\'), R(".")), (chr) -> '\\'+chr
]

@GRAMMAR = GRAMMAR = Grammar HandcompiledRules

@NODES = {
  GNode, Yes, Choice, Rank, Sequence, Lookahead, Existential, Pattern, Not, Ref, Regex, Grammar
}
