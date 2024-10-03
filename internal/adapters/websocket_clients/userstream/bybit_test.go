package userstream_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/bybit"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/userstream"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

func TestNewBybitUserStreamWebsocket(t *testing.T) {
	t.Parallel()
	type args struct {
		logger logster.Logger
		config config.UserStreamConfig
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
				config: config.UserStreamConfig{
					APIKey:        "someKey",
					APISecret:     "someSecret",
					APIPassphrase: "somePassphrase",
					URL:           "wss://stream.bybit.example/",
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := userstream.NewBybitUserStreamWebsocket(tt.args.logger, tt.args.config)
			if got == nil && !tt.wantNil || got != nil && tt.wantNil {
				t.Errorf("NewBybitUserStreamWebsocket() == %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

const (
	testAPIKey    = "testAPIKey"
	testAPISecret = "testSecretKey"
)

type testServerRequest struct {
	Op   bybit.OpType `json:"op"`
	Args []any        `json:"args"`
}

type testServerPongResponse struct {
	ReqID  string   `json:"req_id"`
	Op     string   `json:"op"`
	Args   []string `json:"args"`
	ConnID string   `json:"conn_id"`
}

func echo(w http.ResponseWriter, r *http.Request) {
	//nolint:exhaustruct // don't need to fill all fields
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	for {
		mt, message, _ := c.ReadMessage()
		var req testServerRequest
		_ = json.Unmarshal(message, &req)

		writeResponse(c, mt, req)
	}
}

//nolint:exhaustruct // don't need to fill all fields in structures
func writeResponse(conn *websocket.Conn, mt int, req testServerRequest) {
	//nolint:exhaustive // don't need to check all cases
	switch req.Op {
	case bybit.OpTypeAuth:
		if len(req.Args) != 3 {
			resp, _ := json.Marshal(bybit.Response{
				Success: false,
				RetMsg:  "params count error",
				Op:      "auth",
			})
			_ = conn.WriteMessage(mt, resp)
			return
		}

		key, ok1 := req.Args[0].(string)
		exp, ok2 := req.Args[1].(float64)
		hexSig, ok3 := req.Args[2].(string)
		if !ok1 || !ok2 || !ok3 {
			resp, _ := json.Marshal(bybit.Response{
				Success: false,
				RetMsg:  "invalid params",
				Op:      "auth",
			})
			_ = conn.WriteMessage(mt, resp)
			return
		}

		if key != testAPIKey {
			resp, _ := json.Marshal(bybit.Response{
				Success: false,
				RetMsg:  "not valid apiKey",
				Op:      "auth",
			})
			_ = conn.WriteMessage(mt, resp)
			return
		} else if exp <= float64(time.Now().UnixNano()) {
			resp, _ := json.Marshal(bybit.Response{
				Success: false,
				RetMsg:  "invalid expire time",
				Op:      "auth",
			})
			_ = conn.WriteMessage(mt, resp)
			return
		}

		_, err := hex.DecodeString(hexSig)
		if err != nil {
			resp, _ := json.Marshal(bybit.Response{
				Success: false,
				RetMsg:  "invalid signature",
				Op:      "auth",
			})
			_ = conn.WriteMessage(mt, resp)
			return
		}

		resp, _ := json.Marshal(bybit.Response{
			Success: true,
			RetMsg:  "",
			Op:      "auth",
		})
		_ = conn.WriteMessage(mt, resp)
	case bybit.OpTypeRequestHeartbeat:
		resp, _ := json.Marshal(testServerPongResponse{
			ReqID: "1",
			Op:    "pong",
			Args:  []string{},
		})
		_ = conn.WriteMessage(mt, resp)
	case bybit.OpTypeSubscribe:
		if len(req.Args) < 1 {
			resp, _ := json.Marshal(bybit.Response{
				Success: false,
				RetMsg:  "without args error",
				Op:      "subscribe",
			})
			_ = conn.WriteMessage(mt, resp)
			return
		}

		resp, _ := json.Marshal(bybit.Response{
			Success: true,
			Op:      "subscribe",
		})
		_ = conn.WriteMessage(mt, resp)

		resp, _ = json.Marshal(bybit.OrderUpdate{
			ID:           "1",
			CreationTime: time.Now().Unix(),
			Topic:        bybit.OrderTopic,
			Data: []bybit.Data{{
				Symbol:       "USDTBTC",
				OrderID:      "1",
				Side:         "Sell",
				OrderType:    "Market",
				Price:        "72.5",
				Qty:          "1",
				TimeInForce:  "IOC",
				OrderStatus:  "Filled",
				OrderLinkID:  "123",
				CumExecQty:   "1",
				CumExecValue: "75",
				CumExecFee:   "0.358635",
				CreatedTime:  "1672364262444",
				UpdatedTime:  "1672364262457",
				FeeCurrency:  "USDT",
			}},
		})
		_ = conn.WriteMessage(mt, resp)
	default:
		resp, _ := json.Marshal(bybit.Response{
			Success: false,
			RetMsg:  "unknown operation",
		})
		_ = conn.WriteMessage(mt, resp)
	}
}

func setupTestServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(echo))
	return server
}

//nolint:exhaustruct // don't need to fill all fields
func TestBybitUserStreamWebsocket_Run(t *testing.T) {
	t.Parallel()
	testServer := setupTestServer()
	defer testServer.Close()

	url := "ws" + strings.TrimPrefix(testServer.URL, "http")
	expOrderUpdatesResp := bybit.OrderToDomain(bybit.OrderUpdate{
		ID:           "1",
		CreationTime: time.Now().Unix(),
		Topic:        bybit.OrderTopic,
		Data: []bybit.Data{{
			Symbol:       "USDTBTC",
			OrderID:      "1",
			Side:         "Sell",
			OrderType:    "Market",
			Price:        "72.5",
			Qty:          "1",
			TimeInForce:  "IOC",
			OrderStatus:  "Filled",
			OrderLinkID:  "123",
			CumExecQty:   "1",
			CumExecValue: "75",
			CumExecFee:   "0.358635",
			CreatedTime:  "1672364262444",
			UpdatedTime:  "1672364262457",
			FeeCurrency:  "USDT",
		}},
	})

	ws := userstream.NewBybitUserStreamWebsocket(logster.Discard, config.UserStreamConfig{
		APIKey:    testAPIKey,
		APISecret: testAPISecret,
		URL:       url,
	})

	err := ws.Init(context.Background())
	assert.Nil(t, err)
	err = ws.Run(context.Background())
	assert.Nil(t, err)
	orUps := <-ws.GetOrderUpdates()
	assert.Equal(t, expOrderUpdatesResp[0], orUps)
}
