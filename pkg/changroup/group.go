package changroup

import (
	"sync"
)

type ReleaseFunc func()

type channel[T any] struct {
	ch      chan T
	done    chan struct{}
	send    sync.WaitGroup
	release ReleaseFunc
}

type Group[T any] struct {
	channels *list[*channel[T]]
}

func NewGroup[T any]() *Group[T] {
	return &Group[T]{
		channels: newList[*channel[T]](),
	}
}

// ReleaseAll releases all acquired channels and closes it.
// It's safe to call ReleaseAll several times as well as in parallel with ReleaseFunc.
func (g *Group[T]) ReleaseAll() {
	var all []*channel[T]
	g.channels.ForEach(func(ch *channel[T]) {
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
func (g *Group[T]) Acquire() (<-chan T, ReleaseFunc) {
	ch := g.channels.Insert(&channel[T]{
		ch:      make(chan T),
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

// Send sends value to all acquired channels.
// It guarantees that all channels receive values in the same order.
// And that order is the same as Send calls.
// It may block if some channel is not read until its release.
func (g *Group[T]) Send(value T) {
	wg := sync.WaitGroup{}
	g.channels.ForEach(func(ch *channel[T]) {
		// select is an optimisation to not create goroutine if someone reads the channel (should cover 90% cases)
		select {
		case ch.ch <- value:
		default:
			wg.Add(1)
			ch.send.Add(1)
			go func() {
				defer wg.Done()
				defer ch.send.Done()
				select {
				case ch.ch <- value:
				case <-ch.done:
				}
			}()
		}
	})
	wg.Wait()
}

// SendAsync sends value to all acquired channels, but unlike [Group.Send] doesn't block.
// Also, it doesn't preserve the order of values!
func (g *Group[T]) SendAsync(value T) {
	g.channels.ForEach(func(ch *channel[T]) {
		// select is an optimisation to not create goroutine if someone reads the channel (should cover 90% cases)
		select {
		case ch.ch <- value:
		default:
			ch.send.Add(1)
			go func() {
				defer ch.send.Done()
				select {
				case ch.ch <- value:
				case <-ch.done:
				}
			}()
		}
	})
}
