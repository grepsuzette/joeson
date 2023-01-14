package helpers

// nilable int and bool

type NilableInt struct {
	Int   int
	IsSet bool
}

func NewNilableInt(i int) NilableInt {
	return NilableInt{Int: i, IsSet: true}
}

func (ni *NilableInt) Unset() {
	ni.IsSet = false
	ni.Int = -1
}
func (ni *NilableInt) Set(n int) {
	ni.Int = n
	ni.IsSet = true
}

type NilableBool struct {
	Bool  bool
	IsSet bool
}

func NewNilableBool(b bool) NilableBool {
	return NilableBool{Bool: b, IsSet: true}
}
