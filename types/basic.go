package types

// These all functions can be used in the optional types in a protobuf file.

// String will return the pointer of the string
func String(s string) *string {
	return &s
}

// Int64 will return the pointer of the int64
func Int64(i int64) *int64 {
	return &i
}

// Int will return the pointer of the int
func Int(i int) *int {
	return &i
}

// Float64 will return the pointer of the float64
func Float64(f float64) *float64 {
	return &f
}

// Float32 will return the pointer of the float32
func Float32(f float32) *float32 {
	return &f
}
