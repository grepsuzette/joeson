package joeson

// Ast is the result type of a Parse() operation by a grammar.
//
// The return can be:
//
//   - `nil` if a parser failed (in a way where backtracking will happen),
//
//   - `ParseError` if a parsing must definitely fail (ParseError is an Ast).
//
//   - anything else you want that implements Ast. For that, you will
//     need a rule with a callback.
//
//     For joeson grammar without callbacks as you can see in examples/calculator, the
//     parser_*.go generate Ast nodes that are of types NativeArray, NativeInt,
//     NativeMap, NativeString.
//
//     To generate more specific Ast types, you may take a look at examples/lisp.
//
// Note: Parsers such as sequence, choice, not, pattern are also Ast,
// they are produced when parsing a valid joeson grammar; and they in turn help
// parsing that grammar.
type Ast interface {
	Node
	String() string // text representation of this ast.
}

var (
	_ Ast = &Grammar{}
	_ Ast = &choice{}
	_ Ast = &existential{}
	_ Ast = &lookahead{}
	_ Ast = &not{}
	_ Ast = &pattern{}
	_ Ast = &rank{}
	_ Ast = &regex{}
	_ Ast = &sequence{}
	_ Ast = &str{}
	_ Ast = &NativeArray{}
	_ Ast = NewNativeInt(0)
	_ Ast = &NativeMap{}
	_ Ast = &NativeString{}
	_ Ast = &NativeUndefined{}
	_ Ast = &ParseError{}
)

// WIP: we are converging towards gnolang.Node
// TODO there are some dubious aspects indicated by comments. Line and Label
// seem too narrow for joeson, gnolang being a bit more specific
type (
	Name string
	Node interface {
		assertNode()
		String() string
		Copy() Node
		GetLine() int // line is probably insufficient for joeson
		SetLine(int)
		GetLabel() Name // different from GetRuleLabel()
		SetLabel(Name)
		HasAttribute(key any) bool
		GetAttribute(key any) any
		SetAttribute(key any, value any)
	}
)

type Origin struct {
	Code  *CodeStream
	Start int
	End   int
}

// Attributes (from gnolang)
// All nodes have attributes for general analysis purposes.
type Attributes struct {
	Line int
	// Origin Origin
	Label Name
	data  map[interface{}]interface{}
}

// can delete this i think :(
// func (attr *Attributes) GetOrigin() Origin { return attr.Origin }
// func (attr *Attributes) SetOrigin(code *CodeStream, start int, end int) {
// 	attr.Origin = Origin{Code: code, Start: start, End: end}
// }

func (attr *Attributes) assertNode() {}
func (attr *Attributes) Copy() Node  { panic("Copy() not yet implemented") }

//	func (attr *Attributes) GetLine() int {
//		return attr.Origin.Code.PosToLine(attr.Origin.Start)
//	}
//
// func (attr *Attributes) SetLine(line int)                         { panic("use SetOrigin instead") } // { attr.Line = line }
func (attr *Attributes) GetLine() int                             { return attr.Line }
func (attr *Attributes) SetLine(line int)                         { attr.Line = line }
func (attr *Attributes) GetLabel() Name                           { return attr.Label }
func (attr *Attributes) SetLabel(label Name)                      { attr.Label = label }
func (attr *Attributes) HasAttribute(key interface{}) bool        { _, ok := attr.data[key]; return ok }
func (attr *Attributes) GetAttribute(key interface{}) interface{} { return attr.data[key] }
func (attr *Attributes) SetAttribute(key interface{}, value interface{}) {
	if attr.data == nil {
		attr.data = make(map[interface{}]interface{})
	}
	attr.data[key] = value
}

// prefix(x) + x.String(x)
func String(ast Ast) string {
	if x, isParser := ast.(Parser); isParser {
		return prefix(x) + x.String()
	} else {
		return ast.String()
	}
}

// Port of lib/helpers.js:extend() in a less general way (Ast-specific)
// Extend a source object with the properties of another object (shallow copy).
// The modified `toExtend` object is returned.
func merge(toExtend Ast, withPropertiesOf Ast) Ast {
	// @extend = extend = (object, properties) ->
	//   for key, val of properties
	//     object[key] = val
	//   object
	if toExtend == nil || withPropertiesOf == nil {
		return toExtend
	}
	switch vWithPropertiesOf := withPropertiesOf.(type) {
	case NativeUndefined:
		return toExtend
	case NativeMap:
		switch vToExtend := toExtend.(type) {
		case NativeMap:
			for k, value := range vWithPropertiesOf.Map {
				vToExtend.Set(k, value)
			}
		case Parser:
			for k, value := range vWithPropertiesOf.Map {
				switch k {
				case "label":
					vToExtend.SetRuleLabel(value.(NativeString).Str)
				default:
					panic("unhandled property " + k + " in func (Ast) Merge(). toExtend=" + toExtend.String() + " \n withPropertiesOf=" + withPropertiesOf.String())
				}
			}
		default:
			panic("assert")
		}
		return toExtend
	case Parser:
		switch toExtend.(type) {
		case Parser:
			panic("Unhandled case in func (Ast) Merge()")
		default:
			panic("assert")
		}
	default:
		panic("assert")
	}
}
