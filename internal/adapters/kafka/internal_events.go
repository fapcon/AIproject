package kafka

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"studentgit.kata.academy/quant/torque/internal/domain"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
	"studentgit.kata.academy/quant/torque/internal/adapters/kafka/types/intev"
	"studentgit.kata.academy/quant/torque/pkg/kafka"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const producerHeader = "producer"

var /* const */ self = []byte(strconv.Itoa(rand.Int())) //nolint:gochecknoglobals,gosec // const and no need crypto

func withProducerHeader[K ~string, M proto.Message](
	fn func(K, M) ([]byte, []byte, []sarama.RecordHeader, error),
) func(K, M) ([]byte, []byte, []sarama.RecordHeader, error) {
	//nolint:nonamedreturns // needed for clarity to distinguish what is what
	return func(k K, msg M) (key, value []byte, _ []sarama.RecordHeader, _ error) {
		key, value, headers, err := fn(k, msg)
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(producerHeader),
			Value: self,
		})
		return key, value, headers, err
	}
}

func (k *Kafka) RegisterInternalEventsConsumer() error {
	return k.registerInternalEventsConsumer(&k.Consumers.InternalEvents, k.consumersBuilder.InternalEvents)
}

func (k *Kafka) RegisterGroupInternalEventsConsumer() error {
	return k.registerInternalEventsConsumer(&k.Consumers.GroupInternalEvents, k.consumersBuilder.GroupInternalEvents)
}

func (k *Kafka) registerInternalEventsConsumer(
	c **InternalEventConsumer,
	consumerBuilder func() (Consumer[*intev.InternalEvent], error),
) error {
	return k.once(c, func() error {
		consumer, err := consumerBuilder()
		if err != nil {
			return err
		}

		var internalEvents <-chan intev.AckableInternalEvent
		registerCustomConsumer(
			k,
			&internalEvents,
			consumer,
			func(a kafka.Ackable[*intev.InternalEvent]) intev.AckableInternalEvent {
				return intev.AckableInternalEvent{
					InternalEvent: a.Message,
					Ack:           a.Ack,
				}
			},
		)

		internalEventsDispatcher := intev.NewInternalEventDispatcher(k.logger)
		*c = NewInternalEventConsumer(k.logger, internalEventsDispatcher, consumer)

		k.runs = append(k.runs, func(context.Context) error {
			internalEventsDispatcher.Run(internalEvents)
			return nil
		})

		return nil
	})
}

func decodeInternalEvents(_ context.Context, msg *sarama.ConsumerMessage) (*intev.InternalEvent, kafka.Mark, error) {
	event := new(intev.InternalEvent)
	if err := proto.Unmarshal(msg.Value, event); err != nil {
		return nil, kafka.Skip, fmt.Errorf("can't Unmarshal intev.InternalEvent: %w", err)
	}
	return event, kafka.Process, nil
}

func (b *RealConsumersBuilder) InternalEvents() (Consumer[*intev.InternalEvent], error) {
	decode := func(ctx context.Context, msg *sarama.ConsumerMessage) (*intev.InternalEvent, kafka.Mark, error) {
		for _, header := range msg.Headers {
			if string(header.Key) == producerHeader && string(header.Value) == string(self) {
				return nil, kafka.Skip, nil
			}
		}
		return decodeInternalEvents(ctx, msg)
	}
	return kafka.NewConsumer(b.logger, b.cfg.InternalEvents, decode)
}

func (b *RealConsumersBuilder) GroupInternalEvents() (Consumer[*intev.InternalEvent], error) {
	return kafka.NewConsumer(b.logger, b.cfg.GroupInternalEvents, decodeInternalEvents)
}

type InternalEventConsumer struct {
	logger     logster.Logger
	dispatcher *intev.InternalEventDispatcher
	consumer   Consumer[*intev.InternalEvent]
}

func NewInternalEventConsumer(
	logger logster.Logger,
	dispatcher *intev.InternalEventDispatcher,
	consumer Consumer[*intev.InternalEvent],
) *InternalEventConsumer {
	return &InternalEventConsumer{
		logger:     logger,
		dispatcher: dispatcher,
		consumer:   consumer,
	}
}

func (c *InternalEventConsumer) WaitNewestAcked(ctx context.Context) error {
	return c.consumer.WaitNewestAcked(ctx)
}

func (k *Kafka) RegisterInternalEventsProducer() error {
	return k.once(&k.Producers.InternalEvents, func() error {
		var internalEvents Producer[string, *intev.InternalEvent]
		err := registerProducer(k, &internalEvents, k.producersBuilder.InternalEvents)
		if err != nil {
			return err
		}

		k.Producers.InternalEvents = NewInternalEventProducer(internalEvents)
		return nil
	})
}

func (b *RealProducersBuilder) InternalEvents() (Producer[string, *intev.InternalEvent], error) {
	return kafka.NewProducer(
		b.logger,
		*b.cfg.InternalEvents,
		withProducerHeader(marshalProto[string, *intev.InternalEvent]),
	)
}

type InternalEventProducer struct {
	producer Producer[string, *intev.InternalEvent]
}

func NewInternalEventProducer(producer Producer[string, *intev.InternalEvent]) *InternalEventProducer {
	return &InternalEventProducer{
		producer: producer,
	}
}

func (p *InternalEventProducer) SendSubscribeInstrumentsRequested(d []domain.Instrument) {
	p.producer.Send(
		"instrument_requests",
		&intev.InternalEvent{Event: &intev.InternalEvent_SubscribeInstruments{
			SubscribeInstruments: intev.SubscribeInstrumentsRequestedFromDomain(d),
		}},
	)
}

func (c *InternalEventConsumer) NewOrderRequested() <-chan domain.Ackable[domain.NewOrderRequest] {
	return ChanToDomain[intev.AckableNewOrderRequested, domain.Ackable[domain.NewOrderRequest]](
		c.dispatcher.NewOrderRequested(),
	)
}
