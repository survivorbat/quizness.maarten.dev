package domain

func containsWithKey[T any, K comparable](is T, in []T, keyFunc func(T) K) (int, bool) {
	isKey := keyFunc(is)

	for index, item := range in {
		if isKey == keyFunc(item) {
			return index, true
		}
	}

	return 0, false
}
