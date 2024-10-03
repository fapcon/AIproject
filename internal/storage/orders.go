package storage

import (
	"studentgit.kata.academy/quant/torque/internal/adapters/kafka"
	"studentgit.kata.academy/quant/torque/internal/domain"
)

func (s *Storage) GetNewOrderRequests() <-chan domain.Ackable[domain.NewOrderRequest] {
	s.mustHave(CapabilityNewOrderRequestsConsumer)
	return s.kafka.Consumers.GroupInternalEvents.NewOrderRequested()
}

func (b *Builder) WithNewOrderRequestsConsumer() *Builder {
	b.buildCap(CapabilityNewOrderRequestsConsumer)
	b.buildKafka((*kafka.Kafka).RegisterGroupInternalEventsConsumer)
	return b
}
