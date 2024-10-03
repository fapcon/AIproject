package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Name string

type Collector struct {
	empty   bool
	metrics map[Name]prometheus.Collector
}

// NewCollector creates new metrics collector. If r is nil, then all metrics getters are be noop.
func NewCollector(
	r prometheus.Registerer,
	metrics map[Name]prometheus.Collector,
) (*Collector, error) {
	// register metrics only if r is not nil
	if r == nil {
		return &Collector{
			empty:   true,
			metrics: nil,
		}, nil
	}

	for name, metric := range metrics {
		// validating incoming collector
		//nolint:gocritic // false positive
		switch m := metric.(type) {
		case prometheus.Counter,
			*prometheus.CounterVec,
			prometheus.Gauge,
			*prometheus.GaugeVec,
			prometheus.Histogram,
			*prometheus.HistogramVec,
			prometheus.Summary,
			*prometheus.SummaryVec:
		default:
			return nil, fmt.Errorf("unknown metric type %s: %T", name, m)
		}
		err := r.Register(metric)
		if err != nil {
			return nil, fmt.Errorf("can't register metric %s: %w", name, err)
		}
	}

	return &Collector{
		empty:   false,
		metrics: metrics,
	}, nil
}

func (c *Collector) GetHistogram(name Name, labelValues ...string) prometheus.Observer {
	if c.empty {
		return noopMetric{}
	}

	metric, ok := c.metrics[name]
	if !ok {
		return noopMetric{} // unknown metric
	}

	switch m := metric.(type) {
	case *prometheus.HistogramVec:
		return m.WithLabelValues(labelValues...)
	case prometheus.Histogram:
		return m
	default:
		return noopMetric{} // unknown type
	}
}

func (c *Collector) GetSummary(name Name, labelValues ...string) prometheus.Observer {
	if c.empty {
		return noopMetric{}
	}

	metric, ok := c.metrics[name]
	if !ok {
		return noopMetric{} // unknown metric
	}

	switch m := metric.(type) {
	case *prometheus.SummaryVec:
		return m.WithLabelValues(labelValues...)
	case prometheus.Summary:
		return m
	default:
		return noopMetric{} // unknown type
	}
}

func (c *Collector) GetCounter(name Name, labelValues ...string) prometheus.Counter {
	if c.empty {
		return noopMetric{}
	}

	metric, ok := c.metrics[name]
	if !ok {
		return noopMetric{} // unknown metric
	}

	switch m := metric.(type) {
	case *prometheus.CounterVec:
		return m.WithLabelValues(labelValues...)
	case prometheus.Counter:
		return m
	default:
		return noopMetric{} // unknown type
	}
}

func (c *Collector) GetGauge(name Name, labelValues ...string) prometheus.Gauge {
	if c.empty {
		return noopMetric{}
	}

	metric, ok := c.metrics[name]
	if !ok {
		return noopMetric{} // unknown metric
	}

	switch m := metric.(type) {
	case *prometheus.GaugeVec:
		return m.WithLabelValues(labelValues...)
	case prometheus.Gauge:
		return m
	default:
		return noopMetric{} // unknown type
	}
}
