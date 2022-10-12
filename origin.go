package main

type Origin struct {
	code  string
	start SomePos
	end   SomePos
}

// hum... see also codestream.go, it has type Cursor {line int, pos int}
type SomePos struct {
	line int
	col  int
	pos  int
}
