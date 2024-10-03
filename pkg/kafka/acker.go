package kafka

import (
	"context"
	"sync"
)

type ackedOffset struct {
	offset int64
	acked  bool
}

// Acker converts independent (and parallel) acks into monotonously increasing acks.
type Acker struct {
	ack     func(int64)
	mu      sync.Mutex
	queue   []ackedOffset
	waiting chan struct{}
}

func NewAcker(ack func(int64)) *Acker {
	return &Acker{
		ack:     ack,
		mu:      sync.Mutex{},
		queue:   nil,
		waiting: nil,
	}
}

// Processing puts offset into queue to be acked in the future.
// Typical scenario assumes that offsets are monotonously increasing.
//
// However, it isn't must, because the queue is kept sorted.
// In this case [Acker] doesn't guarantee to call [Acker.ack] in increasing order.
//
// Note: some offsets may be skipped (example of good sequence: 1, 2, 5, 6, 8).
func (a *Acker) Processing(offset int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ao := ackedOffset{
		offset: offset,
		acked:  false,
	}
	if len(a.queue) == 0 || offset >= a.queue[len(a.queue)-1].offset {
		// normal mode: monotonously increasing offsets
		a.queue = append(a.queue, ao)
		return
	}
	if offset <= a.queue[0].offset {
		a.queue = append([]ackedOffset{ao}, a.queue...)
		return
	}
	for i := len(a.queue) - 1; i > 0; i-- {
		if offset >= a.queue[i-1].offset {
			a.queue = append(a.queue[:i], append([]ackedOffset{ao}, a.queue[i:]...)...)
			return
		}
	}
}

// Processed finds offset in queue and marks it acked.
// If queue begins with N consecutive acked offsets,
// N-th offset is acked by [Acker.ack] and N first offsets are removed from queue.
func (a *Acker) Processed(offset int64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i, ao := range a.queue {
		if ao.offset > offset {
			return // already acked or unknown offset
		}
		if ao.offset == offset {
			a.queue[i] = ackedOffset{
				offset: offset,
				acked:  true,
			}
			if i == 0 {
				a.ackAndRemoveFromQueue()
			}
			return
		}
	}
}

// ackAndRemoveFromQueue must be called under a mutex and only if first offset is acked in queue.
func (a *Acker) ackAndRemoveFromQueue() {
	var maxOffset int64
	var idx int
	for i, ao := range a.queue {
		if ao.acked {
			maxOffset = ao.offset
			idx = i
		} else {
			break
		}
	}

	a.ack(maxOffset)
	a.queue = a.queue[idx+1:]

	if a.waiting != nil && len(a.queue) == 0 {
		close(a.waiting)
	}
}

// Wait blocks until queue is empty (i.e. all offsets are acked) or context is done.
func (a *Acker) Wait(ctx context.Context) {
	a.mu.Lock()
	if len(a.queue) == 0 {
		a.mu.Unlock()
		return
	}

	a.waiting = make(chan struct{})
	a.mu.Unlock()

	select {
	case <-a.waiting:
	case <-ctx.Done():
	}
}
