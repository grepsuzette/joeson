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

# vimdiff

A bash script to use vimdiff is supplied (`./diff_go_vs_coffee --help`). 

It's a 100 lines script, so feel free to read the code. Copy and modify it if you want to use a different tool (such as a plain `diff` if you want).

Before using it:

* Make sure to edit the "joePath=~/Work/GNO/JoeScript" in the file so it points to whereever you installed your JoeScript.

* Install https://www.github.com/powerman/vim-plugin-AnsiEsc :: It will give a `:AnsiEsc` command that shows ansi colors in vim. Without this, you're better off using some `diff` variant because of escape sequences.

* You may want this line in your vimfile `autocmd BufRead *.ansi AnsiEsc` to auto-enable AnsiEsc on files with an `.ansi` extension

## Useful keys in vimdiff

* `<Esc>:AnsiEsc<CR>` :: toggle ansi colors
* `h`, `j`, `k`, `l` :: move around (in normal mode)
* `zo`, `zc` :: open fold, close fold
* `<Esc>:q<CR>` :: quit

## Grep + diff

This is probably the most useful option is `--grep <REGEX>`.
Only the lines matching REGEX will be in both columns.

For instance, to only diff the respective lines of the 2 programs beginning by "Loop":

`./diff_go_vs_coffee -g "^Loop"`

