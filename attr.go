package joeson

// Attr helps to implement Ast.
//
// To implement Ast you may simply do those 2 things:
//
// - Embed an `*Attr` field,
// - Have a `String() string` method.
type Attr struct {
	Code  string
	Start int
	End   int
	Line  int
}

func NewAttr() *Attr {
	return &Attr{}
}

// deprecated, replace with NewAttr
func newAttr() *Attr {
	return &Attr{}
}

func (attr *Attr) SetLine(n int) {
	attr.Line = n
	attr.Start = 0
	attr.End = 0
}

func (attr *Attr) GetLine() int {
	return attr.Line
}

func (attr *Attr) SetOrigin(o Origin) {
	attr.Code = o.Code
	attr.Start = o.Start
	attr.End = o.End
	attr.Line = o.Line
}

func (attr *Attr) GetOrigin() Origin {
	return Origin{
		Code:  attr.Code,
		Line:  attr.Line,
		Start: attr.Start,
		End:   attr.End,
	}
}

func (attr *Attr) HasAttribute(key interface{}) bool {
	return false
}

func (attr *Attr) GetAttribute(key interface{}) interface{} {
	return false
}

func (attr *Attr) SetAttribute(key interface{}, value interface{}) {
	panic("unimplemented")
}
