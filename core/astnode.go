package core

import . "grepsuzette/joeson/colors"

type Astnode interface {
	// Parse() reads from ParseContext, updates context's position,
	// returns nil to indicate parse failure.
	Parse(ctx *ParseContext) Astnode

	ContentString() string // colorful representation of an AST node
	GetGNode() *GNode      // nodes without a grammar node (Native*) return nil
	Prepare()              // called after children prepared
	HandlesChildLabel() bool
	// Labels() []string
	// Captures() []Astnode
	ForEachChild(f func(Astnode) Astnode) Astnode
}

func IsRule(x Astnode) bool {
	return x != nil && x.GetGNode() != nil && x.GetGNode().Rule == x
}

func Prefix(x Astnode) string {
	if IsRule(x) {
		return Red(x.GetGNode().Name + ": ")
	} else if x.GetGNode().Label != "" {
		return Cyan(x.GetGNode().Label + ":")
	} else {
		return ""
	}
}

func ContentStringWithPrefix(x Astnode) string {
	if x == nil {
		return "nil"
	} else {
		return Prefix(x) + x.ContentString()
	}
}

func LabelOrName(n Astnode) string {
	if IsRule(n) {
		return Red(n.GetGNode().Name + ": ")
	} else if n.GetGNode().Label != "" {
		return Cyan(n.GetGNode().Label + ":")
	}
	return ""
}

// Extend a source object with the properties of another object (shallow copy).
// Careful, even though we return an Astnode, it is in fact the modified
// `toExtend` object that is returned.
// In original coffee impl., it is called lib/helpers.js:extend()
// (Object.extend)
// this version is specialized for Astnode, and used in place of Object.merge(),
// which does a copy (because we don't really need it)
func Merge(toExtend Astnode, withPropertiesOf Astnode) Astnode {
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
			for _, key := range h.Keys() {
				hToExtend.Set(key, h.Get(key))
			}
		} else if toExtend.GetGNode() != nil {
			for _, key := range h.Keys() {
				var v Astnode = h.Get(key)
				switch key {
				case "label":
					toExtend.GetGNode().Label = v.(NativeString).Str
				default:
					panic("unhandled property " + key + " in func (Astnode) Merge(). toExtend=" + toExtend.ContentString() + " \n withPropertiesOf=" + withPropertiesOf.ContentString())
				}
			}
		} else {
			panic("assert")
		}
		return toExtend
	} else if toExtend.GetGNode() == nil {
		panic("dont know how until we have SetGNode. toExtend=" + toExtend.ContentString() + " \n withPropertiesOf=" + withPropertiesOf.ContentString())
		// TestRaw makes a call with toExtend= NativeMap{join:<NativeUndefined>, value:_PIPE}
	} else {
		panic("Unhandled case in func (Astnode) Merge()")
	}
	// } else {
	// 	gn := withPropertiesOf.GetGNode()
	// 	toExtend.GetGNode().Capture = gn.Capture
	// 	toExtend.GetGNode().CbBuilder = gn.CbBuilder
	// 	toExtend.GetGNode().Debug = gn.Debug
	// 	toExtend.GetGNode().Grammar = gn.Grammar
	// 	toExtend.GetGNode().Id = gn.Id
	// 	toExtend.GetGNode().Index = gn.Index
	// 	toExtend.GetGNode().Label = gn.Label
	// 	toExtend.GetGNode().Name = gn.Name
	// 	toExtend.GetGNode().Parent = gn.Parent
	// 	toExtend.GetGNode().Rule = gn.Rule
	// 	toExtend.GetGNode().Rules = gn.Rules
	// 	toExtend.GetGNode().RulesK = gn.RulesK
	// 	toExtend.GetGNode().SkipCache = gn.SkipCache
	// 	toExtend.GetGNode().SkipLog = gn.SkipLog
	// 	toExtend.GetGNode().Labels_ = gn.Labels_
	// 	toExtend.GetGNode().Captures_ = gn.Captures_
	// 	return toExtend
	// }
}
