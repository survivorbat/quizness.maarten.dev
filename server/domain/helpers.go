package domain

// containsWithKey is a utility function that is similar to Contains(key, list) but
// allows you to use a callback to get a specific key.
func containsWithKey[T any, K comparable](is T, in []T, keyFunc func(T) K) (int, bool) {
	isKey := keyFunc(is)

	for index, item := range in {
		if isKey == keyFunc(item) {
			return index, true
		}
	}

	return 0, false
}
