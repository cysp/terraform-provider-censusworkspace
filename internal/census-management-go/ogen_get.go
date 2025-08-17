package censusmanagement

type valueGetter[T any] interface {
	Get() (T, bool)
}

func getValuePointer[T any](o valueGetter[T]) (v *T) {
	value, ok := o.Get()

	if !ok {
		return nil
	}

	return &value
}

type optNilValueGetter[T any] interface {
	IsSet() bool
	IsNull() bool
	Get() (T, bool)
}

func getOptNilValuePointerPointer[T any](o optNilValueGetter[T]) (v **T) {
	if !o.IsSet() {
		return nil
	}

	value, ok := o.Get()

	if !ok {
		var rv *T = nil
		return &rv
	}

	var rv = &value
	return &rv
}
