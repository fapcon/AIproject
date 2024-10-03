package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const ackTimeoutOnShutdownOrRebalanced = 10 * time.Second

type Mark bool

const (
	Process Mark = true
	Skip    Mark = false
)

type Ackable[T any] struct {
	Message T
	Ack     func()
}

type Handler[T any] struct {
	logger   logster.Logger
	messages chan<- Ackable[T] // messages must be non-buffered channel
	decode   func(context.Context, *sarama.ConsumerMessage) (T, Mark, error)
}

func NewHandler[T any](
	logger logster.Logger,
	messages chan<- Ackable[T],
	decode func(context.Context, *sarama.ConsumerMessage) (T, Mark, error),
) *Handler[T] {
	return &Handler[T]{
		logger:   logger,
		messages: messages,
		decode:   decode,
	}
}

func (p *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	p.logger.
		WithContext(session.Context()).
		WithField("claims", session.Claims()).
		Debugf("Handler: setup")
	return nil
}

func (p *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	p.logger.
		WithContext(session.Context()).
		WithField("claims", session.Claims()).
		Debugf("Handler: cleanup")
	return nil
}

func (p *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	topic := claim.Topic()
	partition := claim.Partition()

	base := p.logger.WithContext(ctx).WithField("topic", topic)

	acker := NewAcker(func(offset int64) {
		// +1 is taken from [sarama.ConsumerGroupSession.MarkMessage]
		session.MarkOffset(topic, partition, offset+1, "")
	})

	for message := range claim.Messages() {
		offset := message.Offset
		acker.Processing(offset)

		logger := base.WithFields(logster.Fields{
			"partition": partition,
			"offset":    offset,
			"timestamp": message.Timestamp,
			"key":       string(message.Key),
		})

		logger.Debugf("Handler: message received")

		msg, mark, err := p.decode(ctx, message)
		if err != nil {
			logger.WithError(err).Warnf("Handler: can't decode message")
		}
		if mark == Process {
			p.messages <- Ackable[T]{
				Message: msg,
				Ack: func() {
					acker.Processed(offset)
					logger.Debugf("Handler: message acked")
				},
			}
		} else {
			acker.Processed(offset)
			logger.Debugf("Handler: message skipped")
		}
	}

	ctx, cancel := context.WithTimeout(ctx, ackTimeoutOnShutdownOrRebalanced)
	defer cancel()
	acker.Wait(ctx)

	return nil
}
