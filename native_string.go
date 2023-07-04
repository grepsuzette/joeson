package joeson

// NativeString is a `string` wrapped in an object that satisfies the `Ast` interface.
type NativeString struct {
	Str string
}

func NewNativeString(s string) NativeString   { return NativeString{s} }
func (ns NativeString) ContentString() string { return ns.Str }
