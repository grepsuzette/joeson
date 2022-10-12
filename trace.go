package main

type TraceSettings struct {
	stack      bool
	loop       bool
	skipSetup  bool
	filterLine int
}

// debugging variables and callbacks
// trace.filterLine<0 means disabled.
var timeStart func(name string) = nil
var timeEnd func(name string) = nil
var trace = TraceSettings{stack: false, loop: false, skipSetup: true, filterLine: -1}

type Origin struct {
	code  string
	start SomePos
	end   SomePos
}
type SomePos struct {
	line int
	col  int
	pos  int
}

// type that is maybe going to be transitory
type Result struct {
	origin *Origin
	m      map[string]*Result
}

func (r Result) toString() string {
	return "TODO result.toString"
}
