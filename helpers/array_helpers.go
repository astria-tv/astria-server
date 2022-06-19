package helpers

// ElementExists returns tru if the given element exists in the provided list.
func ElementExists[T comparable](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
