# Testing

This is a library, there is no *main()*.


```
# Run all tests recursively:
$ go test ./...

# Run named test:
$ go test . --run TestHandcompiled -v

# Run a certain test from examples/calculator:
$ go test examples/calculator/calculator_test.go --run test_12 -v
```

![](./docs/assets/calculator_test1.png)

## Control traces with the TRACE environment variable

The `TRACE` variable is read from the tests (and distrib/diff_go_vs_coffee and the modified coffeescript version). It modifies the `TraceOptions` fields accordingly. 

For example: `TRACE=all,skipsetup go test . --run=TestSquareroot -v`

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

![](./docs/assets/calculator_test2.png)

This time with `TRACE=loop,grammar`:

```
$ TRACE=loop,grammar go test examples/calculator/calculator_test.go --run Test_12 -v
```

![](./docs/assets/calculator_test3.png)


