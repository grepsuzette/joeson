assert = require 'assert'
joe = require '../joescript_grammar'
jsx = require './javascript'

test = (code, expected) ->
  node = joe.GRAMMAR.parse code
  proc = []
  translated = jsx.translate(node)
  assert.equal translated.replace(/[\n ]+/g, ' '), expected.replace(/[\n ]+/g, ' ')

test """if true then 1 + 1 else 2 + 2""", 'if(true) { (1 + 1); } else { (2 + 2); };'
test """
if true
  1 + 1
  b = 2
  2
""", 'var b; if(true) { (1 + 1); b = 2; 2; };'
test """1 + 1""", '(1 + 1);'
test """
if 1
  2
else
  3""", 'if(1) { 2; } else { 3; };'
test """
foo =
  if 1
    2
  else
    3""", 'var foo; if(1) { foo = 2; } else { foo = 3; };'
test """(a) -> return a""", 'function(a) { return a; };'
test """(b) -> b""", 'function(b) { return b; };'
