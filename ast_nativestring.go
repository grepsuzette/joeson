package joeson

// NativeInt and NativeString denote terminal nodes
// and at the same time satisfy the `Ast` interface. Joeson.coffee used
// Number or string directly.
type NativeString struct {
	Str string
}

func NewNativeString(s string) NativeString {
	return NativeString{s}
}
func (ns NativeString) GetGNode() *GNode            { return nil }
func (ns NativeString) HandlesChildLabel() bool     { return false }
func (ns NativeString) Prepare()                    {}
func (ns NativeString) ContentString() string       { return ns.Str }
func (ns NativeString) Parse(ctx *ParseContext) Ast { panic("uncallable") }

func (n NativeString) ForEachChild(f func(Ast) Ast) Ast { return n }
