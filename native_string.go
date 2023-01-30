package joeson

// NativeInt and NativeString denote terminal nodes
// and at the same time satisfy the `Ast` interface. Joeson.coffee used
// Number or string directly.
type NativeString struct {
	Str string
}

func NewNativeString(s string) NativeString   { return NativeString{s} }
func (ns NativeString) ContentString() string { return ns.Str }
