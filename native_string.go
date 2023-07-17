package joeson

// NativeString wraps `string` and implements Ast.
type NativeString struct {
	Str string
}

func NewNativeString(s string) NativeString   { return NativeString{s} }
func (ns NativeString) ContentString() string { return ns.Str }
