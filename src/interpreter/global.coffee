{clazz, colors:{red, blue, cyan, magenta, green, normal, black, white, yellow}} = require('cardamom')
{inspect} = require 'util'
assert = require 'assert'
_ = require 'underscore'
joe = require('joeson/src/joescript').NODES
{pad, escape, starts, ends} = require 'joeson/lib/helpers'
{extend, isVariable} = require('joeson/src/joescript').HELPERS
{debug, info, warn, error:fatal} = require('nogg').logger 'server'

{GOD} = require 'joeson/src/interpreter'
{JObject, JArray, JUser, JUndefined, JNull, JNaN} = require 'joeson/src/interpreter/object'

@print = ($, [obj]) ->
  $.output(obj.__html__($) + '<br/>')
  return JUndefined

module.exports = new JObject creator:GOD, data:@