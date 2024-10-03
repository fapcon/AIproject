package slices

func Map[E, V any](slice []E, convert func(E) V) []V {
	return MapFilter(slice, func(e E) (V, bool) {
		return convert(e), true
	})
}

func MapFilter[E, V any](slice []E, convert func(E) (V, bool)) []V {
	res := make([]V, 0, len(slice))
	for _, e := range slice {
		v, ok := convert(e)
		if ok {
			res = append(res, v)
		}
	}
	return res
}

func ToMap1[E any, K comparable, V any](slice []E, convert func(E) (K, V)) map[K]V {
	res := make(map[K]V, len(slice))
	for _, e := range slice {
		k, v := convert(e)
		res[k] = v
	}
	return res
}

func Merge[S ~[]V, V any](s ...S) S {
	var l int
	for _, slice := range s {
		l += len(slice)
	}
	r := make(S, 0, l)
	for _, slice := range s {
		r = append(r, slice...)
	}
	return r
}
