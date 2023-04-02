package internal

func Contains[E comparable](s []E, v E) bool {
	return ContainsFunc(s, func(e E) bool {
		return e == v
	})
}

func ContainsFunc[E any](s []E, f func(E) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}

	return false
}
