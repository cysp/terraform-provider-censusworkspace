package testing

type multilevelMap2[K1, K2 comparable, T any] struct {
	m map[K1]map[K2]T
}

func NewMultilevelMap2[K1, K2 comparable, T any]() multilevelMap2[K1, K2, T] {
	return multilevelMap2[K1, K2, T]{
		m: make(map[K1]map[K2]T),
	}
}

func (mlm multilevelMap2[K1, K2, T]) Get(a K1, b K2) (T, bool) {
	if innerMap, exists := mlm.m[a]; exists {
		if value, exists := innerMap[b]; exists {
			return value, true
		}
	}

	var zero T
	return zero, false
}

func (mlm multilevelMap2[K1, K2, T]) Set(a K1, b K2, value T) {
	if _, exists := mlm.m[a]; !exists {
		mlm.m[a] = make(map[K2]T)
	}

	mlm.m[a][b] = value
}
