package censusmanagement

func (o OptNilString) ValueStringPointer() (v *string) {
	return getValuePointer(o)
}
