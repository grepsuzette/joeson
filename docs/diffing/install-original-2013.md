# Install and update original JoeScript (coffeescript version)

Start by installing **coffeescript**; a recent version will do.

Clone the JoeScript repo:

```bash
git clone https://github.com/jaekwon/JoeScript
cd JoeScript
```

Then run this line (the defaults for cardamom and findit in package.json won't work):

`npm -i cardamom@^0.0.9 findit@^0.1.2`

It would normally work with`coffee -c src/joeson.coffee`.
However coffeescript made some breaking chance at some point, 
so you will get the following error: "unexpected indentation".

Luckily it's easy to fix:

1. Add a backslash (`\`) at the end of lines 90-92, like below:
```
      return if trace.filterLine? and line isnt trace.filterLine
      codeSgmnt = "#{ white ''+line+','+@code.col \ <- here
                }\t#{ black pad right:5, (p=escape(@code.peek beforeChars:5))[p.length-5...] \ <- here
                  }#{ green pad left:20, (p=escape(@code.peek afterChars:20))[0...20] \ <- and here
```

2. Add a backslash at the end of lines 207-210. Save. `coffee -c src/joeson.coffee` should now work.
3. Edit tests/joeson_test.coffee. Add a backslash at the end of lines 85 and 86.
4. Test it: `coffee -c src/joeson.coffee && coffee tests/joeson_test.coffee` should work.

Done.
