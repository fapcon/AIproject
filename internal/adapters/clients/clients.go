package clients

import (
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"studentgit.kata.academy/quant/torque/internal/domain/f"

	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type Clients struct {
	logger   logster.Logger
	OKXApi   *OKXApi
	BybitAPI *BybitAPI
}

func NewClients(logger logster.Logger) *Clients {
	return &Clients{
		logger: logger,
		// fields below are filled during registration
		OKXApi:   nil,
		BybitAPI: nil,
	}
}

func (c *Clients) RegisterOKXAPIClient(config config.HTTPClient) {
	c.OKXApi = NewOKXAPI(c.logger, config)
}

func (c *Clients) RegisterBybitAPIClient(config config.HTTPClient) {
	c.BybitAPI = NewBybitAPI(c.logger, config)
}

func LogElapsed(logger logster.Logger, start time.Time, method, path string) {
	logger.WithField("elapsed_ms", time.Since(start).Milliseconds()).
		WithField(f.Method, method).WithField(f.Path, path).Debugf("Log time elapsed")
}
