package kafka

import (
	"context"
	"fmt"
	"io"

	"github.com/IBM/sarama"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"studentgit.kata.academy/quant/torque/pkg/fp/chans"
	"studentgit.kata.academy/quant/torque/pkg/kafka"
	"studentgit.kata.academy/quant/torque/pkg/logster"

	"studentgit.kata.academy/quant/torque/internal/domain"
)

type Consumer[T any] interface {
	Messages() <-chan kafka.Ackable[T]
	WaitNewestAcked(context.Context) error
	Init(context.Context) error
	Run(context.Context) error
	Topic() string
	Close() error
}

type Messages[T any] <-chan domain.Ackable[T]

func (c Messages[T]) Messages() <-chan domain.Ackable[T] {
	return c
}

type Producer[K, T any] interface {
	Send(K, T)
	Topic() string
	Close() error
}

type SyncProducer[K, T any] interface {
	Send(K, T) error
	Topic() string
	Close() error
}

type consumers struct {
	InternalEvents      *InternalEventConsumer
	GroupInternalEvents *InternalEventConsumer
}

type producers struct {
	InternalEvents *InternalEventProducer
}

type syncProducers struct {
}

type Kafka struct {
	logger           logster.Logger
	consumersBuilder ConsumersBuilder
	producersBuilder ProducersBuilder
	registered       map[any]struct{}

	inits  []func(context.Context) error
	runs   []func(context.Context) error
	closes []func() error

	Consumers     consumers
	Producers     producers
	SyncProducers syncProducers
}

func NewKafka(
	logger logster.Logger,
	consumersBuilder ConsumersBuilder,
	producersBuilder ProducersBuilder,
) *Kafka {
	return &Kafka{ //nolint:exhaustruct // other fields are filled later
		logger:           logger,
		consumersBuilder: consumersBuilder,
		producersBuilder: producersBuilder,
		registered:       make(map[any]struct{}),
	}
}

// Init should be called after all registrations.
func (k *Kafka) Init(ctx context.Context) error {
	for _, f := range k.inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Kafka) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, f := range k.runs {
		f := f
		eg.Go(func() error {
			return f(ctx)
		})
	}
	return eg.Wait()
}

func (k *Kafka) Close() error {
	var errs error
	for _, closeFn := range k.closes {
		if err := closeFn(); err != nil {
			multierr.AppendInto(&errs, fmt.Errorf("Kafka.Close(): %w", err))
		}
	}
	return errs
}

func (k *Kafka) once(target any, f func() error) error {
	if _, ok := k.registered[target]; ok {
		return nil
	}

	k.registered[target] = struct{}{}

	return f()
}

func registerCustomConsumer[T any, A any, C ~<-chan A](
	k *Kafka,
	field *C,
	consumer Consumer[T],
	convert func(kafka.Ackable[T]) A,
) {
	*field = chans.Map(consumer.Messages(), convert)
	k.inits = append(k.inits, consumer.Init)
	k.runs = append(k.runs, consumer.Run)
	k.closes = append(k.closes, consumer.Close)
}

func registerProducer[T io.Closer](k *Kafka, field *T, build func() (T, error)) error {
	producer, err := build()
	if err != nil {
		return err
	}
	*field = producer
	k.closes = append(k.closes, producer.Close)
	return nil
}

//nolint:nonamedreturns // for clarity
func marshalProto[K ~string, M proto.Message](k K, msg M) (key, value []byte, _ []sarama.RecordHeader, _ error) {
	value, err := proto.Marshal(msg)
	if err != nil {
		return nil, nil, nil, err
	}
	return []byte(k), value, nil, nil
}

func ChanToDomain[
	T interface{ ToDomain() A },
	A domain.Ackable[D],
	D any,
](
	in <-chan T,
) <-chan A {
	return chans.Map(in, func(t T) A {
		return t.ToDomain()
	})
}
