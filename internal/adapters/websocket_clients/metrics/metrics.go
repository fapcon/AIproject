package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"studentgit.kata.academy/quant/torque/internal/domain/observation"
	"studentgit.kata.academy/quant/torque/pkg/metrics"
)

func NewWebsocketMetrics() map[metrics.Name]prometheus.Collector {
	return map[metrics.Name]prometheus.Collector{
		observation.WebsocketReconnect: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: observation.Namespace,
			Subsystem: observation.Subsystem,
			Name:      string(observation.WebsocketReconnect),
			Help:      "Number of websocket reconnects",
		}, []string{"type"}),
		observation.WebsocketEventRoundTripDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: observation.Namespace,
			Subsystem: observation.Subsystem,
			Name:      string(observation.WebsocketEventRoundTripDuration),
			Help:      "Websocket event round trip duration microseconds",
			Buckets:   prometheus.ExponentialBuckets(0.1, 2, 5),
		}, []string{"event_type"}),
		observation.WebsocketEventDecodeDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: observation.Namespace,
			Subsystem: observation.Subsystem,
			Name:      string(observation.WebsocketEventDecodeDuration),
			Help:      "Websocket event decode duration microseconds",
			Buckets:   prometheus.ExponentialBuckets(100, 2, 15),
		}),
	}
}
