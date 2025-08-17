package censusmanagement

// NewOptNilPointerString returns new OptNilString with value set to v.
func NewOptNilPointerString(v *string) OptNilString {
	if v == nil {
		return OptNilString{
			Set:  true,
			Null: true,
		}
	}

	return OptNilString{
		Set:   true,
		Value: *v,
	}
}

// NewOptNilPointerPointerString returns new OptNilString with value set to v.
func NewOptNilPointerPointerString(v **string) OptNilString {
	if v == nil {
		return OptNilString{}
	}

	if *v == nil {
		return OptNilString{
			Set:  true,
			Null: true,
		}
	}

	return OptNilString{
		Set:   true,
		Value: **v,
	}
}

func (o OptNilString) ValueStringPointer() (v *string) {
	return getValuePointer(o)
}

func (o OptNilString) ValueStringPointerPointer() (v **string) {
	return getOptNilValuePointerPointer(o)
}
