package helpers

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

// func (o *NullInt) UnmarshalJSON(data []byte) error {
//     if string(data) != "null" {
//         if err := json.Unmarshal(data, &o.Int); err != nil {
//             return err
//         }
//         o.IsValid = true
//     }
//     return nil
// }

// func (o NullInt) MarshalJSON() ([]byte, error) {
//     if o.IsValid {
//         return json.Marshal(o.Int)
//     }
//     return json.Marshal(nil)
// }
