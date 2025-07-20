package client

// NewNilPointerString returns new NilString with value set to v.
func NewNilPointerString(v *string) NilString {
	if v == nil {
		return NilString{}
	}

	return NilString{
		Value: *v,
	}
}
