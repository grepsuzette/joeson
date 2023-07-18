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
type (
	Ast interface {
		Locator
		String() string // text representation of this ast.
	}
	Locator interface {
		GetLocation() Origin
		SetLocation(Origin)
	}
	Origin struct {
		Code  *CodeStream
		Start int
		End   int
	}
)

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
	_ Ast = &NativeInt{}
	_ Ast = &NativeMap{}
	_ Ast = &NativeString{}
	_ Ast = &NativeUndefined{}
	_ Ast = &ParseError{}
)

type Attributes struct{ Location Origin }

func (attr *Attributes) GetLocation() Origin  { return attr.Location }
func (attr *Attributes) SetLocation(o Origin) { attr.Location = o }

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
