package joeson

// NativeString wraps `string` and implements Ast.
type NativeString struct {
	*Origin
	Str string
}

func NewNativeString(s string) NativeString { return NativeString{&Origin{}, s} }
func (ns NativeString) String() string      { return ns.Str }
func (ns NativeString) assertNode()         {}
