package utils

func InArray(list []string, file string) bool {
	for _, l := range list {
		if l == file {
			return true
		}
	}

	return false
}
