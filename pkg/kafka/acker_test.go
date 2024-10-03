package kafka_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"studentgit.kata.academy/quant/torque/pkg/kafka"
)

func TestAcker(t *testing.T) {
	t.Parallel()
	t.Run("one message", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)

		ack.Expect(10)
		acker.Processed(10)
	})
	t.Run("two messages - same order", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		acker.Processing(20)

		ack.Expect(10)
		acker.Processed(10)

		ack.Expect(20)
		acker.Processed(20)
	})
	t.Run("two messages - reversed order", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		acker.Processing(20)

		// no ack
		acker.Processed(20)

		ack.Expect(20)
		acker.Processed(10)
	})
	t.Run("several messages - random order", func(t *testing.T) {
		t.Parallel()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		x := []int64{20, 30, 40, 50, 60, 70, 80}
		for i := range r.Perm(len(x)) {
			acker.Processing(x[i])
		}
		acker.Processing(90)
		acker.Processing(100)

		for i := range r.Perm(len(x)) {
			acker.Processed(x[i])
		}
		acker.Processed(100)

		ack.Expect(80)
		acker.Processed(10)

		ack.Expect(100)
		acker.Processed(90)
	})
	t.Run("ignored unknown offset", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		acker.Processing(20)
		acker.Processing(30)

		acker.Processed(5)
		acker.Processed(15)
		acker.Processed(25)
		acker.Processed(35)
	})
	t.Run("calling processed before processing does nothing", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)

		acker.Processed(20)

		ack.Expect(10)
		acker.Processed(10)

		acker.Processing(20)
	})
	t.Run("call processing with lowest offset", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(20)
		acker.Processing(30)
		acker.Processing(10)

		acker.Processed(20)

		ack.Expect(20)
		acker.Processed(10)
	})
	t.Run("call processing with not increasing offset", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		acker.Processing(30)
		acker.Processing(20)

		acker.Processed(20)

		ack.Expect(20)
		acker.Processed(10)
	})
}

func TestAcker_Wait(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)

		done := make(chan struct{})
		go func() {
			acker.Wait(context.Background())
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			assert.FailNow(t, "timeout")
		}
	})
	t.Run("already processed", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		ack.Expect(10)
		acker.Processed(10)

		done := make(chan struct{})
		go func() {
			acker.Wait(context.Background())
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			assert.FailNow(t, "timeout")
		}
	})
	t.Run("wait", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)
		acker.Processing(20)
		acker.Processing(30)
		acker.Processed(30)
		acker.Processed(20)

		done := make(chan struct{})
		go func() {
			acker.Wait(context.Background())
			close(done)
		}()

		time.Sleep(10 * time.Millisecond)
		select {
		case <-done:
			assert.FailNow(t, "already done")
		default:
		}

		ack.Expect(30)
		acker.Processed(10)

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			assert.FailNow(t, "timeout")
		}
	})
	t.Run("context.Done", func(t *testing.T) {
		t.Parallel()
		ack := NewAck(t)
		acker := kafka.NewAcker(ack.Ack)
		acker.Processing(10)

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() {
			acker.Wait(ctx)
			close(done)
		}()

		time.Sleep(10 * time.Millisecond)
		select {
		case <-done:
			assert.FailNow(t, "already done")
		default:
		}

		cancel()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			assert.FailNow(t, "timeout")
		}
	})
}

type Ack struct {
	m mock.Mock
}

func NewAck(t *testing.T) *Ack {
	m := new(Ack)
	m.m.Test(t)
	t.Cleanup(func() {
		m.m.AssertExpectations(t)
	})
	return m
}

func (m *Ack) Ack(offset int64) {
	m.m.Called(offset)
}

func (m *Ack) Expect(offset int64) {
	m.m.On("Ack", offset).Once()
}
