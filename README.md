Golang Packrat parser with [Left Recursion](https://raw.githubusercontent.com/jaekwon/JoeScript/master/docs/ipplrs_douglass.pdf) ported from [https://github.com/jaekwon/JoeScript](https://github.com/jaekwon/JoeScript), for the [Gnoland](https://github.com/gnolang/gno) project.

A good doc is still TBD.
Here are some bits:

- [docs/testing](docs/testing.md)
- [docs/faq](docs/faq.md)
- [docs/diffing](docs/diffing.md)
- [docs/internals](docs/internals.md)

# Testing

This is a library, there is no `main` function.

> Run all tests recursively:

```
go test ./...
```

> Run named test:
```
go test . --run TestHandcompiled -v
```

> Run a certain test from examples/calculator:
```
go test examples/calculator/calculator_test.go --run test_12 -v
```

![calculator_test1](https://user-images.githubusercontent.com/350354/216583474-4fd47a26-54a1-400a-aba6-96af1b06188f.png)

# Enable and control verbosity

Verbosity is controlled by `TraceOptions`:

```go
type TraceOptions struct {
	Stack      bool // print detailed parsing steps
	Loop       bool // print all rules
	Grammar    bool // print grammar information and all rules
	FilterLine int  // to filter only the Nth line to parse when n != -1 and Stack is true
	SkipSetup  bool // mute traces during the setup of the joeson grammar
}
```

The first way to enable some verbosity is to build a specific TraceOptions as shown before:

```go
gm := joeson.GrammarFromLines(rules, "myGrammarTitle", joeson.GrammarOptions{TraceOptions: joeson.Verbose()})
```

The second way is to set the environment, using `$TRACE` which is read from the tests and controls the tracing. So from the CLI, here is one way to enable traces: `TRACE=all,skipsetup go test . --run=TestSquareroot -v`

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

# Error handling

A Parser doesn't panic but returns `Ast` or `nil`.

In custom parsing callbacks, 

* return `nil` whenever parsing failed and the parser declines the responsability for what is being parsed. 
  * This will pass on the hand to the next rule in a rank.
* return `ParseError` whenever there is a problem in what is being parsed (`ParseError` implements `Ast`).
* return a regular `Ast` otherwise.

## Example: handle error to parse an octal byte value that must NOT exceed 255

This example comes from Go. In Go you can specify an octal triplet using `\ooo` notation, 
where `o` is an octal digit. For instance to print a star (0x2a, 42 or 052 in
octal) you can use `fmt.Println("\052")`. 

The problem with that format is the `[0-7]{3,3}` triplet can easily outflow
a byte, which a parser would want to detect:
```go
i(named("octal_digit", "[0-7]")),
i(named("octal_digits", "octal_digit ('_'? octal_digit)+")),
i(named("octal_byte_value", "'\\\\' octal_digit{3,3}"), func(ast joeson.Ast) joeson.Ast {
    // check <= 255
    n := joeson.NewNativeIntFrom(ast).Int()
    if n > 255 {
        return NewParseError(ctx, "illegal: octal value over 255")
    } else {
        return n
    }
}),
```

This way there is a single return and no panic.
