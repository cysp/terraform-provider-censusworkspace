package censusmanagement

type valueGetter[T any] interface {
	Get() (T, bool)
}

func getValuePointer[T any](o valueGetter[T]) *T {
	value, ok := o.Get()

	if !ok {
		return nil
	}

	return &value
}
