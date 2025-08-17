package censusmanagement

func (o NilString) ValueStringPointer() (v *string) {
	return getValuePointer(o)
}
