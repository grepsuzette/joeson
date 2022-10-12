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
