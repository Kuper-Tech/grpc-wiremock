package sliceutils

func FirstOf[T any](slice []T) T {
	var result T

	if len(slice) > 0 {
		result = slice[0]
	}

	return result
}

func SliceToMap[T comparable](slice []T) map[T]struct{} {
	mapAsSet := map[T]struct{}{}

	for _, el := range slice {
		mapAsSet[el] = struct{}{}
	}

	return mapAsSet
}
