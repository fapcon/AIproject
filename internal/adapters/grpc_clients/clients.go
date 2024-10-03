package grpcclients

import (
	"context"

	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type Clients struct {
	logger       logster.Logger
	PistonClient *PistonClient
}

func NewGrpcClients(logger logster.Logger) *Clients {
	return &Clients{
		logger: logger,
		// fields below are filled during registration
		PistonClient: nil,
	}
}

func (c *Clients) RegisterPistonClient(address string) {
	c.PistonClient = NewPistonClient(c.logger, address)
}

func (c *Clients) Run(ctx context.Context) error {
	if c.PistonClient != nil {
		err := c.PistonClient.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Clients) Init(ctx context.Context) error {
	if c.PistonClient != nil {
		err := c.PistonClient.Init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
