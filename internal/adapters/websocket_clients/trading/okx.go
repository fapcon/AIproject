package trading

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/mailru/easyjson"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	okxLimiter "studentgit.kata.academy/quant/torque/internal/adapters/clients"
	okxTypes "studentgit.kata.academy/quant/torque/internal/adapters/clients/types/okx"
	"studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients/types/okx"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/internal/domain/observation"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/metrics"
	websocketclient "studentgit.kata.academy/quant/torque/pkg/websocket_client"
)

type OKXTradingWebsocket struct {
	logger            logster.Logger
	client            *websocketclient.Client
	apiKey            string
	apiSecret         string
	apiPassphrase     string
	rateLimiterMap    *okxLimiter.RateLimiterMap
	ulidGenerator     func() string
	createOrderAck    chan domain.OKXOrderResponseAck
	modifyOrderAck    chan domain.OKXOrderResponseAck
	cancelOrderAck    chan domain.OKXOrderResponseAck
	modifyOrderFailed chan domain.OKXOrderResponseAck
	cancelOrderFailed chan domain.OKXOrderResponseAck
	createOrderFailed chan domain.OKXOrderResponseAck

	collector *metrics.Collector
}

func NewOkxTradingWebsocket(
	logger logster.Logger, config config.UserStreamConfig,
) *OKXTradingWebsocket {
	return &OKXTradingWebsocket{
		logger:            logger.WithField(f.Module, "okx_trading_websocket"),
		client:            websocketclient.NewClient(config.URL, logger),
		apiKey:            config.APIKey,
		apiSecret:         config.APISecret,
		apiPassphrase:     config.APIPassphrase,
		ulidGenerator:     func() string { return ulid.Make().String() },
		rateLimiterMap:    okxLimiter.NewRateLimiterMap(),
		createOrderAck:    make(chan domain.OKXOrderResponseAck),
		modifyOrderAck:    make(chan domain.OKXOrderResponseAck),
		cancelOrderAck:    make(chan domain.OKXOrderResponseAck),
		modifyOrderFailed: make(chan domain.OKXOrderResponseAck),
		cancelOrderFailed: make(chan domain.OKXOrderResponseAck),
		createOrderFailed: make(chan domain.OKXOrderResponseAck),
		collector:         metrics.NewNoopCollector(),
	}
}

func (c *OKXTradingWebsocket) IsConnected() bool {
	return c.client.IsConnected()
}

func (c *OKXTradingWebsocket) WithMetrics(collector *metrics.Collector) {
	c.collector = collector
}

func (c *OKXTradingWebsocket) Run(ctx context.Context) error {
	go c.dispatchRequests(ctx)
	c.logger.Infof("okx trading websocket is running")
	c.SetReconnectMessage()
	return nil
}

func (c *OKXTradingWebsocket) SetReconnectMessage() {
	c.client.SendMessageAfterReconnect([][]byte{c.getAuth()})
}

func (c *OKXTradingWebsocket) Init(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return err
	}
	err = c.SendAuth(c.getAuth())
	if err != nil {
		return err
	}
	return nil
}

func (c *OKXTradingWebsocket) Close() {
	c.logger.Infof("closing user stream websocket...")
	close(c.createOrderAck)
	close(c.modifyOrderAck)
	close(c.cancelOrderAck)
	close(c.cancelOrderFailed)
	close(c.modifyOrderFailed)
	close(c.createOrderFailed)
	c.client.Shutdown()
	c.logger.Infof("user stream websocket closed")
}

func (c *OKXTradingWebsocket) GetCreateOrderAck() <-chan domain.OKXOrderResponseAck {
	return c.createOrderAck
}

func (c *OKXTradingWebsocket) GetModifyOrderAck() <-chan domain.OKXOrderResponseAck {
	return c.modifyOrderAck
}

func (c *OKXTradingWebsocket) GetCancelOrderAck() <-chan domain.OKXOrderResponseAck {
	return c.cancelOrderAck
}

func (c *OKXTradingWebsocket) GetCancelOrderFailed() <-chan domain.OKXOrderResponseAck {
	return c.cancelOrderFailed
}

func (c *OKXTradingWebsocket) GetModifyOrderFailed() <-chan domain.OKXOrderResponseAck {
	return c.modifyOrderFailed
}

func (c *OKXTradingWebsocket) GetCreateOrderFailed() <-chan domain.OKXOrderResponseAck {
	return c.createOrderFailed
}

func (c *OKXTradingWebsocket) dispatchRequests(ctx context.Context) {
	for {
		select {
		case msg, ok := <-c.client.GetMessages():
			if !ok {
				c.logger.Infof("message channel closed, stopping dispatching requests")
				return
			}
			timeStart := time.Now()
			c.decodeMessage(msg)
			c.collector.GetHistogram(
				observation.WebsocketEventDecodeDuration,
			).Observe(float64(time.Since(timeStart).Microseconds()))

		case <-ctx.Done():
			c.logger.Infof("stopping dispatching requests")
			c.Close()
			return
		}
	}
}

