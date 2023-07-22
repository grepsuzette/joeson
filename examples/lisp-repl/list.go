package main

type List struct {
	attr
	List []Expr
}

func list(a ...Expr) List { return List{attr{}, a} }
func nilList() List       { return List{attr{}, []Expr{}} }

func (o List) assertNode() {}
func (o List) Expr() Expr  { return Expr{attr{}, kindList, "", 0, o, ""} }
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
