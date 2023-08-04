package joeson

type Attr struct {
	Code     string
	Start    int
	End      int
	Line     int
	RuleName string
}

func newAttr() *Attr {
	return &Attr{}
}

func isPredefinedAttr(attr string) bool {
	return false
}

func (attr *Attr) SetLine(n int) {
	attr.Line = n
	attr.Start = 0
	attr.End = 0
}

func (attr *Attr) GetLine() int {
	return attr.Line
}

func (attr *Attr) setRuleName(rulename string) {
	attr.RuleName = rulename
}

func (attr *Attr) SetOrigin(o Origin) {
	attr.Code = o.Code
	attr.Start = o.Start
	attr.End = o.End
	attr.Line = o.Line
	attr.RuleName = o.RuleName
}

func (attr *Attr) GetOrigin() Origin {
	return Origin{
		Code:     attr.Code,
		Line:     attr.Line,
		Start:    attr.Start,
		End:      attr.End,
		RuleName: attr.RuleName,
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
