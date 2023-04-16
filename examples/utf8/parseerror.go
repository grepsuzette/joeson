package main

type ParseError struct {
	string
}

func NewParseError(s string) ParseError {
	return ParseError{s}
}

func (e ParseError) ContentString() string {
	return e.string
}
