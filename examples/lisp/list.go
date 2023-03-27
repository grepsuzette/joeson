package main

type List []Expr

func list(a ...Expr) List { return a }

func (o List) Expr() Expr { return Expr{kindList, "", 0, o, ""} }
func (o List) ContentString() string {
	s := "("
	first := true
	for _, v := range o {
		if !first {
			s += ", "
		}
		s += v.ContentString()
		first = false
	}
	s += ")"
	return s
}
