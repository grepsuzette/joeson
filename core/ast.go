package core

import . "grepsuzette/joeson/colors"

type Ast interface {
	// Parse() reads from ParseContext, updates context's position,
	// returns nil to indicate parse failure.
	Parse(ctx *ParseContext) Ast

	ContentString() string // colorful representation of an AST node
	GetGNode() *GNode      // nodes without a grammar node (Native*) return nil
	Prepare()              // called after children prepared
	HandlesChildLabel() bool
	ForEachChild(f func(Ast) Ast) Ast
}

func IsRule(x Ast) bool {
	return x != nil && x.GetGNode() != nil && x.GetGNode().Rule == x
}

// Show "<name>: " if `x` is a rule, or "<label>:", or empty string
func Prefix(x Ast) string {
	if IsRule(x) {
		return Red(x.GetGNode().Name + ": ")
	} else if x.GetGNode().Label != "" {
		return Cyan(x.GetGNode().Label + ":")
	} else {
		return ""
	}
}

// This is Prefix(x) + x.ContentString(x)
func String(x Ast) string {
	return Prefix(x) + x.ContentString()
}

// This version is specialized for Astnode, and used instead of Object.merge()
// Extend a source object with the properties of another object (shallow copy).
// Careful, even though we return an Astnode, it is in fact the modified
// `toExtend` object that is returned.
// In original coffee impl., it is called lib/helpers.js:extend()
// (Object.extend)
func Merge(toExtend Ast, withPropertiesOf Ast) Ast {
	// @extend = extend = (object, properties) ->
	//   for key, val of properties
	//     object[key] = val
	//   object
	if toExtend == nil || withPropertiesOf == nil {
		return toExtend
	} else if _, isUndefined := withPropertiesOf.(NativeUndefined); isUndefined {
		return toExtend
	} else if h, isMap := withPropertiesOf.(NativeMap); isMap {
		if hToExtend, isMap := toExtend.(NativeMap); isMap {
			for k, v := range h {
				hToExtend.Set(k, v)
			}
		} else if toExtend.GetGNode() != nil {
			for k, v := range h {
				switch k {
				case "label":
					toExtend.GetGNode().Label = v.(NativeString).Str
				default:
					panic("unhandled property " + k + " in func (Astnode) Merge(). toExtend=" + toExtend.ContentString() + " \n withPropertiesOf=" + withPropertiesOf.ContentString())
				}
			}
		} else {
			panic("assert")
		}
		return toExtend
	} else if toExtend.GetGNode() == nil {
		panic("dont know how until we have SetGNode. toExtend=" + toExtend.ContentString() + " \n withPropertiesOf=" + withPropertiesOf.ContentString())
	} else {
		panic("Unhandled case in func (Astnode) Merge()")
	}
}
