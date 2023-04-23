package main

type Duo struct {
	a string
	b string
}

func duo(a, b string) Duo { return Duo{a, b} }
