package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type okxProxyStorage interface {
	GetNewOrderRequests() <-chan domain.Ackable[domain.NewOrderRequest]
}

type OKXProxy struct {
	logger  logster.Logger
	storage okxProxyStorage

	newOrderRequests <-chan domain.Ackable[domain.NewOrderRequest]
}

func NewOKXProxy(logger logster.Logger, storage okxProxyStorage) *OKXProxy {
	return &OKXProxy{
		logger:  logger.WithField(f.Module, "okx_proxy"),
		storage: storage,

		newOrderRequests: storage.GetNewOrderRequests(),
	}
}

func (c *OKXProxy) Run(context.Context) error {
	g := errgroup.Group{}
	g.Go(processAckable(c.logger, "ProcessNewOrder", c.newOrderRequests, c.processNewOrderRequest))
	return g.Wait()
}

func (c *OKXProxy) processNewOrderRequest(_ context.Context, request domain.NewOrderRequest) error {
	c.logger.Infof("request %+v", request)
	return nil
}
