package prometheus

import (
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"studentgit.kata.academy/quant/torque/internal/domain/observation"
	"studentgit.kata.academy/quant/torque/pkg/fp/maps"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const (
	Gateway = "gateway"
)

type (
	observationsSet  map[observation.Name]*sync.Once
	countersMap      map[observation.Name]prometheus.Counter
	countersVecMap   map[observation.Name]*prometheus.CounterVec
	histogramsMap    map[observation.Name]prometheus.Histogram
	histogramsVecMap map[observation.Name]*prometheus.HistogramVec
	gaugesMap        map[observation.Name]prometheus.Gauge
	gaugesVecMap     map[observation.Name]*prometheus.GaugeVec
)

type Collector struct {
	logger logster.Logger

	observations  observationsSet
	counters      countersMap
	countersVec   countersVecMap
	histograms    histogramsMap
	histogramsVec histogramsVecMap
	gauges        gaugesMap
	gaugesVec     gaugesVecMap

	registerer prometheus.Registerer
}

func NewCollector(logger logster.Logger, r prometheus.Registerer, subsystem string) (*Collector, error) {
	counters := countersMap{
		observation.RestAccountInfoRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Gateway,
			Subsystem: subsystem,
			Name:      observation.RestAccountInfoRequest,
			Help:      "Number of account info requests",
		}),
		observation.RestOpenOrdersRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Gateway,
			Subsystem: subsystem,
			Name:      observation.RestOpenOrdersRequest,
			Help:      "Number of open orders requests",
		}),
		observation.RestOrderCancelRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Gateway,
			Subsystem: subsystem,
			Name:      observation.RestOrderCancelRequest,
			Help:      "Number of cancel orders requests",
		}),
		observation.RestOrderRequest: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: Gateway,
			Subsystem: subsystem,
			Name:      observation.RestOrderRequest,
			Help:      "Number of order requests",
		}),
	}

	countersVec := countersVecMap{
		observation.MakerDataUpdate: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Gateway,
			Subsystem: subsystem,
			Name:      observation.MakerDataUpdate,
			Help:      "Number of market data updates requests",
		}, []string{"symbol"}),
	}

	histograms := histogramsMap{}

	histogramsVec := histogramsVecMap{}

	gauges := gaugesMap{}

	gaugesVec := gaugesVecMap{
		observation.OrderBook: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Gateway,
			Subsystem: subsystem,
			Name:      observation.OrderBook,
			Help:      "Number of bids and asks in order book",
		}, []string{"symbol"}),
	}

	observations := slices.ToMap1(slices.Merge(
		maps.Keys(counters),
		maps.Keys(countersVec),
		maps.Keys(histograms),
		maps.Keys(histogramsVec),
		maps.Keys(gaugesVec),
		maps.Keys(gauges),
	), func(e observation.Name) (observation.Name, *sync.Once) {
		return e, &sync.Once{}
	})

	if len(counters)+len(countersVec)+
		len(histograms)+len(histogramsVec)+len(gauges)+len(gaugesVec) > len(observations) {
		return nil, errors.New("multiple initialization of one of observations")
	}

	return &Collector{
		logger:        logger,
		observations:  observations,
		counters:      counters,
		countersVec:   countersVec,
		histograms:    histograms,
		histogramsVec: histogramsVec,
		gauges:        gauges,
		gaugesVec:     gaugesVec,
		registerer:    r,
	}, nil
}

func (c *Collector) IncreaseCounter(observation observation.Name, inc uint) {
	if inc == 0 {
		return
	}
	if counter, ok := getMetric(c, c.counters, observation); ok {
		counter.Add(float64(inc))
	}
}

func (c *Collector) IncreaseLabelCounter(observation observation.Name, label string, inc uint) {
	if inc == 0 {
		return
	}
	if counter, ok := getMetric(c, c.countersVec, observation); ok {
		counter.WithLabelValues(label).Add(float64(inc))
	}
}

func (c *Collector) ObserveHistogram(observation observation.Name, value float64) {
	if histogram, ok := getMetric(c, c.histograms, observation); ok {
		histogram.Observe(value)
	}
}

func (c *Collector) SetGauge(observation observation.Name, value float64) {
	if gauge, ok := getMetric(c, c.gauges, observation); ok {
		gauge.Set(value)
	}
}

func (c *Collector) SetLabelGauge(observation observation.Name, label string, value float64) {
	if gauge, ok := getMetric(c, c.gaugesVec, observation); ok {
		gauge.WithLabelValues(label).Set(value)
	}
}

func getMetric[T prometheus.Collector](c *Collector, m map[observation.Name]T, observation observation.Name) (T, bool) {
	metric, ok1 := m[observation]
	once, ok2 := c.observations[observation]
	if !ok1 || !ok2 {
		c.logger.WithField("observation", observation).Errorf("Unknown observation")
		var zero T
		return zero, false
	}
	once.Do(func() {
		err := c.registerer.Register(metric)
		if err != nil {
			c.logger.WithField("observation", observation).WithError(err).Errorf("Can't register metric")
		}
	})
	return metric, true
}
