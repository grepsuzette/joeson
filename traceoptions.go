package joeson

import (
	"os"
	"strconv"
	"strings"
)

// Trace options. They produce various traces during parsing.
type TraceOptions struct {
	Stack      bool // print detailed parsing steps
	Loop       bool // print all rules
	Grammar    bool // print grammar information and all rules
	FilterLine int  // to filter only the Nth line to parse when n != -1 and Stack is true
	SkipSetup  bool // mute traces during the setup of the joeson grammar
}

// Define what the default TraceOptions is when GrammarFromLines() is called optionless.
// By default if $TRACE envvar is specified, the options will be read from it
// starting from a Mute() state. If $TRACE is not specified, we only show succint
// trace to be beginner-friendly.
func DefaultTraceOptions() TraceOptions {
	return CheckEnvironmentForTraceOptionsOrUse(TraceOptions{
		Stack:      false, // why default is false: too verbose
		Loop:       false, // why default is false: prefer to use grammar
		Grammar:    true,
		SkipSetup:  true,
		FilterLine: -1,
	})
}

// Mute() creates a TraceOptions with all traces disabled.
func Mute() TraceOptions {
	return TraceOptions{Stack: false, Loop: false, Grammar: false, FilterLine: -1, SkipSetup: false}
}

// Verbose() creates a TraceOptions with Stack and Grammar tracing enabled.
func Verbose() TraceOptions {
	return TraceOptions{Stack: true, Loop: false, Grammar: true, FilterLine: -1, SkipSetup: true}
}

// If $TRACE (or $trace) environment variable is defined, derive the
// trace options from it, starting from a Mute() initial state. If
// the envvar is missing, defaultOpts is returned instead.
func CheckEnvironmentForTraceOptionsOrUse(defaultOpts TraceOptions) TraceOptions {
	return checkEnvironmentForTraceOptions(Mute())
}

// Given an initial set of trace options, this will attempt to read
// the environment variable $TRACE to extend it if possible.
// For instance `TRACE=loop,stack go test . --run TestHandcompiled -v`
// will show all rules in the grammar and the parsing steps.
//
// See the code for the option details, as the array does not render
// well in godoc.
func checkEnvironmentForTraceOptions(initial ...TraceOptions) TraceOptions {
	// Possible values (several are possible, comma-separated):
	//
	// | Name       | Behavior                                   |
	// | ---------- | ------------------------------------------ |
	// | none       | disable everything                         |
	// | stack      | print detailed parsing steps               |
	// | loop       | print all rules in the grammar             |
	// | line='N'   | only the stack trace for the nth line      |
	// | grammar    | print grammar information and all rules    |
	// | skipsetup  | mute traces during joeson grammar setup    |
	// | all        | print all that makes sense                 |
	var opt TraceOptions
	if len(initial) > 0 {
		opt = initial[0]
	} else {
		opt = DefaultTraceOptions()
	}
	env := os.Getenv("TRACE")
	if env == "" {
		env = os.Getenv("trace")
	}
	if env != "" {
		// as soon as the TRACE envvar is defined, we reset every option
		opt = Mute()
		for _, v := range strings.Split(env, ",") {
			switch v {
			case "none":
				opt = Mute()
			case "stack":
				opt.Stack = true
			case "loop":
				opt.Loop = true
			case "skipsetup":
				opt.SkipSetup = true
			case "grammar":
				opt.Grammar = true
			case "all":
				opt = Verbose()
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
					panic("unrecognized TRACE option: " + v + ". Recognized: TRACE=none,stack,loop,grammar,line=N,skipsetup")
				}
			}
		}
	}
	return opt
}
