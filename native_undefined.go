package joeson

// Some parsers (namedly Not, Pattern, Sequence, Existential) need
// a value different from `nil` (which represents parsing failure).
type NativeUndefined struct {
	*Attributes
	*gnodeimpl
}

func NewNativeUndefined() NativeUndefined {
	return NativeUndefined{&Attributes{}, NewGNode()}
}
func (nu NativeUndefined) assertNode()    {}
func (nu NativeUndefined) String() string { return "<NativeUndefined>" }

func (nu NativeUndefined) Parse(ctx *ParseContext) Ast               { return nu }
func (nu NativeUndefined) gnode() *gnodeimpl                         { return nu.gnodeimpl }
func (nu NativeUndefined) prepare()                                  {}
func (nu NativeUndefined) handlesChildLabel() bool                   { return false }
func (nu NativeUndefined) ForEachChild(f func(Parser) Parser) Parser { return nu }

func notNilAndNotNativeUndefined(x Ast) bool {
	if x == nil {
		return false
	}
	if _, isUndefined := x.(NativeUndefined); isUndefined {
		return false
	}
	return true
}
