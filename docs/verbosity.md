# Enable and control verbosity

## Per-parse operation verbosity

This will trace parsing of "153-4318-591" from beginning to end:
```go
gm := j.GrammarFromLines("telephone", /* rules */)
gm.ParseString("153-4318-591", j.Debug{true})`
```
## Per-rule verbosity

This will trace all parsing happening only for rule "phonenumber", but for *all* ParseString or ParseTokens operations:
```go
gm := j.GrammarFromLines("phonebook", []j.Line{
    o(named("fullname", `NAME+_{1,}`)),
    o(named("phonenumber", `[0-9]+'-'{1,}`), j.Debug{true}),
    i(named("NAME", `[a-zA-Z]*`), func (it j.Ast) j.Ast { return j.NewNativeStringFrom(it) }),
})
```
## Global verbosity

Usually however, especially when writing a grammar, it's just convenient to enable all tracing.  Verbosity is controlled by `TraceOptions`. Various settings are available:
```go
type TraceOptions struct {
	Stack      bool // print detailed parsing steps
	Grammar    bool // print grammar information and all rules
	Loop       bool // print all rules (this is all too similar to Grammar, one of them has to go)
	FilterLine int  // show only the Nth parse iteration (when n != -1 and Stack is true)
	SkipSetup  bool // mute traces during the setup of the joeson grammar
}
```
The first way to enable some verbosity is to build a specific TraceOptions as shown before:
```go
gm := joeson.GrammarWithOptionsFromLines( "myGrammarTitle", joeson.GrammarOptions{TraceOptions: joeson.Verbose()}, rules )
```
The second way is to set the environment, using `$TRACE` which is read from the tests and controls the tracing. From the CLI: `TRACE=all,skipsetup go test . --run=TestSquareroot -v`

Here are the possible values:

| Name       | Behavior                                          |
| ---------- | ------------------------------------------------  |
| none       | disable everything                                |
| stack      | print detailed parsing steps                      |
| loop       | print all rules in the grammar                    |
| grammar    | print grammar information and all rules           |
| skipsetup  | mute traces during joeson grammar setup           |
| all        | enable everything                                 |


Here is the same command-line as before with a `TRACE=all` prefix, affecting "all" to the environment variable `$TRACE`.

```
$ TRACE=all go test examples/calculator/calculator_test.go --run Test_12 -v
```

It shows very detailed trace:

![calculator_test2](https://user-images.githubusercontent.com/350354/216583646-0009d83d-36c2-457b-8cc3-e2aa0012edac.png)

This time with `TRACE=loop,grammar`:

```
$ TRACE=loop,grammar go test examples/calculator/calculator_test.go --run Test_12 -v
```

![calculator_test3](https://user-images.githubusercontent.com/350354/216583710-3a9fe967-2264-4b6a-8786-46a0f7d3edfc.png)

To disable the usage of environment, pass some TraceOptions of your choice to GrammarFromLines().

Here is a table to help:

| Function                                                         | Stack | Grammar | Read env? |
| ---------------------------------------------------------------- | ----- | ------- | --------- |
| DefaultTraceOptions() TraceOptions                               | 0     | 1       | yes       |
| TraceOptionsFromEnvironmentOrUse(opts TraceOptions) TraceOptions | 0     | 0       | yes       |
| Mute() TraceOptions                                              | 0     | 0       | no        |
| Verbose() TraceOptions                                           | 1     | 1       | no        |

