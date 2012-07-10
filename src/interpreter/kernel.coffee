trace = debug:no, logCode:no

{clazz, colors:{red, blue, cyan, magenta, green, normal, black, white, yellow}} = require('cardamom')
{inspect} = require 'util'
assert = require 'assert'
{debug, info, warn, fatal} = require('nogg').logger __filename.split('/').last()

{randid, pad, escape, starts, ends} = require 'joeson/lib/helpers'
{
  NODES:joe
  HELPERS: {extend, isVariable}
} = require('joeson/src/joescript')

{@NODES, @HELPERS} = {NODES:{JObject, JArray, JUser, JUndefined, JNull, JNaN, JBoundFunc}} = require 'joeson/src/interpreter/object'

# installs instructions to joescript prototypes
require 'joeson/src/interpreter/instructions'

JStackItem = @JStackItem = clazz 'JStackItem', ->
  init: ({@node}) ->
    # figure out which function this node is declared in
    # used for printing a stack trace.
    # TODO make it lazy
    declaringFunc = @node.parent
    declaringFunc = declaringFunc.parent while declaringFunc? and declaringFunc not instanceof joe.Func
    @declaringFunc = declaringFunc
  toString: -> "'#{@node?.toJavascript?()}' (source:#{@declaringFunc}, line:#{@node._origin?.line}, col:#{@node._origin?.col})"

# A runtime context. (Represents a thread/process of execution)
# user:     Owner of the process
# scope:    All the local variables, a dual of the lexical scope.
# i9ns:     Instructions, a "stack" that also stores intermediate data.
# error:    Last thrown error
JThread = @JThread = clazz 'JThread', ->

  # kernel: JKernel to which this thread belongs
  # start:  The start node of program to run
  # user:   The user associated with this thread
  # scope:  Immediate local lexical scope object
  init: ({@kernel, @start, @user, @scope, @callback}) ->
    assert.ok @kernel instanceof JKernel,  "JThread wants kernel"
    assert.ok @start  instanceof joe.Node, "JThread wants function"
    assert.ok @user   instanceof JObject,  "JThread wants user"
    @scope ?= new JObject creator:@user
    assert.ok @scope  instanceof JObject,  "JThread scope not JObject"
    @id = randid()
    @i9ns = [] # i9n stack
    @last = JUndefined # last return value.
    @state = null
    @push this:@start, func:@start.interpret
    # if @user is GOD then @will = -> yes # optimization

  # Main run loop iteration.
  # return:
  #   'error'     for uncaught errors. see @error
  #   'return'    for the final return value. see @last
  #   'wait'      to wait for IO.
  #   null        all other intermediate cases.
  runStep: ->
    if @i9ns.length is 0
      return @state='return'
    {func, this:that, target, targetKey, targetIndex} = i9n = @i9ns[@i9ns.length-1]
    info blue "             -- runStep --" if trace.debug
    @printScope @scope if trace.debug
    @printStack() if trace.debug
    throw new Error "Last i9n.func undefined!" if not func?
    throw new Error "target and targetKey must be present together" if (target? or targetKey?) and not (target? and targetKey?)
    #key = "#{that.constructor.name}.#{func._name}"
    #timeit key, =>
    #  @last = func.call that, this, i9n, @last
    @last = func.call that ? i9n, this, i9n, @last
    switch @state
      when null
        info "             #{blue 'last ->'} #{@last}" if trace.debug
        if targetIndex?
          target[targetKey][targetIndex] = @last
        else if target?
          target[targetKey] = @last
        return null
      when 'error'
        info "             #{red 'throw ->'} #{@last}" if trace.debug
        loop # unwind loop
          dontcare = @pop()
          i9n = @peek()
          if not i9n?
            return 'error'
          else if i9n.this instanceof joe.Try and not i9n.isHandlingError
            i9n.isHandlingError = true
            i9n.func = joe.Try::interpretCatch
            last = @error
            return @state=null
      when 'return'
        info "             #{yellow 'return ->'} #{@last}" if trace.debug
        loop # unwind loop
          dontcare = @pop()
          i9n = @peek()
          if not i9n?
            return 'return'
          else if i9n.this instanceof joe.Invocation
            assert.ok i9n.func is joe.Invocation::interpretFinal, "Unexpected i9n.func #{i9n.func?._name or i9n.func?._name}"
            return @state=null
      when 'wait'
        info "             #{yellow 'wait ->'} #{inspect @waitKey}" if trace.debug
        existing = @kernel.waitlist[waitKey]
        if existing?
          @kernel.waitlist[waitKey].push thread
        else
          @kernel.waitlist[waitKey] = [thread]
        return 'wait'
      else
        throw new Error "Unexpected state #{@state}"

  ### STACKS ###

  pop: -> @i9ns.pop()

  peek: -> @i9ns[@i9ns.length-1]

  push: (i9n) -> @i9ns.push i9n

  callStack: ->
    stack = []
    for item in @i9ns when item.this instanceof joe.Invocation
      stack.push JStackItem node:item.this
    return stack

  ### FLOW CONTROL ###

  throw: (name, message) ->
    @error = name:name, message:message, stack:@callStack()
    @state = 'error'
    return

  return: (result) ->
    assert.ok result?, "result value can't be undefined. Maybe JUndefined?"
    @state = 'return'
    return result # return the result of this to set @last.

  wait: (waitKey) ->
    @waitKey = waitKey
    @state = 'wait'
    return

  awaken: (waitKey) ->
    throw new Error "TODO"
    # set @last
    # reanimate from waitlist

  exit: ->
    if @callback?
      @callback()
    else
      @cleanup()

  cleanup: ->
    # pass

  ### ACCESS CONTROL ###

  # Look at the object's acl to determine
  # if the action is permitted.
  will: (action, obj) ->
    return yes # TODO
    #return yes if obj.creator is @user
    #acl = obj.acl ? obj
    #throw new Error 'TODO determine permissing using ACL'

  toString: -> "[JThread]"

  ### DEBUG ###

  printStack: (stack=@i9ns) ->
    assert.ok stack instanceof Array
    for i9n, i in stack
      i9nCopy = Object.clone i9n
      delete i9nCopy.this
      delete i9nCopy.func
      info        "#{ blue pad right:12, "#{i9n.this?.constructor.name}"
                 }.#{ yellow i9n.func?._name
             }($, {#{ white Object.keys(i9nCopy).join ','
            }}, _) #{ black escape i9n.this }"

  printErrorStack: ->
    stackTrace = @error.stack.map((x)->'  at '+x).join('\n') or '  -- no stack trace available --'
    warn("#{@error.name ? 'UnknownError'}: #{@error.message ? ''}\n  Most recent call last:\n#{stackTrace}")

  printScope: (scope, lvl=0) ->
    for key, value of scope.data when key isnt '__proto__'
      try
        valueStr = value.__str__(@)
      catch error
        valueStr = "<ERROR IN __STR__: #{error}>"
      info "#{black pad left:13, lvl}#{red key}#{ blue ':'} #{valueStr}"
    @printScope scope.data.__proto__, lvl+1 if scope.data.__proto__?

  # for convenience, jml is available on a thread.
  jml: (args...) ->
    attributes = undefined
    if args[0] instanceof Object and args[0] not instanceof JObject
      attributes = args.shift()
    if args.length is 1 and args[0] instanceof Array
      elements = args[0]
    else
      elements = args
    if attributes?
      elements[''+key] = value for key, value of attributes
    return new JArray creator:@user, data:elements

