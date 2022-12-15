package helpers

// A value which is either set to an Int, or not set

type NullInt struct {
	Int   int
	IsSet bool
}

func NewNullInt(i int) NullInt {
	return NullInt{Int: i, IsSet: true}
}
func NewUndefinedNullInt() NullInt {
	return NullInt{Int: 0, IsSet: false}
}
func (ni *NullInt) Unset() {
	ni.IsSet = false
	ni.Int = -1
}
func (ni *NullInt) Set(n int) {
	ni.Int = n
	ni.IsSet = true
}

type NullBool struct {
	Bool  bool
	IsSet bool
}

func NewNullBool(b bool) NullBool {
	return NullBool{Bool: b, IsSet: true}
}
