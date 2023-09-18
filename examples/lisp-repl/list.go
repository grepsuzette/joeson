package main

import j "github.com/grepsuzette/joeson"

type List struct {
	*j.Attr
	List []Expr
}

func list(a ...Expr) List { return List{j.NewAttr(), a} }
func nilList() List       { return List{j.NewAttr(), []Expr{}} }

func (o List) assertNode() {}
func (o List) Expr() Expr  { return Expr{j.NewAttr(), kindList, "", 0, o, ""} }
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
