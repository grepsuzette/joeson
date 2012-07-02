{clazz, colors:{red, blue, cyan, magenta, green, normal, black, white, yellow}} = require('cardamom')
{inspect} = require 'util'
assert    = require 'assert'

joe = require('joeson/src/joescript').NODES
{extend, isWord, isVariable} = require('joeson/src/joescript').HELPERS
{randid}  = require 'joeson/lib/helpers'

# A heirarchical lexical scope structure.
@LScope = LScope = clazz 'LScope', ->
  init: (parent) ->
    # this is to make scopes and nodes non-circular.
    Object.defineProperty this, 'parent', value:parent, enumerable:no
    @variables  = [] #
    @parameters = []
    @children   = [] # child LScopes
    @parent.children.push this if @parent?
  declares: (name) ->
    name = ''+name unless name instanceof joe.Undetermined
    return name in @variables
  isDeclared: (name) ->
    name = ''+name unless name instanceof joe.Undetermined
    return true if name in @variables
    return true if @parent?.isDeclared(name)
    return false
  willDeclare: (name) ->
    name = ''+name unless name instanceof joe.Undetermined
    return true if name in @variables
    return true if @children.some (child)->child.willDeclare(name)
    return false
  ensureVariable: (name) ->
    name = ''+name unless name instanceof joe.Undetermined
    @variables.push name unless @isDeclared name
  declareVariable: (name, isParameter=no) ->
    name = ''+name unless name instanceof joe.Undetermined
    @variables.push name unless name in @variables
    @parameters.push name unless name in @parameters if isParameter
  nonparameterVariables$: get: ->
    @variables.subtract @parameters

# Node::installScope: (plugin)
# Installs lexical scopes on nodes and collect variables and parameters.
#
# CONTRACT:
#   Calling installScope on a node with scope already installed
#   should be a safe operation that re-installs the scope.
#   After node transformations like node.toJSNode(), you need to re-install.
@install = ->
  return if joe.Node::installScope? # already defined.

  init = (node, options) ->
    # Dependency validation
    if options.create or not options.parent?
      node.scope = node.ownScope = new LScope options.parent?.scope
    else
      node.scope = options.parent.scope

  joe.Node::extend
    installScope: (options={}) ->
      init @, options
      @withChildren (child, parent) ->
        child.installScope?(create:no, parent:parent)
      return this
    determine: ->
      that = this
      @withChildren (child, parent, key, desc, index) ->
        if child instanceof joe.Undetermined
          child.determine()
          if index?
            that[key][index] = child.word
          else
            that[key] = child.word
        else if child instanceof joe.Node
          child.determine()
      @

  joe.Try::extend
    installScope: (options={}) ->
      init @, options
      @catchBlock.installScope(create:yes, parent:this) if @catchVar? and @catchBlock?
      @catchBlock.scope.declareVariable(@catchVar) if @catchVar?
      @withChildren (child, parent, key) ->
        child.installScope?(create:no, parent:parent) unless key is 'catchBlock'
      return this

  joe.Func::extend
    installScope: (options={}) ->
      init @, options
      @block.installScope(create:yes, parent:this) if @block?
      @block.scope.declareVariable(name, yes) for name in @params?.targetNames||[]
      @withChildren (child, parent, key) ->
        child.installScope?(create:no, parent:parent) unless key is 'block'
      return this

  joe.Assign::extend
    installScope: (options={}) ->
      init @, options
      @scope.ensureVariable(@target) if isVariable @target
      @withChildren (child, parent) ->
        child.installScope?(create:no, parent:parent)
      return this

  joe.Undetermined::extend
    determine: ->
      return if @word? # already determined.
      assert.ok @scope?, "Scope must be available to determine an Undetermined"
      loop
        word = @prefix+'_'+randid(4)
        if not @scope.isDeclared(word) and not @scope.willDeclare(word)
          return @word=joe.Word(word)
