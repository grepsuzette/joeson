package core

import (
	"os"
	"strconv"
	"strings"
)

// debugging variables and callbacks
// to rewrite soon in a more idiomatic way
var TimeStart func(name string) = nil
var TimeEnd func(name string) = nil

// trace options. They produce various traces during parsing.
type TraceOptions struct {
	Stack      bool // print detailed parsing steps
	Loop       bool // print all rules
	Grammar    bool // print grammar information and all rules
	FilterLine int  // to filter only the Nth line to parse when n != -1 and Stack is true
}

// The default TraceOptions. Use nvvar TRACE= with `go test`
func DefaultTraceOptions() TraceOptions {
	return TraceOptions{
		Stack:      true,
		Loop:       true,
		Grammar:    false,
		FilterLine: -1,
	}
}

// With this function, it's possible to extend `initial` with the envvar TRACE.
// For instance `TRACE=all go test . --run TestRaw -v` could be used to force
// all traces from the command-line without changing the code.
//
// Possible values (several are possible, comma-separated):
//
// | Name       | Behavior                                |
// | ---------- | --------------------------------------- |
// | none       | disable everything                      |
// | stack      | print detailed parsing steps            |
// | loop       | print all rules in the grammar          |
// | line='N'   | only the stack trace for the nth line   |
// | grammar    | print grammar information and all rules |
// | all        | print all that makes sense              |
//
// For instance `TRACE=loop,stack,line=4 go test . --run TestHandcompiled -v`
func CheckEnvironmentForTraceOptions(initial ...TraceOptions) TraceOptions {
	var opt TraceOptions
	if len(initial) > 0 {
		opt = initial[0]
	} else {
		// opt = TraceOptions{
		// 	Stack:      Trace.Stack,
		// 	Loop:       Trace.Loop,
		// 	Grammar:    Trace.Grammar,
		// 	FilterLine: Trace.FilterLine,
		// }
		opt = DefaultTraceOptions()
	}
	env := os.Getenv("TRACE")
	if env != "" {
		// as soon as the TRACE envvar is defined, we reset every option
		// and expect what to activate
		opt.Stack = false
		opt.Loop = false
		opt.Grammar = false
		opt.FilterLine = -1
		for _, v := range strings.Split(env, ",") {
			switch v {
			case "none":
				opt.Stack = false
				opt.Loop = false
				opt.Grammar = false
			case "stack":
				opt.Stack = true
			case "loop":
				opt.Loop = true
			case "grammar":
				opt.Grammar = true
			case "all":
				opt.Stack = true
				opt.Loop = true
				opt.Grammar = true
			default:
				if strings.Index(v, "line=") == 0 || strings.Index(v, "filterline=") == 0 {
					// line=4 or filterline=4 must result in opt.FilterLine = 4
					a := strings.Split(v, "=")
					if len(a) != 2 {
						panic("syntax: line=N, where N is a number")
					} else {
						if n, err := strconv.Atoi(a[1]); err == nil {
							opt.FilterLine = n
						} else {
							panic("syntax: line=N, where N is a number")
						}
					}
				} else {
					panic("unrecognized TRACE option: " + v + ". Recognized: TRACE=none,stack,loop,grammar,line=N")
				}
			}
		}
	}
	return opt
}
