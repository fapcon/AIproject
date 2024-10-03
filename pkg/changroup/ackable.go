package changroup

import (
	"sync"
)

type Ackable[T any] struct {
	Value T
	Ack   func()
}

func NewAckable[T any](value T, ack func()) Ackable[T] {
	return Ackable[T]{
		Value: value,
		Ack:   ack,
	}
}

type AckableGroup[T any] struct {
	channels *list[*channel[Ackable[T]]]
}

func NewAckableGroup[T any]() *AckableGroup[T] {
	return &AckableGroup[T]{
		channels: newList[*channel[Ackable[T]]](),
	}
}

// ReleaseAll releases all acquired channels and closes it.
// It's safe to call ReleaseAll several times as well as in parallel with ReleaseFunc.
func (g *AckableGroup[T]) ReleaseAll() {
	var all []*channel[Ackable[T]]
	g.channels.ForEach(func(ch *channel[Ackable[T]]) {
		all = append(all, ch)
	})
	for _, ch := range all {
		ch.release() // there will be deadlock if call it inside ForEach.
	}
}

// Acquire creates new channel and adds it to group.
//
// ReleaseFunc is returned as the second value.
// It should be called to remove the channel from group and close it.
// It's safe to call ReleaseFunc several times as well as in parallel with ReleaseAll.
func (g *AckableGroup[T]) Acquire() (<-chan Ackable[T], ReleaseFunc) {
	ch := g.channels.Insert(&channel[Ackable[T]]{
		ch:      make(chan Ackable[T]),
		done:    make(chan struct{}),
		send:    sync.WaitGroup{},
		release: nil, // is filled below
	})

	once := sync.Once{}
	ch.elem.release = func() {
		once.Do(func() {
			ch.Delete()
			close(ch.elem.done)
			ch.elem.send.Wait()
			close(ch.elem.ch)
		})
	}

	return ch.elem.ch, ch.elem.release
}

// Send sends copy of value to all acquired channels.
// Each copy has its own [Ackable.Ack].
// Original [Ackable.Ack] is called after all copies are acked.
// Send doesn't wait for ack.
// It guarantees that all channels receive values in the same order.
// And that order is the same as Send calls.
// It may block if some channel is not read until its release.
func (g *AckableGroup[T]) Send(value Ackable[T]) {
	send := sync.WaitGroup{}
	ack := sync.WaitGroup{}
	g.channels.ForEach(func(ch *channel[Ackable[T]]) {
		ack.Add(1)
		once := sync.Once{}
		v := NewAckable(value.Value, func() { once.Do(ack.Done) })

		// select is an optimisation to not create goroutine if someone reads the channel (should cover 90% cases)
		select {
		case ch.ch <- v:
		default:
			send.Add(1)
			ch.send.Add(1)
			go func() {
				defer send.Done()
				defer ch.send.Done()
				select {
				case ch.ch <- v:
				case <-ch.done:
					ack.Done()
				}
			}()
		}
	})
	go func() {
		ack.Wait()
		value.Ack()
	}()
	send.Wait()
}

// SendAsync sends value to all acquired channels, but unlike [AckableGroup.Send] doesn't block.
// Also, it doesn't preserve the order of values!
func (g *AckableGroup[T]) SendAsync(value Ackable[T]) {
	ack := sync.WaitGroup{}
	g.channels.ForEach(func(ch *channel[Ackable[T]]) {
		ack.Add(1)
		once := sync.Once{}
		v := NewAckable(value.Value, func() { once.Do(ack.Done) })

		// select is an optimisation to not create goroutine if someone reads the channel (should cover 90% cases)
		select {
		case ch.ch <- v:
		default:
			ch.send.Add(1)
			go func() {
				defer ch.send.Done()
				select {
				case ch.ch <- v:
				case <-ch.done:
					ack.Done()
				}
			}()
		}
	})
	go func() {
		ack.Wait()
		value.Ack()
	}()
}
