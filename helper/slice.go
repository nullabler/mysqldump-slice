package helper

func SliceIsExist[T comparable](key T, list []T) bool {
	for _, item := range list {
		if key == item {
			return true
		}
	}

	return false
}
