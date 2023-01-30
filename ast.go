package joeson

type Ast interface {
	ContentString() string // colorful representation
}

// Special kind of AST that is able to parse
// A compiled grammar is an AST whose nodes satisfy Parser.
type Parser interface {
	Ast
	Parse(ctx *ParseContext) Ast // Parse, update context's position, a return of nil indicates a parse failure.
	GetGNode() *GNode            // Grammar node, it can never be nil (even in the case of NativeUndefined a dummy is used)
	Prepare()                    // Called after children prepared
	HandlesChildLabel() bool
	ForEachChild(f func(Parser) Parser) Parser // depth-first walk enabler
}

func IsRule(parser Parser) bool {
	return parser.GetGNode().Rule == parser
}

// Show "<name>: " if `x` is a rule, or "<label>:", or empty string
func prefix(parser Parser) string {
	if IsRule(parser) {
		return red(parser.GetGNode().Name + ": ")
	} else if parser.GetGNode().Label != "" {
		return cyan(parser.GetGNode().Label + ":")
	} else {
		return ""
	}
}

// This is Prefix(x) + x.ContentString(x)
func String(ast Ast) string {
	switch x := ast.(type) {
	case Parser:
		return prefix(x) + x.ContentString()
	default:
		return x.ContentString()
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
					vToExtend.GetGNode().Label = value.(NativeString).Str
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
