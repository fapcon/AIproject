package marketdata_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	wbs "github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/marketdata"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/gateio"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type TestSuite struct {
	logger                 logster.Logger
	testServer             *httptest.Server
	testServerURL          string
	realServerURL          string
	upgrader               wbs.Upgrader
	websocket              *marketdata.GateIOMarketDataWebsocket
	wbsConn                *wbs.Conn
	ordBookCh              <-chan domain.GateIOOrderBook
	mu                     sync.Mutex
	instruments            []domain.GateIOInstrument
	instrumentsToSubscribe []string
	cfg                    config.MarketDataConfig
}

func NewTestSuite() *TestSuite { //nolint: gocognit // testing
	logger := logster.Discard
	testSuite := TestSuite{ //nolint: exhaustruct // testing
		logger: logger,
		upgrader: wbs.Upgrader{ //nolint: exhaustruct // testing
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		mu: sync.Mutex{},
	}
	marketDataHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			conn, err := testSuite.upgrader.Upgrade(w, r, nil)
			if err != nil {
				logger.Fatal("test server handler init error")
				return
			}
			testSuite.wbsConn = conn

			defer conn.Close()

			for {
				msgType, msg, err1 := conn.ReadMessage()
				if err1 != nil {
					logger.Fatal("test server handler init error")
					return
				}
				msgUnm := &gateio.MDRequest{}    //nolint:exhaustruct // omitting for test reason
				resp := &gateio.CommonResponse{} //nolint:exhaustruct // omitting for test reason
				err1 = easyjson.Unmarshal(msg, msgUnm)
				if err1 != nil {
					resp.Error.Code = -1
					resp.Error.Message = "wrong request format"
				}
				if msgUnm.Channel != gateio.SpotOrderBook {
					resp.Error.Code = -2
					resp.Error.Message = "wrong connection channel"
				}
				if len(msgUnm.Payload) < 3 {
					resp.Error.Code = -5
					resp.Error.Message = "wrong payload - at least 3 values"
				}
				switch msgUnm.Event {
				case gateio.EventTypeSubscribe:
					resp.Event = string(gateio.EventTypeSubscribe)
				case gateio.EventTypeUnsubscribe:
					resp.Event = string(gateio.EventTypeUnsubscribe)
				case gateio.EventTypeUpdate:
					resp.Event = string(gateio.EventTypeUpdate)
				default:
					resp.Error.Code = -4
					resp.Error.Message = "request event inconsistency"
				}

				if resp.Error.Message == "" &&
					msgUnm.Event != gateio.EventTypeUpdate {
					resp.Result = struct{ status string }{"success"}
				}

				if resp.Error.Message == "" &&
					msgUnm.Event == gateio.EventTypeUpdate {
					return
				}

				respMrsh, err1 := easyjson.Marshal(resp)
				if err1 != nil {
					logger.Fatal("message marshaling error: ", err)
					return
				}

				err1 = conn.WriteMessage(msgType, respMrsh)
				if err1 != nil {
					logger.Fatal("message writing error: ", err)
					return
				}

				if resp.Event == string(gateio.EventTypeSubscribe) {
					pair := msgUnm.Payload[0][0]

					go func() {
						bookUpdate := gateio.OrderBookUpdate{} //nolint:exhaustruct // omitting for test reason
						bookUpdate.Event = gateio.EventTypeUpdate
						bookUpdate.Result.S = pair
						bookUpdate.Time = time.Now().Unix()
						bookUpdate.Result.Asks = [][]string{
							{"73477.2", "0.32954"},
							{"73478", "0.01612"},
							{"73478.7", "0.01034"},
							{"73482.2", "0.0004"},
							{"73482.5", "0.16"},
						}
						bookUpdate.Result.Bids = [][]string{
							{"73466.6", "0.04352"},
							{"73467", "0.02933"},
							{"73475.2", "0.13609"},
							{"73476.4", "0.15785"},
							{"73477.1", "0.96608"},
						}
						msg, err = easyjson.Marshal(bookUpdate)
						if err != nil {
							logger.Fatal("message marshaling error: ", err)
							return
						}
						testSuite.mu.Lock()
						defer testSuite.mu.Unlock()
						err = conn.WriteMessage(msgType, msg)
						if err != nil {
							logger.Fatal("message writing error: ", err)
							return
						}
					}()
				}
			}
		})

	testSuite.testServer = httptest.NewServer(marketDataHandler)
	testSuite.realServerURL = "wss://api.gateio.ws/ws/v4/"
	testSuite.testServerURL = "ws" + strings.
		TrimPrefix(testSuite.testServer.URL, "http")

	testSuite.upgrader = wbs.Upgrader{ //nolint: exhaustruct // omitting for test reason
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	testSuite.instrumentsToSubscribe = []string{
		"BTC_USDT", "ETH_USDT"}

	testSuite.instruments = make(
		[]domain.GateIOInstrument,
		len(testSuite.instrumentsToSubscribe))

	for i, value := range testSuite.instrumentsToSubscribe {
		testSuite.instruments[i] = domain.GateIOSpotInstrument(value)
	}

	testSuite.cfg = config.MarketDataConfig{
		URL:                    testSuite.testServerURL,
		InstrumentsToSubscribe: testSuite.instrumentsToSubscribe,
	}

	testSuite.websocket = marketdata.
		NewGateIOMarketDataWebsocket(logger, testSuite.cfg)

	testSuite.ordBookCh = testSuite.websocket.GetOrderBookUpdate()

	return &testSuite
}
func TestGateIOMarketDataWebsocket_Init(t *testing.T) {
	testSuite := NewTestSuite()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "init", args: args{context.Background()}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testSuite.websocket.Init(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("GateIOMarketDataWebsocket.Init() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
	if testSuite.wbsConn == nil {
		t.Error("GateIOMarketDataWebsocket.Init() error = no connection logged")
	}
}
func TestGateIOMarketDataWebsocket_Run(t *testing.T) {
	testSuite := NewTestSuite()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "run", args: args{context.Background()}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, testSuite.websocket.Init(tt.args.ctx))

			if err := testSuite.websocket.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("GateIOMarketDataWebsocket.Run() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
	time.Sleep(50 * time.Microsecond)
	// time.Sleep(500 * time.Microsecond) // Necessary for real server testing
	timeout := time.After(10 * time.Millisecond)
	select {
	case <-testSuite.ordBookCh:
		return
	case <-timeout:
		t.Error("GateIOMarketDataWebsocket.Run() error = no websocket running logged")
	}
}
func TestGateIOMarketDataWebsocket_SubscribeInstruments(t *testing.T) {
	testSuite := NewTestSuite()
	type args struct {
		instruments []domain.GateIOInstrument
		ctx         context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "subscribe",
			args: args{
				instruments: testSuite.instruments,
				ctx:         context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		assert.NoError(t, testSuite.websocket.Init(tt.args.ctx))
		assert.NoError(t, testSuite.websocket.Run(tt.args.ctx))
		t.Run(tt.name, func(t *testing.T) {
			if err := testSuite.websocket.SubscribeInstruments(tt.args.instruments); (err != nil) != tt.wantErr {
				t.Errorf("GateIOMarketDataWebsocket.SubscribeInstruments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	time.Sleep(50 * time.Microsecond)
	// time.Sleep(500 * time.Microsecond) // Necessary for real server testing
	timeout := time.After(10 * time.Millisecond)
	select {
	case resp := <-testSuite.ordBookCh:
		if resp.Type != domain.UpdateType(gateio.EventTypeUpdate) {
			t.Error("GateIOMarketDataWebsocket.SubscribeInstruments() error = no update status in respond")
		}
		return
	case <-timeout:
		t.Error("GateIOMarketDataWebsocket.SubscribeInstruments() error = no respond after susbscribe request")
	}
}
func TestGateIOMarketDataWebsocket_GetOrderBookUpdate(t *testing.T) {
	testSuite := NewTestSuite()
	type args struct {
		instruments []domain.GateIOInstrument
		ctx         context.Context
	}
	tests := []struct {
		name string
		args args
		want <-chan domain.GateIOOrderBook
	}{
		{ //nolint:exhaustruct // omitting for test reason
			name: "get updates",
			args: args{
				instruments: testSuite.instruments,
				ctx:         context.Background(),
			},
		},
	}
	for _, tt := range tests {
		assert.NoError(t, testSuite.websocket.Init(tt.args.ctx))
		assert.NoError(t, testSuite.websocket.Run(tt.args.ctx))
		assert.NoError(t, testSuite.websocket.SubscribeInstruments(tt.args.instruments))
		t.Run(tt.name, func(t *testing.T) {
			got := testSuite.websocket.GetOrderBookUpdate()
			time.Sleep(50 * time.Microsecond)
			// time.Sleep(500 * time.Microsecond) // Necessary for real server testing
			for range testSuite.instruments {
				timeout := time.After(10 * time.Millisecond)
				select {
				case resp := <-got:
					if resp.Type != domain.UpdateType(gateio.EventTypeUpdate) {
						t.Error("GateIOMarketDataWebsocket.SubscribeInstruments() error = no update status in respond")
					}
					return
				case <-timeout:
					t.Error("GateIOMarketDataWebsocket.SubscribeInstruments() error = no respond after susbscribe request")
				}
			}
		})
	}
}
func TestGateIOMarketDataWebsocket_UnsubscribeInstruments(t *testing.T) {
	testSuite := NewTestSuite()
	type args struct {
		instruments []domain.GateIOInstrument
		ctx         context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unsubscribe",
			args: args{
				instruments: testSuite.instruments,
				ctx:         context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		assert.NoError(t, testSuite.websocket.Init(tt.args.ctx))
		assert.NoError(t, testSuite.websocket.Run(tt.args.ctx))
		assert.NoError(t, testSuite.websocket.SubscribeInstruments(tt.args.instruments))
		time.Sleep(50 * time.Microsecond)
		// time.Sleep(500 * time.Microsecond) // Necessary for real server testing
		for range testSuite.instruments {
			timeout := time.After(10 * time.Millisecond)
			select {
			case resp := <-testSuite.ordBookCh:
				if resp.Type != domain.UpdateType(gateio.EventTypeUpdate) {
					t.Error("GateIOMarketDataWebsocket.SubscribeInstruments() error = no updates status in respond")
				}
				return
			case <-timeout:
				t.Error("GateIOMarketDataWebsocket.SubscribeInstruments() error = no respond after susbscribe request")
			}
		}
		assert.NoError(t, testSuite.websocket.UnsubscribeInstruments(tt.args.instruments))
		t.Run(tt.name, func(t *testing.T) {
			if err := testSuite.websocket.UnsubscribeInstruments(tt.args.instruments); (err != nil) != tt.wantErr {
				t.Errorf("GateIOMarketDataWebsocket.UnsubscribeInstruments() error = %v, wantErr %v", err, tt.wantErr)
			}
			timeout := time.After(10 * time.Millisecond)
			select {
			case resp := <-testSuite.ordBookCh:
				t.Error("GateIOMarketDataWebsocket.UnsubscribeInstruments() error = updates got after unsubscribe: ", resp)
				return
			case <-timeout:
				return
			}
		})
	}
}
