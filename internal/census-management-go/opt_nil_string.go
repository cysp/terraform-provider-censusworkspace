package client

func (o OptNilString) ValueStringPointer() (v *string) {
	value, ok := o.Get()

	if !ok {
		return nil
	}

	return &value
}
