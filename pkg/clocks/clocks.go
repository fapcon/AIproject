package clocks

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Clock should be used instead of [time.Now], [time.After],
// [time.NewTimer], [time.Tick] and [time.NewTicker].
// This allows to mock time in tests.
type Clock interface {
	// Now is a mockable [time.Now].
	Now() time.Time

	// After is a mockable [time.After]; is not removed by GC until timer fires.
	After(fireAfter time.Duration) <-chan time.Time

	// NewTimer is a mockable [time.NewTimer].
	NewTimer(fireAfter time.Duration) (<-chan time.Time, func() bool)

	// Tick is a mockable [time.Tick]; it leaks.
	Tick(fireEach time.Duration) <-chan time.Time

	// NewTicker is a mockable [time.NewTicker].
	NewTicker(fireEach time.Duration) (<-chan time.Time, func())
}

// RealClock uses [time.Now], [time.NewTimer] and [time.NewTicker].
// Should be used in production.
type RealClock struct{}

func (c RealClock) Now() time.Time {
	return time.Now()
}

func (c RealClock) After(fireAfter time.Duration) <-chan time.Time {
	return time.After(fireAfter)
}

func (c RealClock) NewTimer(fireAfter time.Duration) (<-chan time.Time, func() bool) {
	t := time.NewTimer(fireAfter)
	return t.C, t.Stop
}

func (c RealClock) Tick(fireEach time.Duration) <-chan time.Time {
	return time.Tick(fireEach) //nolint:staticcheck // yes, it leaks
}

func (c RealClock) NewTicker(fireEach time.Duration) (<-chan time.Time, func()) {
	t := time.NewTicker(fireEach)
	return t.C, t.Stop
}

type testTimer struct {
	fireAt  time.Time
	ch      chan time.Time
	stopped *atomic.Bool
}

type testTicker struct {
	lastFiredAt time.Time
	interval    time.Duration
	ch          chan time.Time
	stopped     *atomic.Bool
}

// TestClock is a mock for [Clock].
// Should be used in tests.
type TestClock struct {
	mu      sync.Mutex
	now     time.Time
	timers  map[int]testTimer
	counter int
	tickers []testTicker
}

// NewTestClock returns new mock for [Clock].
// Should be used in tests.
func NewTestClock(t time.Time) *TestClock {
	return &TestClock{
		mu:      sync.Mutex{},
		now:     t,
		timers:  make(map[int]testTimer),
		counter: 0,
		tickers: nil,
	}
}

func (c *TestClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *TestClock) After(fireAfter time.Duration) <-chan time.Time {
	ch, _ := c.NewTimer(fireAfter)
	return ch
}

func (c *TestClock) NewTimer(fireAfter time.Duration) (<-chan time.Time, func() bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan time.Time, 1)
	stopped := &atomic.Bool{}

	if fireAfter <= 0 {
		ch <- c.now
		return ch, func() bool { return false }
	}

	c.counter++
	c.timers[c.counter] = testTimer{
		fireAt:  c.now.Add(fireAfter),
		ch:      ch,
		stopped: stopped,
	}

	stop := func() bool {
		return stopped.CompareAndSwap(false, true)
	}

	return ch, stop
}

func (c *TestClock) Tick(fireEach time.Duration) <-chan time.Time {
	ch, _ := c.NewTicker(fireEach)
	return ch
}

func (c *TestClock) NewTicker(fireEach time.Duration) (<-chan time.Time, func()) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan time.Time, 1)
	stopped := &atomic.Bool{}

	c.tickers = append(c.tickers, testTicker{
		lastFiredAt: c.now,
		interval:    fireEach,
		ch:          ch,
		stopped:     stopped,
	})

	stop := func() {
		stopped.Store(true)
	}

	return ch, stop
}

// Set moves TestClock forward to requested moment in time.
// Panics if t is before now.
func (c *TestClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if t.Before(c.now) {
		panic(fmt.Sprintf("you can not travel back in time, now: %s, destination: %s", c.now, t))
	}
	c.now = t

	for i, timer := range c.timers {
		if !timer.fireAt.After(c.now) {
			timer.stopped.Store(true)
			timer.ch <- c.now   // no need to select with default because buffer + delete
			delete(c.timers, i) // yes, it's safe to delete from map while iterating over it
		}
	}

	for i, ticker := range c.tickers {
		if !ticker.stopped.Load() && !ticker.lastFiredAt.Add(ticker.interval).After(c.now) {
			//nolint:durationcheck // ok to multiply durations
			c.tickers[i].lastFiredAt = ticker.lastFiredAt.Add(
				(c.now.Sub(ticker.lastFiredAt) / ticker.interval) * ticker.interval,
			)
			select {
			case ticker.ch <- c.now:
			default:
			}
		}
	}
}

// Forward moves TestClock forward for requested duration.
// Panics if negative duration is passed.
func (c *TestClock) Forward(d time.Duration) {
	c.Set(c.now.Add(d))
}

func GetSecondsTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}
