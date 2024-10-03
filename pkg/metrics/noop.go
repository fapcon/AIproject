package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

func NewNoopCollector() *Collector {
	c, _ := NewCollector(nil, nil)
	return c
}

// noopMetric implements Counter, Gauge, Histogram and Summary.
type noopMetric struct{}

func (n noopMetric) Observe(float64) {}

func (n noopMetric) Desc() *prometheus.Desc { return nil }

func (n noopMetric) Write(*io_prometheus_client.Metric) error { return nil }

func (n noopMetric) Describe(chan<- *prometheus.Desc) {}

func (n noopMetric) Collect(chan<- prometheus.Metric) {}

func (n noopMetric) Set(float64) {}

func (n noopMetric) Inc() {}

func (n noopMetric) Dec() {}

func (n noopMetric) Add(float64) {}

func (n noopMetric) Sub(float64) {}

func (n noopMetric) SetToCurrentTime() {}
