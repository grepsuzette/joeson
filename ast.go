package joeson

// Ast is the result type of a Parse() operation by a grammar.
// It can be nil if a parser failed,
// otherwise it can almost anything. Presently it only requires a ContentString() method.
// (perhaps people would find it preferable to use `any` instead of Ast)
//
// See Parser, which is an interface for special Ast, capable of parsing
// (they are the byproduct of the joeson grammar).
type Ast interface {
	ContentString() string // a text representation of this ast. whatever you want.
}

type ParseError struct {
	ctx         *ParseContext
	ErrorString string
}

func (pe ParseError) ContentString() string {
	return "ERROR " + pe.ErrorString + " " + pe.ctx.String()
}

func NewParseError(ctx *ParseContext, s string) ParseError {
	return ParseError{ctx, s}
}

func IsParseError(ast Ast) bool {
	switch ast.(type) {
	case ParseError:
		return true
	default:
	}
	return false
}

// prefix(x) + x.ContentString(x)
func String(ast Ast) string {
	if x, isParser := ast.(Parser); isParser {
		return prefix(x) + x.ContentString()
	} else {
		return ast.ContentString()
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
			for k, value := range vWithPropertiesOf {
				vToExtend.Set(k, value)
			}
		case Parser:
			for k, value := range vWithPropertiesOf {
				switch k {
				case "label":
					vToExtend.SetLabel(value.(NativeString).Str)
				default:
					panic("unhandled property " + k + " in func (Ast) Merge(). toExtend=" + toExtend.ContentString() + " \n withPropertiesOf=" + withPropertiesOf.ContentString())
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
