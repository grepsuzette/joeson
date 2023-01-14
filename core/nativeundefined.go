package core

/*
 Some nodes such as Not use undefined,
 as a terminal node (similar to
 strings and numbers in coffeescript impl.).
 nil would represent parsing failure, therefore this type is required.

 For instance consider the code of Not.parse() in coffeescript:
  parse: @$wrap ($) ->
    pos = $.code.pos
    res = @it.parse $
    $.code.pos = pos
    if res isnt null
      return null
    else
      return undefined
*/
type NativeUndefined struct{}

func NewNativeUndefined() NativeUndefined              { return NativeUndefined{} }
func (nu NativeUndefined) ContentString() string       { return "<NativeUndefined>" }
func (nu NativeUndefined) GetGNode() *GNode            { return nil }
func (nu NativeUndefined) Prepare()                    {}
func (nu NativeUndefined) HandlesChildLabel() bool     { return false }
func (nu NativeUndefined) Labels() []string            { return []string{} }
func (nu NativeUndefined) Captures() []Ast             { return []Ast{} }
func (nu NativeUndefined) Parse(ctx *ParseContext) Ast { panic("unparsable?") }

// no Native* object must walk through children: see node.coffee:78 `if ptr.child instanceof Node`
func (n NativeUndefined) ForEachChild(f func(Ast) Ast) Ast { return n } // undefined has no children, thus f is not called

func NotNilAndNotNativeUndefined(x Ast) bool {
	if x == nil {
		return false
	}
	if _, isUndefined := x.(NativeUndefined); isUndefined {
		return false
	}
	return true
}
