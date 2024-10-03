package main

import (
	"context"
	"os"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"time"
)

//go:generate go run ./../../../../local_run/piston_server/main.go
func main() {
	logger := logster.New(os.Stdout, logster.Config{
		Project:           "x",
		Format:            "text",
		Level:             "debug",
		Env:               "local",
		DisableStackTrace: true,
		System:            "",
		Inst:              "",
	})

	logger.Infof("Starting server")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	piston := NewPistonClientMarketTest(logger, "localhost:5800")
	if err := piston.Init(ctx); err != nil {
		logger.WithError(err).Errorf("failed to init")
		return
	}

	// Receive messages.
	go func() {
		for msg := range piston.MarketDataSubscribeRequests {
			logger.Infof("Received message: %+v", msg)
		}
	}()

	// Send messages.
	go func() {
		for i := 0; i < 5; i++ {
			<-time.After(5 * time.Second)
			logger.Infof("Sending orderbook")
			err := piston.SendOrderBook(ctx, domain.OrderBook{
				Exchange:   domain.OKXExchange,
				Instrument: "BTC/USDT",
				Asks:       map[string]float64{"123": 1, "456": 2, "789": 3, "321": 4, "654": 5},
				Bids:       map[string]float64{"123": 1, "456": 2, "789": 3, "321": 4, "654": 5},
			})
			if err != nil {
				logger.WithError(err).Errorf("failed to send orderbook")
			}
		}
	}()

	// Run the piston.
	if err := piston.Run(ctx); err != nil {
		logger.WithError(err).Errorf("failed to run")
		return
	}

}
