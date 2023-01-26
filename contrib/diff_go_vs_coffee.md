diff_go_vs_coffee is a bash script using the tool vimdiff (part of the vim editor distribution)
to show the trace differences between the go and coffeescript implementation of joeson.

# Before using this tool

* Make sure to edit the "joePath=~/Work/GNO/JoeScript" in the file so it points to whereever you installed your JoeScript.

* Install https://www.github.com/powerman/vim-plugin-AnsiEsc :: It will give a `:AnsiEsc` command that shows ansi colors in vim. Without this, you're better off using some `diff` variant because of escape sequences.

* You may want this line in your vimfile `autocmd BufRead *.ansi AnsiEsc` to auto-enable AnsiEsc on files with an `.ansi` extension

## Default options

TODO Improve this.

Some work is still required here, but in the meantime it's worth a note.

The default grep pattern is set to `grepPattern="^${BOLD}[[:digit:]]"` where `${BOLD}` is "\e[1m", which means it will filters only the Stack traces. 

To show and diff everything, use `-g ""`.

The default function pattern is `goTestFunc="^TestRaw"`.

## Useful keys in vimdiff

* `<Esc>:AnsiEsc<CR>` :: toggle ansi colors
* `h`, `j`, `k`, `l` :: move around (in normal mode)
* `<Ctrl-w> + w` :: jump to the opposite split window
* `zo`, `zc` :: open fold, close fold
* `<Esc>:q<CR>` :: quit

## Grep + diff

This is probably the most useful option is `--grep <REGEX>`.
Only the lines matching REGEX will be in both columns.

For instance, to only diff the respective lines of the 2 programs beginning by "Loop":

`./diff_go_vs_coffee -g "^Loop"`


