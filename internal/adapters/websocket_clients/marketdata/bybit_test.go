package marketdata_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/marketdata"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/sig"
)

func TestNewBybitMarketDataWebsocket(t *testing.T) {
	type args struct {
		logger logster.Logger
		config config.MarketDataConfig
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "notNil",
			args: args{
				logger: logster.Discard,
				config: config.MarketDataConfig{
					URL:                    "wss://stream.bybit.com/v5/public/spot",
					InstrumentsToSubscribe: []string{"BTCUSDT"},
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := marketdata.NewBybitMarketDataWebsocket(tt.args.logger, tt.args.config)
			if got == nil && !tt.wantNil || got != nil && tt.wantNil {
				t.Errorf("NewByBitUserStreamWebsocket() == %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}
func TestMarketDataWebSocket(t *testing.T) {
	t.Parallel()
	logger := logster.Discard

	g, ctx := errgroup.WithContext(context.Background())

	defer logster.LogShutdownDuration(ctx, logger)()
	g.Go(func() error {
		return sig.ListenSignal(ctx, logger)
	})

	testConfig := config.MarketDataConfig{
		URL:                    "wss://stream.bybit.com/v5/public/spot",
		InstrumentsToSubscribe: []string{"BTCUSDT"},
	}

	websocket := marketdata.NewBybitMarketDataWebsocket(logger, testConfig)

	err := websocket.Init(ctx)
	assert.NoError(t, err)
	go func() {
		errRun := websocket.Run(ctx)
		assert.NoError(t, errRun)
	}()
	select {
	case orderBook, ok := <-websocket.GetOrderBookUpdate():
		if !ok {
			t.Error("The channel is closed")
			return
		}
		assert.NotNil(t, orderBook.Symbol)
		assert.NotEmpty(t, orderBook.Bids)
		assert.NotEmpty(t, orderBook.Asks)
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for order book update")
	}
}

func TestMarketDataWebSocketMock(t *testing.T) {
	t.Parallel()
	logger := logster.Discard

	g, ctx := errgroup.WithContext(context.Background())

	defer logster.LogShutdownDuration(ctx, logger)()
	g.Go(func() error {
		return sig.ListenSignal(ctx, logger)
	})
	//nolint:exhaustruct // don't need to fill all fields
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("error upgrading connection: %v", err)
			return
		}
		defer conn.Close()

		updateMsg := []byte(
			`{
				"topic":"orderbook.BTCUSD.50", 
				"type":"snapshot", 
				"ts":1615758851000, 
				"data":{
					"s":"BTCUSD",
					"b":[["50000","1"]],
					"a":[["50010","2"]],
					"u":123,
					"seq":456},
				"cts":1615758852000
				}`)
		err = conn.WriteMessage(websocket.TextMessage, updateMsg)
		if err != nil {
			t.Errorf("error writing message: %v", err)
			return
		}
	}))

	defer server.Close()

	testConfig := config.MarketDataConfig{
		URL:                    "ws" + strings.TrimPrefix(server.URL, "http"),
		InstrumentsToSubscribe: []string{"BTCUSDT"},
	}

	websocket := marketdata.NewBybitMarketDataWebsocket(logger, testConfig)

	err := websocket.Init(ctx)
	assert.NoError(t, err)

	go func() {
		errMock := websocket.Run(ctx)
		assert.NoError(t, errMock)
	}()
	select {
	case orderBook, ok := <-websocket.GetOrderBookUpdate():
		if !ok {
			return
		}
		assert.NotEmpty(t, orderBook.Type)
		assert.Equal(t, domain.SnapshotBybit, orderBook.Type)
		assert.NotNil(t, orderBook.Symbol)
		assert.Equal(t, domain.BybitSpotInstrument("BTCUSD"), orderBook.Symbol)
		assert.NotEmpty(t, orderBook.Bids)
		assert.Equal(t, map[string]float64{"50000": 1.0}, orderBook.Bids)
		assert.NotEmpty(t, orderBook.Asks)
		assert.Equal(t, map[string]float64{"50010": 2.0}, orderBook.Asks)
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for order book update")
	}
}
