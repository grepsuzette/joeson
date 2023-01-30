package joeson

// Some nodes such as Not use undefined as a terminal node,
// whereas nil represents parsing failure.
type NativeUndefined struct{}

func NewNativeUndefined() NativeUndefined        { return NativeUndefined{} }
func (nu NativeUndefined) ContentString() string { return "<NativeUndefined>" }

// make NativeUndefined satisfy Parser as well
func (nu NativeUndefined) Parse(ctx *ParseContext) Ast               { return nu }
func (nu NativeUndefined) GetGNode() *GNode                          { return nil }
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
