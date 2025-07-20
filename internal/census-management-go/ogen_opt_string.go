package censusmanagement

// NewOptPointerString returns new OptString with value set to v.
func NewOptPointerString(v *string) OptString {
	if v == nil {
		return OptString{}
	}

	return NewOptString(*v)
}

func (o OptString) ValueStringPointer() (v *string) {
	return getValuePointer(o)
}
