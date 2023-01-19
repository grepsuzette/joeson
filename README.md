This is a Golang Packrat parser with [Left Recursion](https://raw.githubusercontent.com/jaekwon/JoeScript/master/docs/ipplrs_douglass.pdf) support ported from [https://github.com/jaekwon/JoeScript](https://github.com/jaekwon/JoeScript)

This is a library, there is no binary. Only tests and examples.

This doc is still TODO.

The port is as literal as possible.

# Usage

Run all tests recursively:

```
$ go test ./...
```

Run named test:

```
$ go test . --run TestHandcompiled -v
```

List tests in a file:

```
$ grep Test joeson_test.go 
func TestHandcompiled(t *testing.T) {
func TestRaw(t *testing.T) {
func Test100Times(t *testing.T) {
func TestDebugLabel(t *testing.T) {
```

Run a certain test from examples/calculator:

```
$ go test examples/calculator/calculator_test.go --run Test_73_plus_4 -v
```

## Trace options

It's possible to choose trace categories to show right from the command-line, using the `TRACE` env variable.

The structure returned by `grammar.Options()` is used by the parser to govern what to show:

```golang
// trace options. These options can produce traces as the parsing goes.
type TraceOptions struct {
	Stack      bool // print detailed parsing steps
    Grammar    bool // print grammar information and all rules
	Loop       bool // print all rules
	FilterLine int  // -1 to disable, no effect if !Stack. Only show parsing stack for the nth line
}
```

Here is how to change them using the `TRACE` environment variable (lower-case `trace=` also works, note):
```
TRACE=loop,stack,line=4 go test . --run TestHandcompiled -v
```

Here are the categories:

| Name       | Behavior                                          |
| ---------- | ------------------------------------------------  |
| none       | disable everything                                |
| stack      | print detailed parsing steps                      |
| loop       | print all rules in the grammar                    |
| line='N'   | only the stack trace for the nth line of code (⁺) |
| grammar    | print grammar information and all rules           |
| skipsetup  | mute traces during joeson grammar setup           |
| all        | print all that makes sense                        |

* ⁺: doesn't work well, whereas in js or go: all lines of code with the line system are line 0. We should probably use parseContext.Counter instead (TODO)
