package joeson

// GetAttribute, SetAttribute and HasAttribute
// will respond to predefined attributes. This is a bit
// strange. Our thinking is it facilitates low-coupling
// with other packages (at least in early stage).
// TODO We will remove this if it proves unnecessary.
type Attr struct {
	Code     string
	Start    int
	End      int
	Line     int
	RuleName string
	// gnode    *gnodeimpl
	// meta     map[interface{}]interface{}   // removed, see reason in parser.go
}

// var predefinedAttrs = [5]string{"Code", "Start", "End", "Line", "RuleName"}

func newAttr() Attr {
	return Attr{}
}

func isPredefinedAttr(attr string) bool {
	// for _, s := range predefinedAttrs {
	// 	if s == attr {
	// 		return true
	// 	}
	// }
	return false
}

func (attr Attr) setRuleName(rulename string) {
	attr.RuleName = rulename
}

func (attr Attr) SetOrigin(o Origin) {
	attr.Code = o.Code
	attr.Start = o.Start
	attr.End = o.End
	attr.Line = o.Line
	attr.RuleName = o.RuleName
}

func (attr Attr) GetOrigin() Origin {
	return Origin{
		Code:     attr.Code,
		Line:     attr.Line,
		Start:    attr.Start,
		End:      attr.End,
		RuleName: attr.RuleName,
	}
}

func (attr Attr) HasAttribute(key interface{}) bool {
	// if s, is := key.(string); is {
	// 	if isPredefinedAttr(s) {
	// 		return true
	// 	}
	// }
	// _, ok := attr.meta[key]
	// return ok
	return false
}

func (attr Attr) GetAttribute(key interface{}) interface{} {
	// if s, ok := key.(string); ok {
	// 	switch s {
	// 	case "RuleName":
	// 		return attr.RuleName
	// 	case "Line":
	// 		return attr.Line
	// 	case "Start":
	// 		return attr.Start
	// 	case "End":
	// 		return attr.End
	// 	case "Code":
	// 		return attr.Code
	// 	}
	// }
	// return attr.meta[key]
	return false
}

func (attr Attr) SetAttribute(key interface{}, value interface{}) {
	panic("unimplemented")
	// if attr.meta == nil {
	// 	attr.meta = make(map[interface{}]interface{})
	// }
	// if s, ok := key.(string); ok {
	// 	switch s {
	// 	case "RuleName":
	// 		attr.RuleName = value.(string)
	// 	case "Line":
	// 		attr.Line = value.(int)
	// 	case "Start":
	// 		attr.Start = value.(int)
	// 	case "End":
	// 		attr.End = value.(int)
	// 	case "Code":
	// 		attr.Code = value.(string)
	// 	}
	// } else {
	// 	attr.meta[key] = value
	// }
}