func (c *OKXTradingWebsocket) getAuth() []byte {
	timestamp := time.Now().Unix()
	data := okx.LoginMessage{
		Op: okx.EventTypeLogin,
		Args: []okx.LoginArgs{
			{
				APIKey:     c.apiKey,
				Passphrase: c.apiPassphrase,
				Timestamp:  timestamp,
				Sign:       c.generateSignature(fmt.Sprintf("%dGET/users/self/verify", timestamp)),
			},
		},
	}
	jsonData, _ := easyjson.Marshal(data)
	return jsonData
}

func (c *OKXTradingWebsocket) SendAuth(msg []byte) error {
	return c.client.Send(msg)
}

func (c *OKXTradingWebsocket) CancelOrder(
	_ context.Context, symbol domain.OKXInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) error {
	if !c.rateLimiterMap.Get(okxTypes.CancelOrder).Allow() {
		return domain.ErrCancelOrderRateLimitExceed
	}
	data := okx.CancelOrder{
		ID: c.ulidGenerator(),
		Op: okx.EventTypeCancel,
		Args: []okx.CancelArgs{
			{
				InstID:  string(symbol.Symbol()),
				OrdID:   string(orderID),
				ClOrdID: string(clientOrderID),
			},
		},
	}
	msg, _ := easyjson.Marshal(data)
	return c.client.Send(msg)
}

func (c *OKXTradingWebsocket) CreateOrder(_ context.Context, request domain.OrderRequest) error {
	if !c.rateLimiterMap.Get(okxTypes.Order).Allow() {
		return domain.ErrOrderRateLimitExceed
	}
	//nolint:exhaustruct // TargetCurrency is filled bellow
	createOrderArgs := okx.CreateOrderArgs{
		Instrument:    string(request.Symbol),
		ClientOrderID: string(request.ClientOrderID),
		Side:          okx.SideFromDomain(request.Side),
		Type:          okx.OrderTypeFromDomain(request.Type),
		Quantity:      request.Quantity.String(),
		Price:         request.Price.String(),
		TdMode:        "cash",
		// TIF is a part of Type in OKX
	}
	if request.Type == domain.OrderTypeMarket {
		createOrderArgs.TargetCurrency = "base_ccy"
	}
	data := okx.CreateOrder{
		ID: c.ulidGenerator(),
		Op: okx.EventTypeOrder,
		Args: []okx.CreateOrderArgs{
			createOrderArgs,
		},
	}
	msg, _ := easyjson.Marshal(data)
	return c.client.Send(msg)
}

func (c *OKXTradingWebsocket) ModifyOrder(
	_ context.Context,
	symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	if !c.rateLimiterMap.Get(okxTypes.AmendOrder).Allow() {
		return domain.ErrAmendOrderRateLimitExceed
	}
	data := okx.ModifyOrder{
		ID: c.ulidGenerator(),
		Op: okx.EventTypeAmend,
		Args: []okx.ModifyArgs{
			{
				Instrument:    string(symbol.Symbol()),
				OrderID:       string(orderID),
				ClientOrderID: string(clientOrderID),
				NewPrice:      newPrice.String(),
				NewSize:       newSize.String(),
				RequestID:     string(clientOrderID),
			},
		},
	}
	msg, _ := easyjson.Marshal(data)
	return c.client.Send(msg)
}

func (c *OKXTradingWebsocket) generateSignature(data string) string {
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *OKXTradingWebsocket) decodeMessage(msg []byte) {
	var update okx.OrderAck

	c.logger.Debugf("received trading message %s", msg)

	err := easyjson.Unmarshal(msg, &update)
	if err != nil {
		c.logger.WithError(err).Errorf("error while unmarshalling message")
	}

	if update.Event == okx.EventTypeLogin {
		c.logger.WithField(f.OKXLoginResponse, string(msg)).Infof("login successful")
		return
	}

	ack := okx.OrderAckToDomain(update)
	if update.Code != "0" {
		c.logger.WithField(f.Message, string(msg)).
			WithField(f.OKXError, update.Msg).Errorf("received error response from okx")
		switch update.Op {
		case okx.EventTypeOrder:
			c.createOrderFailed <- ack
		case okx.EventTypeAmend:
			c.modifyOrderFailed <- ack
		case okx.EventTypeCancel:
			c.cancelOrderFailed <- ack
		case okx.EventTypeLogin,
			okx.EventTypeSubscribe,
			okx.EventTypeUnsubscribe,
			okx.EventTypeError,
			okx.EventTypeUpdate,
			okx.EventTypeSnapshot:
		default:
			c.logger.Errorf("unknown request type %s", update)
		}
		return
	}
	switch update.Op {
	case okx.EventTypeCancel:
		c.cancelOrderAck <- ack
	case okx.EventTypeOrder:
		c.createOrderAck <- ack
	case okx.EventTypeAmend:
		c.modifyOrderAck <- ack
	case okx.EventTypeLogin,
		okx.EventTypeSubscribe,
		okx.EventTypeUnsubscribe,
		okx.EventTypeError,
		okx.EventTypeUpdate,
		okx.EventTypeSnapshot:
	default:
		c.logger.Errorf("unknown request type %s", update.Op)
	}
}
