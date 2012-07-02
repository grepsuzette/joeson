http = require 'http'
connect = require 'connect'
{debug, info, warn, error:fatal} = (nogg=require('nogg')).logger 'server'
assert = require 'assert'
sugar = require 'sugar'

# logging
nogg.configure
  'default': [
    {file: 'logs/app.log',    level: 'debug'},
    {file: 'stdout',          level: 'debug'}]
  #'foo': [
  #  {file: 'foo.log',    level: 'debug'},
  #  {forward: 'default'}]
  'access': [
    {file: 'logs/access.log', formatter: null}]

# uncaught exceptions
process.on 'uncaughtException', (err) ->
  warn """\n
^^^^^^^^^^^^^^^^^
http://debuggable.com/posts/node-js-dealing-with-uncaught-exceptions:4c933d54-1428-443c-928d-4e1ecbdd56cb
#{err.message}
#{err.stack}
vvvvvvvvvvvvvvvvv
"""

# server
c = connect()
  .use(connect.logger())
  #.use(connect.staticCache())
  .use('/s', connect.static(__dirname + '/static'))
  .use(connect.favicon())
  .use(connect.cookieParser('TODO determine just how secret this is'))
  .use(connect.session({ cookie: { maxAge: 1000*60*60*24*30 }}))
  .use(connect.query())
  .use(connect.bodyParser())
c.use (req, res) ->
  res.writeHead 200, {'Content-Type': 'text/html'}
  res.end """
<html>
<link rel='stylesheet' type='text/css' href='http://fonts.googleapis.com/css?family=Anonymous+Pro'/>
<link rel='stylesheet' type='text/css' href='/s/style.css'/>
<script src='/s/jquery-1.7.2.js'></script>
<script src='/s/boot.js'></script>
<body>
  hello
</body>
</html>
"""

# server app
app = http.createServer(c)
io = require('socket.io').listen app
app.listen 8080

# kernel
{JKernel, JTypes} = require 'joeson/src/interpreter'
KERNEL = new JKernel
info "initialized kernel runloop"

# make output object
makeOut = (socket, threadId) ->
  write = (html) ->
    assert.ok typeof html is 'string', "makeOut/write wants a string, but got #{typeof html}"
    info "makeOut/write:", threadId:threadId, html
    socket.emit 'output', html:html, threadId:threadId
  write.close = ->
    socket.emit 'output', command:'close', threadId:threadId
  return write

# connect.io <-> kernel
io.sockets.on 'connection', (socket) ->

  # login
  # Caller (client) will receive a user object __str__
  socket.on 'login', ({name}) ->
    info "user #{name} login"

  # start code
  socket.on 'start', ({code,threadId}) ->
    info "received code #{code}, thread id #{threadId}"
    output = makeOut socket, threadId
    KERNEL.run
      user:   user
      code:   code
      output: output
      input:  undefined # not implemented
      callback: ->
        switch @state
          when 'return'
            unless @last is JTypes.JUndefined
              output(@last.__repr__(@).__html__(@))
            #output.close()
          when 'error'
            if @error.stack.length
              @printStack @error.stack
              stackTrace = @error.stack.map((x)->'  at '+x).join('\n')
              output("#{@error.name ? 'UnknownError'}: #{@error.message ? ''}\n  Most recent call last:\n#{stackTrace}")
            else
              output("#{@error.name ? 'UnknownError'}: #{@error.message ? ''}")
            #output.close()
          else
            throw new Error "Unexpected state #{@state} during kernel callback"
        @cleanup()
