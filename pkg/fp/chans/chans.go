package chans

func Map[E, V any](input <-chan E, convert func(E) V) <-chan V {
	out := make(chan V)
	go func() {
		for e := range input {
			out <- convert(e)
		}
		close(out)
	}()
	return out
}

func MapFilter[E, V any](input <-chan E, f func(E) (V, bool)) <-chan V {
	out := make(chan V)
	go func() {
		for e := range input {
			v, ok := f(e)
			if ok {
				out <- v
			}
		}
		close(out)
	}()
	return out
}

func ToSlice[E any](input <-chan E) []E {
	out := make([]E, 0, cap(input))
	for e := range input {
		out = append(out, e)
	}
	return out
}
