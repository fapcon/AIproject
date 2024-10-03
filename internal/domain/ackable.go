package domain

type Ackable[T any] struct {
	Value T
	ack   func()
}

func NewAckable[T any](value T, ack func()) Ackable[T] {
	return Ackable[T]{
		Value: value,
		ack:   ack,
	}
}

func (a Ackable[T]) Ack() {
	a.ack()
}
