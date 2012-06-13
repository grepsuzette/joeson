{clazz, colors:{red, blue, cyan, magenta, green, normal, black, white, yellow}} = require('cardamom')
{inspect} = require 'util'
assert    = require 'assert'
_         = require 'underscore'

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
    return true if _.any @children, (child)->child.willDeclare(name)
    return false
  ensureVariable: (name) ->
    name = ''+name unless name instanceof joe.Undetermined
    @variables.push name unless @isDeclared name
  declareVariable: (name, isParameter=no) ->
    name = ''+name unless name instanceof joe.Undetermined
    @variables.push name unless name in @variables
    @parameters.push name unless name in @parameters if isParameter
  nonparameterVariables$: get: ->
    _.difference @variables, @parameters

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

  joe.Try::extend
    installScope: (options={}) ->
      init @, options
      @catchBlock.installScope(create:yes, parent:this) if @catchVar? and @catchBlock?
      @catchBlock.scope.declareVariable(@catchVar) if @catchVar?
      @withChildren (child, parent, attr) ->
        child.installScope?(create:no, parent:parent) unless attr is 'catchBlock'
      return this

  joe.Func::extend
    installScope: (options={}) ->
      init @, options
      @block.installScope(create:yes, parent:this) if @block?
      @block.scope.declareVariable(name, yes) for name in @params?.targetNames||[]
      @withChildren (child, parent, attr) ->
        child.installScope?(create:no, parent:parent) unless attr is 'block'
      return this

  joe.Assign::extend
    installScope: (options={}) ->
      init @, options
      @scope.ensureVariable(@target) if isVariable @target
      @withChildren (child, parent) ->
        child.installScope?(create:no, parent:parent)
      return this

  joe.Undetermined::extend
    word$:
      get: ->
        return "[Undetermined #{@prefix}]" if not @scope?
        loop
          _word = @prefix + randid(12) # lol.
          if not @scope.isDeclared(_word) and not @scope.willDeclare(_word)
            @word=_word
            return _word