package kafka

import (
	"studentgit.kata.academy/quant/torque/internal/adapters/kafka/types/intev"
	"studentgit.kata.academy/quant/torque/pkg/kafka"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type ProducersBuilder interface {
	InternalEvents() (Producer[string, *intev.InternalEvent], error)
}

type ProducerConfigs struct {
	InternalEvents *kafka.ProducerConfig
}

type RealProducersBuilder struct {
	logger logster.Logger
	cfg    ProducerConfigs
}

func NewRealProducersBuilder(logger logster.Logger, cfg ProducerConfigs) *RealProducersBuilder {
	return &RealProducersBuilder{
		logger: logger,
		cfg:    cfg,
	}
}
