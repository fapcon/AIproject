package kafka

import (
	"studentgit.kata.academy/quant/torque/internal/adapters/kafka/types/intev"
	"studentgit.kata.academy/quant/torque/pkg/kafka"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type ConsumersBuilder interface {
	InternalEvents() (Consumer[*intev.InternalEvent], error)
	GroupInternalEvents() (Consumer[*intev.InternalEvent], error)
}

type ConsumerConfigs struct {
	InternalEvents      *kafka.ConsumerConfig
	GroupInternalEvents *kafka.ConsumerConfig
}

type RealConsumersBuilder struct {
	logger logster.Logger
	cfg    ConsumerConfigs
}

func NewRealConsumersBuilder(logger logster.Logger, cfg ConsumerConfigs) *RealConsumersBuilder {
	return &RealConsumersBuilder{
		logger: logger,
		cfg:    cfg,
	}
}
