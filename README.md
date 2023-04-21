Golang Packrat parser with [Left Recursion](https://raw.githubusercontent.com/jaekwon/JoeScript/master/docs/ipplrs_douglass.pdf) ported from [https://github.com/jaekwon/JoeScript](https://github.com/jaekwon/JoeScript), for the [Gnoland](https://github.com/gnolang/gno) project.

A good doc is still TBD.
Here are some bits:

- [docs/testing](docs/testing.md)
- [docs/faq](docs/faq.md)
- [docs/diffing](docs/diffing.md)
- [docs/internals](docs/internals.md)

## Error handling

A Parser doesn't panic but returns `Ast` or `nil`.

In custom parsing callbacks, 

* return `nil` whenever parsing failed and the parser declines the responsability for what is being parsed. 
  * This will pass on the hand to the next rule in a rank.
* return `ParseError` whenever there is a problem in what is being parsed (`ParseError` satisfies `Ast`).
* return a regular `Ast` otherwise.

### An example: parse an octal byte value that must NOT exceed 255 (ParseError otherwise)

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
