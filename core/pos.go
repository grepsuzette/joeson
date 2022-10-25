package core

type Cursor struct {
	line int
	col  int
	pos  int
}

type Origin struct {
	code  string
	start Cursor
	end   Cursor
}
