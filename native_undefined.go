package joeson

// Some parsers (namedly Not, Pattern, Sequence, Existential) need
// a value different from `nil` (which represents parsing failure).
type NativeUndefined struct{ *gnodeimpl }

func NewNativeUndefined() NativeUndefined {
	return NativeUndefined{NewGNode()}
}
func (nu NativeUndefined) ContentString() string { return "<NativeUndefined>" }

func (nu NativeUndefined) Parse(ctx *ParseContext) Ast               { return nu }
func (nu NativeUndefined) gnode() *gnodeimpl                         { return nu.gnodeimpl }
func (nu NativeUndefined) Prepare()                                  {}
func (nu NativeUndefined) HandlesChildLabel() bool                   { return false }
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
