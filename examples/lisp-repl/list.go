package main

type List []Expr

func list(a ...Expr) List { return a }

func (o List) Expr() Expr { return Expr{kindList, "", 0, o, ""} }
func (o List) String() string {
	s := "("
	first := true
	for _, v := range o {
		if !first {
			s += ", "
		}
		s += v.String()
		first = false
	}
	s += ")"
	return s
}
