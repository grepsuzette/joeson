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

![calculator_test1](https://user-images.githubusercontent.com/350354/216583474-4fd47a26-54a1-400a-aba6-96af1b06188f.png)

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

![calculator_test2](https://user-images.githubusercontent.com/350354/216583646-0009d83d-36c2-457b-8cc3-e2aa0012edac.png)

This time with `TRACE=loop,grammar`:

```
$ TRACE=loop,grammar go test examples/calculator/calculator_test.go --run Test_12 -v
```

![calculator_test3](https://user-images.githubusercontent.com/350354/216583710-3a9fe967-2264-4b6a-8786-46a0f7d3edfc.png)


