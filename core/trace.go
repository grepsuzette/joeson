package core

// direct port from original coffeescript impl.
// probably should rewrite

// debugging variables and callbacks
// trace.filterLine<0 means disabled.
var TimeStart func(name string) = nil
var TimeEnd func(name string) = nil

var Trace = TraceSettings{
	Stack:      true,
	Loop:       true,
	SkipSetup:  false,
	FilterLine: -1,
}

type TraceSettings struct {
	Stack      bool
	Loop       bool
	SkipSetup  bool
	FilterLine int
}
