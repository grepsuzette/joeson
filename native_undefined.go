package joeson

// Some parsers (namedly Not, Pattern, Sequence, Existential) need
// a value different from `nil` (which represents parsing failure).
type NativeUndefined struct {
	Attr
	*gnodeimpl
}

func NewNativeUndefined() NativeUndefined {
	return NativeUndefined{newAttr(), newGNode()}
}
func (nu NativeUndefined) assertNode()    {}
func (nu NativeUndefined) String() string { return "<NativeUndefined>" }

func (nu NativeUndefined) Parse(ctx *ParseContext) Ast               { return nu }
func (nu NativeUndefined) gnode() *gnodeimpl                         { return nu.gnodeimpl }
func (nu NativeUndefined) prepare()                                  {}
func (nu NativeUndefined) handlesChildLabel() bool                   { return false }
func (nu NativeUndefined) ForEachChild(f func(Parser) Parser) Parser { return nu }

func isNotUndefined(x Ast) bool {
	if x == nil {
		return false
	}
	_, isUndefined := x.(NativeUndefined)
	return !isUndefined
}
