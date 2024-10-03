package fp

type Optional[V any] struct {
	Value   V
	Present bool
}

func NewOptional[V any](value V) Optional[V] {
	return Optional[V]{
		Value:   value,
		Present: true,
	}
}

func Identity[V any](value V) V {
	return value
}

func (o Optional[V]) Get() (V, bool) {
	return o.Value, o.Present
}

func (o Optional[V]) GetValue() V {
	return o.Value
}

func (o Optional[V]) IsPresent() bool {
	return o.Present
}
