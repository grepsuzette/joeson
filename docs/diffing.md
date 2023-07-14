**Diffing** has been a way to bring convergence and remains a tool of choice if there are problems. (Reminder: this was started as a port of a [coffeescript program](https://github.com/jaekwon/JoeScript)).

A bash wrapper to use the interactive diff tool [vimdiff](https://www.freecodecamp.org/news/compare-two-files-in-linux-using-vim/) was put in **contrib/diff_go_vs_coffee**.

# How to install JoeScript

Before you diff anything, you will need to install a version of JoeScript of course.

There are two version of JoeScript's `joeson.coffee`. 

1. a [fork enabling identical outputs](https://github.com/grepsuzette/JoeScript), recommended for diffing,
2. the [original version](https://github.com/jaekwon/JoeScript).

## Installing the trace-enabled JoeScript version

The code is similar to the original one save for side effects and added tests:

- TRACE environment variable allows to choose what trace category to output,
- TEST environment variable allows to choose the test to run,
- additionnal named tests were added,
- type names in the traces were modified to be compatible with the go version.

Install it like so:

```bash
git clone https://github.com/grepsuzette/JoeScript
cd JoeScript
npm -i
coffee -c src/joeson.coffee && coffee -c tests_joeson.coffee
```

## Installing the original version

Install **coffeescript**.
Clone the JoeScript repo:

```bash
git clone https://github.com/jaekwon/JoeScript
cd JoeScript
npm -i cardamom@^0.0.9 findit@^0.1.2
```

Coffeescript made some minor breaking changes, 
so it gives a couple of "unexpected indentation" errors when running `coffee -c src/joeson.coffee`. 

To fix it:

1. Add a backslash (`\`) at the end of lines 90-92, as shown below:
```
      return if trace.filterLine? and line isnt trace.filterLine
      codeSgmnt = "#{ white ''+line+','+@code.col \ <- here
                }\t#{ black pad right:5, (p=escape(@code.peek beforeChars:5))[p.length-5...] \ <- here
                  }#{ green pad left:20, (p=escape(@code.peek afterChars:20))[0...20] \ <- and here
```

2. Same, add a backslash at the end of lines 207-210. Save. 
3. Run `coffee -c src/joeson.coffee`. Now it should transpile without error.
4. Edit tests/joeson_test.coffee. Add a backslash at the end of lines 85 and 86. Save.
5. Test it: `coffee -c src/joeson.coffee && coffee tests/joeson_test.coffee`, it should work.

# debugging methodology

1. Compare the traces of the coffeescript and this implementation (TRACE=stack,grammar environment, see the README.md).
2. If there is trully a divergence, go on below.
3. Note the Counter at which point there is a divergence, the Counter comes last for instance for the following trace `0,0	     123 + 456           ] | |  * UnaryExpr: /([0-9])/g*{1,} 2` it would be ctx.Counter == 2.
4. Using the debugger of your choice (e.g. Visual Code), add a Conditional breakpoint in packrat.go, inside the callback in loopify(), and enter a condition like this: `ctx.Counter >= 2 && ctx.GrammarName == "<the name of your grammar>"`. This helps to avoid the breakpoint to be called when the joeson grammar itself or your own grammar is being parsed, and will save you a lot of time.
5. Actually find the problem.


