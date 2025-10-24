package util

func Convert[E, R any](data []E, f func(E) R) []R {

	mapped := make([]R, len(data))

	for i, e := range data {
		mapped[i] = f(e)
	}

	return mapped
}

func Fold[E, R any](data []E, initial R, f func(R, E) R) R {

	res := initial

	for _, e := range data {
		res = f(res, e)
	}

	return res
}
