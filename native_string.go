package joeson

// NativeString wraps `string` and implements Ast.
type NativeString struct {
	Attr
	Str string
}

func NewNativeString(s string) NativeString { return NativeString{newAttr(), s} }
func (ns NativeString) String() string      { return ns.Str }
func (ns NativeString) assertNode()         {}
