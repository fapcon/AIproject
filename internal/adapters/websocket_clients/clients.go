package websocketclients

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/marketdata"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/trading"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/userstream"

	"studentgit.kata.academy/quant/torque/config"

	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type Clients struct {
	logger                    logster.Logger
	OKXSpotTradingClient      *trading.OKXTradingWebsocket
	OKXMarketDataWebsocket    *marketdata.OKXMarketDataWebsocket
	OKXUserStreamWebsocket    *userstream.OKXUserStreamWebsocket
	BybitMarketDataWebsocket  *marketdata.BybitMarketDataWebsocket
	BybitUserStreamWebsocket  *userstream.BybitUserStreamWebsocket
	GateIOMarketDataWebsocket *marketdata.GateIOMarketDataWebsocket
}

func NewWebsocketClients(logger logster.Logger) *Clients {
	return &Clients{
		logger: logger,
		// fields below are filled during registration
		OKXSpotTradingClient:      nil,
		OKXMarketDataWebsocket:    nil,
		OKXUserStreamWebsocket:    nil,
		BybitMarketDataWebsocket:  nil,
		BybitUserStreamWebsocket:  nil,
		GateIOMarketDataWebsocket: nil,
	}
}

func (c *Clients) RegisterOKXMarketDataWebsocketClient(config config.MarketDataConfig) {
	c.OKXMarketDataWebsocket = marketdata.NewOKXMarketDataWebsocket(c.logger, config)
}

func (c *Clients) RegisterGateIOMarketDataWebsocketClient(config config.MarketDataConfig) {
	c.GateIOMarketDataWebsocket = marketdata.NewGateIOMarketDataWebsocket(c.logger, config)
}

func (c *Clients) RegisterOKXUserStreamWebsocketClient(config config.UserStreamConfig) {
	c.OKXUserStreamWebsocket = userstream.NewOkxUserStreamWebsocket(c.logger, config)
}

func (c *Clients) RegisterOKXTradingWebsocketClient(config config.UserStreamConfig) {
	c.OKXSpotTradingClient = trading.NewOkxTradingWebsocket(c.logger, config)
}

func (c *Clients) RegisterBybitMarketDataWebsocketClient(config config.MarketDataConfig) {
	c.BybitMarketDataWebsocket = marketdata.NewBybitMarketDataWebsocket(c.logger, config)
}

func (c *Clients) RegisterBybitUserStreamWebsocketClient(config config.UserStreamConfig) {
	c.BybitUserStreamWebsocket = userstream.NewBybitUserStreamWebsocket(c.logger, config)
}

func (c *Clients) Close() {

}

func (c *Clients) Run(ctx context.Context) error {
	if c.OKXMarketDataWebsocket != nil {
		err := c.OKXMarketDataWebsocket.Run(ctx)
		if err != nil {
			return err
		}
	}
	if c.GateIOMarketDataWebsocket != nil {
		err := c.GateIOMarketDataWebsocket.Run(ctx)
		if err != nil {
			return err
		}
	}
	if c.OKXUserStreamWebsocket != nil {
		err := c.OKXUserStreamWebsocket.Run(ctx)
		if err != nil {
			return err
		}
	}
	if c.OKXSpotTradingClient != nil {
		err := c.OKXSpotTradingClient.Run(ctx)
		if err != nil {
			return err
		}
	}
	if c.BybitMarketDataWebsocket != nil {
		err := c.BybitMarketDataWebsocket.Run(ctx)
		if err != nil {
			return err
		}
	}
	if c.BybitUserStreamWebsocket != nil {
		err := c.BybitUserStreamWebsocket.Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Clients) Init(ctx context.Context) error {
	if c.OKXMarketDataWebsocket != nil {
		err := c.OKXMarketDataWebsocket.Init(ctx)
		if err != nil {
			return err
		}
	}
	if c.GateIOMarketDataWebsocket != nil {
		err := c.GateIOMarketDataWebsocket.Init(ctx)
		if err != nil {
			return err
		}
	}
	if c.OKXUserStreamWebsocket != nil {
		err := c.OKXUserStreamWebsocket.Init(ctx)
		if err != nil {
			return err
		}
	}
	if c.OKXSpotTradingClient != nil {
		err := c.OKXSpotTradingClient.Init(ctx)
		if err != nil {
			return err
		}
	}
	if c.BybitMarketDataWebsocket != nil {
		err := c.BybitMarketDataWebsocket.Init(ctx)
		if err != nil {
			return err
		}
	}
	if c.BybitUserStreamWebsocket != nil {
		err := c.BybitUserStreamWebsocket.Init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
