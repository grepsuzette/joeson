{clazz} = require 'cardamom'

@inspect = (x) -> require('util').inspect x, false, 100

@CodeStream = CodeStream = clazz 'CodeStream', ->
  init: (@text, @pos=0, @buffer=null) ->
    @maxSeen=0

  pos$:
    enum: true
    conf: true
    get: -> (@_pos)
    set: (newPos) ->
      if @_pos isnt newPos
        @buffer = null
        @_pos = newPos
      if newPos > @maxSeen
        @maxSeen = newPos
        #console.log "#{@maxSeen}/#{@text.length} #{@text[@maxSeen..@maxSeen+10]}"
      @_pos

  # Get until the string `end` is encountered.
  # Change @pos accordingly, including the `end`.
  getUntil: (end, ignoreEOF=yes) ->
    index = @text.indexOf end, @pos
    if index is -1
      if ignoreEOF
        index = @text.length
      else
        throw new EOFError
    else
      index += end.length
    return @text[@pos...(@pos=index)]

  # Returns a string of up to chars chars or words words,
  # words:
  #  - are delimited by space(s)
  #  - include preceding space(s)
  #  - exclude trailing space(s)
  peek: ({chars, words, lines}) ->
    if chars?
      return @buffer = @text[@pos...@pos+chars]
    else if words?
      origPos = @pos
      while words > 0
        words -= 1 unless @getUntil(' ') is ' '
      @pos -= 1 if @pos > origPos and @text[@pos] is ' '
      result = @text[origPos...@pos]
      @pos = origPos
      return @buffer = result
    else if lines?
      throw new Error 'Not implemented yet'
    else
      return @buffer = @text[@pos]

  # Match a string or a regex
  # If matching a regex, don't forget to
  # set the buffer by calling .peek.
  # Regex returns null if match failed,
  # otherwise returns match[0] which may be ''
  match: ({regex, string}) ->
    if string?
      peek = @peek chars: string.length
      return null if peek isnt string
      @pos += string.length
      string
    else if regex?
      if not @buffer?
        throw new Error 'Buffer was null during regex match. Forget to peek?'
      matched = @buffer.match regex
      return null if matched is null
      @pos += matched[0].length
      matched[0]