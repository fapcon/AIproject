package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"studentgit.kata.academy/quant/torque/pkg/health"
)

func DefaultTechOptions(logger promhttp.Logger, registry prometheus.Gatherer) RouterOption {
	return RouterOptions(
		WithRecover(),
		WithReadinessHandler(),
		WithDebugHandler(),
		WithMetricsHandler(logger, registry),
	)
}

func RouterOptions(options ...RouterOption) func(chi.Router) {
	return func(r chi.Router) {
		for _, option := range options {
			option(r)
		}
	}
}

type RouterOption func(chi.Router)

func WithReadinessHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/healthcheck", health.Routes())
	}
}

func WithDebugHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/debug", middleware.Profiler())
	}
}

func WithMetricsHandler(logger promhttp.Logger, registry prometheus.Gatherer) RouterOption {
	return func(r chi.Router) {
		r.Mount("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: logger}))
	}
}

// WithRecover adds recover middleware, which can catch panics from handlers.
func WithRecover() RouterOption {
	return func(r chi.Router) {
		r.Use(middleware.Recoverer)
	}
}
