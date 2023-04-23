package joeson

// Some nodes such as Not use undefined as a terminal node,
// whereas nil represents parsing failure.
// note NativeUndefined satisfies not only Ast but also Parser
type NativeUndefined struct{ *gnodeimpl }

func NewNativeUndefined() NativeUndefined {
	return NativeUndefined{NewGNode()}
}
func (nu NativeUndefined) ContentString() string { return "<NativeUndefined>" }

func (nu NativeUndefined) Parse(ctx *ParseContext) Ast               { return nu }
func (nu NativeUndefined) getgnode() *gnodeimpl                      { return nu.gnodeimpl }
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
