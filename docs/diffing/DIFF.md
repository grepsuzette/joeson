# Diffing

a) In 2013, https://github.com/jaekwon/JoeScript was written (coffeescript). See [how to install](install-original-2013.md)
b) This library is a very literal port of its joeson.coffee and joeson_test.coffee.

As of January 2023, their outputs are identical and should remain so until the go version becomes trusted enough (at which point divergence won't matter).

Diffing has been a way to bring convergence and remains a tool of choice if there are problems.

# Installing and choosing a JoeScript version

There exists 2 version of JoeScript's joeson.coffee. 

1. the [2013 version](https://github.com/jaekwon/JoeScript) ([how to install it](install-original-2013.md)),
2. a [fork enabling identical outputs](https://github.com/grepsuzette/JoeScript).

For diffing (very low-level debugging) purposes, the second one is better. It is also easier to install:

```bash
git clone https://github.com/grepsuzette/JoeScript
cd JoeScript
npm -i
 # transpile from coffeescript to javascript
coffee -c src/joeson.coffee && coffee -c tests_joeson.coffee
```

## vimdiff

A bash script to use vimdiff is supplied.
It's put in contrib/diff_go_vs_coffee (see also its [doc](../../contrib/diff_go_vs_coffee.md)).

