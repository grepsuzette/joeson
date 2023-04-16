package joeson

// Abstract syntax tree, the result of a Parse() operation by a grammar.
// It can be almost anything you want. Presently it only requires a ContentString()
// method. Remark: perhaps people would find it preferable to use `any` instead of Ast.
//
// In particular, the byproduct of the joeson grammar are special Ast nodes
// that also satisfies Parser and are in turn capable of
// parsing things (see parser.go).
//
type Ast interface {
	ContentString() string // A text representation of this AST. Whatever you want.
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
