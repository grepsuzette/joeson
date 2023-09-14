package joeson

// Some parsers (namedly Not, Pattern, Sequence, Existential) need
// a value different from `nil` (which represents parsing failure).
//
// NativeUndefined need both to
// - implement Parser
// - get instantiated fast
type NativeUndefined struct {
	*Attr
	*rule
}

func NewNativeUndefined() NativeUndefined {
	nu := NativeUndefined{newAttr(), newRule()}
	return nu
}
func (nu NativeUndefined) assertNode()    {}
func (nu NativeUndefined) String() string { return "<NativeUndefined>" }

func (nu NativeUndefined) parse(ctx *ParseContext) Ast               { return nu }
func (nu NativeUndefined) getRule() *rule                            { return nu.rule }
func (nu NativeUndefined) prepare()                                  {}
func (nu NativeUndefined) handlesChildLabel() bool                   { return false }
func (nu NativeUndefined) forEachChild(f func(Parser) Parser) Parser { return nu }

func IsUndefined(x Ast) bool {
	if x == nil {
		return true
	}
	_, isUndefined := x.(NativeUndefined)
	return isUndefined
}
