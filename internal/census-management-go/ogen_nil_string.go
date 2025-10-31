package censusmanagement

func (o NilString) ValueStringPointer() *string {
	return getValuePointer(o)
}

func (o OptNilString) ValueStringPointer() *string {
	return getValuePointer(o)
}
