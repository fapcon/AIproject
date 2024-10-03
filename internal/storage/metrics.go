package storage

import (
	"time"

	"studentgit.kata.academy/quant/torque/internal/domain/observation"
)

func (s *Storage) ObserveDuration(observation observation.Name, elapsed time.Duration) {
	s.prometheus.ObserveHistogram(observation, elapsed.Seconds())
}

func (s *Storage) ObserveValue(observation observation.Name, value float64) {
	s.prometheus.ObserveHistogram(observation, value)
}

func (s *Storage) SetValue(observation observation.Name, value float64) {
	s.prometheus.SetGauge(observation, value)
}

func (s *Storage) SetLabelValue(observation observation.Name, label string, value float64) {
	s.prometheus.SetLabelGauge(observation, label, value)
}

func (s *Storage) AddToCounter(o observation.Name, inc int) {
	s.prometheus.IncreaseCounter(o, uint(inc))
}

func (s *Storage) IncCounter(o observation.Name) {
	s.prometheus.IncreaseCounter(o, 1)
}

func (s *Storage) IncLabelCounter(o observation.Name, label string) {
	s.prometheus.IncreaseLabelCounter(o, label, 1)
}
