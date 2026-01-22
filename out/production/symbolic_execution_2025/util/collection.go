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

func FirstOrNil[E any](data []*E) *E {
	if len(data) == 0 {
		return nil
	}

	return data[0]
}

func Last[E any](data []E) E {
	return data[len(data)-1]
}

func IndexOf[E any](data []*E, v *E) int {
	for i, e := range data {
		if v == e {
			return i
		}
	}

	return -1
}
