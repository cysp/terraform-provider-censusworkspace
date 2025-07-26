package client

// NewOptPointerString returns new OptString with value set to v.
func NewOptPointerString(v *string) OptString {
	if v == nil {
		return OptString{}
	}

	return NewOptString(*v)
}