# Multi-user time-shared interpreter.
@JKernel = JKernel = clazz 'JKernel', ->

  init: ({@cache}={}) ->
    @threads = []
    @cache ?= {}       # TODO should be weak etc.
    @index = 0
    @ticker = 0

  # Start processing another thread
  # user:     The same user object as returned by login.
  # callback: Called with thread after it exits.
  run: ({user, code, scope, callback}) ->
    {WORLD, ANON} = require 'joeson/src/interpreter/global'
    user ?= ANON
    scope ?= WORLD.create user
    assert.ok user?, "User must be provided."
    assert.ok user instanceof JUser, "User not instanceof JUser, got #{user?.constructor.name}"
    assert.ok scope?, "Scope missing for user #{user.name}"
    try
      if typeof code is 'string'
        info "received code:\n#{code}" if trace.debug or trace.logCode
        node = require('joeson/src/joescript').parse code
        info "unparsed node:\n" + node.serialize() if trace.debug or trace.logCode
        node = node.toJSNode(toValue:yes).installScope().determine()
        info "parsed node:\n" + node.serialize() if trace.debug or trace.logCode
      else
        assert.ok code instanceof joe.Node
        node = code
      thread = new JThread
        kernel:@
        start:node
        user:user
        scope:scope
        callback:callback
      @threads.push thread
      if @threads.length is 1
        @index = 0 # might have been NaN
        @runloop()
    catch error
      if node?
        warn "Error in user code start:", error.stack, "\nfor node:\n", node.serialize()
      else
        warn "Error parsing code:", error.stack, "\nfor code text:\n", code

  runloop$: ->
    @ticker++
    thread = @threads[@index]
    debug "tick #{@ticker}. #{@threads.length} threads, try #{thread.id}" if trace.debug
    try
      # TODO this reduces nextTick overhead, which is more significant when server is running (vs just testing)
      # kinda like a linux "tick", values is adjustable.
      for i in [0..20]
        exitCode = thread.runStep()
        if exitCode?
          @threads[@index..@index] = [] # splice out
          @index = @index % @threads.length # oops, sometimes NaN
          process.nextTick @runloop if @threads.length > 0
          thread.exit()
          return
      @index = (@index + 1) % @threads.length
      process.nextTick @runloop
    catch error
      fatal "Error thrown in runStep. Stopping execution, setting error. stack:\n" + (error.stack ? error)
      thread.throw 'InternalError', "#{error.name}:#{error.message}"
      @threads[@index..@index] = [] # splice out
      @index = @index % @threads.length # oops, sometimes NaN
      process.nextTick @runloop if @threads.length > 0
      thread.exit()
      return