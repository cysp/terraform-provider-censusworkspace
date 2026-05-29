package censusmanagement

func (o OptNilString) ValueStringPointer() *string {
	if value, ok := o.Get(); ok {
		return &value
	}

	return nil
}
