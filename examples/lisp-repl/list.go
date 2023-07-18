package main

import (
	j "github.com/grepsuzette/joeson"
)

type List struct {
	*j.Attributes
	List []Expr
}

func list(a ...Expr) List { return List{&j.Attributes{}, a} }
func nilList() List       { return List{&j.Attributes{}, []Expr{}} }

func (o List) assertNode() {}
func (o List) Expr() Expr  { return Expr{&j.Attributes{}, kindList, "", 0, o, ""} }
func (o List) Length() int { return len(o.List) }
func (o List) String() string {
	s := "("
	first := true
	for _, v := range o.List {
		if !first {
			s += ", "
		}
		s += v.String()
		first = false
	}
	s += ")"
	return s
}
