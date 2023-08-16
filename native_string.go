package joeson

// NativeString is an alias for`string` but implements Ast.
type NativeString string

func NewNativeString(s string) NativeString { return NativeString(s) }
func (ns NativeString) String() string      { return string(ns) }
func (ns NativeString) assertNode()         {}

func (ns NativeString) SetLine(n int)                                   {}
func (ns NativeString) GetLine() int                                    { return 1 }
func (ns NativeString) SetOrigin(o Origin)                              {}
func (ns NativeString) GetOrigin() Origin                               { return Origin{} }
func (ns NativeString) HasAttribute(key interface{}) bool               { return false }
func (ns NativeString) GetAttribute(key interface{}) interface{}        { return nil }
func (ns NativeString) SetAttribute(key interface{}, value interface{}) { panic("N/A") }
