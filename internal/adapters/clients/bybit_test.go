package clients_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/clients"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const bybitTestingURL = "https://api-testnet.bybit.com"

//nolint:gochecknoglobals //for testing purposes
var (
	instance *Suite
	once     sync.Once
)

type Suite struct {
	*testing.T
	api *clients.BybitAPI
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	once.Do(func() {
		t.Helper()
		cfg := config.HTTPClient{
			APIKey:        "6G538aBqXBV7lDj7ff",
			APISecret:     "mT3OGxi6oGZuvlfy9nHct3IgDzuL5zXAyy4J",
			APIPassphrase: "",
			URL:           bybitTestingURL,
			Timeout:       time.Second * 5,
		}
		client := clients.NewBybitAPI(logster.Discard, cfg)
		instance = &Suite{api: client, T: t}
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})
	return ctx, instance
}

func TestCreateOrderSpot_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	// проверим отсутствие активных ордеров
	symbol := domain.Symbol("BTCUSDT")
	openedOrders, err := suite.api.GetSpotOpenOrders(ctx, symbol)
	require.Error(t, err)
	assert.ErrorContainsf(t, err, "no opened orders found", err.Error())
	assert.Len(t, openedOrders, 0)
	// создадим ордер
	createdOrder, err := suite.api.CreateSpotOrder(ctx, domain.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(0.001),
		Price:         decimal.NewFromFloat(25000),
		ClientOrderID: "",
	})
	require.NoError(t, err)
	// проверим созданный ордер
	var expectedOrderResponseAck domain.BybitOrderResponseAck
	assert.IsType(t, expectedOrderResponseAck, createdOrder)
	// ордер должен появится на бирже
	openedOrders, err = suite.api.GetSpotOpenOrders(ctx, symbol)
	require.NoError(t, err)
	assert.Len(t, openedOrders, 1)

	// отменим конкретно его
	cancelResp, err := suite.api.CancelOrder(ctx, "BTCUSDT", createdOrder.OrderID, createdOrder.ClientOrderID)
	require.NoError(t, err)
	// проверим отмененный ордер
	assert.Equal(t, createdOrder.OrderID, cancelResp.OrderID)
	assert.Equal(t, createdOrder.ClientOrderID, cancelResp.ClientOrderID)
}

func TestAmendOrderSpot_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	request := domain.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(0.001),
		Price:         decimal.NewFromFloat(25000),
		ClientOrderID: "",
	}
	// проверим отсутствие активных ордеров
	openedOrders, err := suite.api.GetSpotOpenOrders(ctx, request.Symbol)
	require.Error(t, err)
	assert.ErrorContainsf(t, err, "no opened orders found", err.Error())
	assert.Len(t, openedOrders, 0)
	// создадим ордер
	createdOrder, err := suite.api.CreateSpotOrder(ctx, request)
	require.NoError(t, err)
	// проверим созданный ордер
	var expectedOrderResponseAck domain.BybitOrderResponseAck
	assert.IsType(t, expectedOrderResponseAck, createdOrder)
	// пытаемся изменить ордер
	deltaQty := decimal.NewFromFloat(0.001)
	deltaPrice := decimal.NewFromFloat(100)
	newPrice := decimal.Sum(request.Price, deltaPrice)
	newQty := decimal.Sum(request.Quantity, deltaQty)
	err = suite.api.AmendSpotOrder(
		ctx,
		"BTCUSDT",
		createdOrder.OrderID,
		createdOrder.ClientOrderID,
		newQty,
		newPrice,
	)

	require.NoError(t, err)
	// получаем открытые ордера
	openedOrders, err = suite.api.GetSpotOpenOrders(ctx, request.Symbol)
	require.NoError(t, err)
	// проверяем изменения
	assert.Equal(t, newPrice.String(), openedOrders[0].Price.String())
	assert.Equal(t, newQty.String(), openedOrders[0].OrigQty.String())
	// отменим все ордера для теста CancelFutureOrders
	err = suite.api.CancelSpotOrders(ctx, request.Symbol)
	require.NoError(t, err)
}

func TestCreateOrderFutures_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	// проверим отсутствие активных ордеров
	symbol := domain.Symbol("BTCUSDT")
	openedOrders, err := suite.api.GetFuturesOpenOrders(ctx, symbol)
	require.Error(t, err)
	assert.ErrorContainsf(t, err, "no opened orders found", err.Error())
	assert.Len(t, openedOrders, 0)
	// создадим ордер
	createdOrder, err := suite.api.CreateFuturesOrder(ctx, domain.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(0.001),
		Price:         decimal.NewFromFloat(50000),
		ClientOrderID: "",
	})
	require.NoError(t, err)
	// проверим созданный ордер
	var expectedOrderResponseAck domain.BybitOrderResponseAck
	assert.IsType(t, expectedOrderResponseAck, createdOrder)
	// ордер должен появится на бирже
	openedOrders, err = suite.api.GetFuturesOpenOrders(ctx, symbol)
	require.NoError(t, err)
	assert.Len(t, openedOrders, 1)

	// отменим конкретно его
	err = suite.api.CancelFuturesOrders(ctx, symbol)
	require.NoError(t, err)
}

func TestAmendOrderFutures_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	request := domain.OrderRequest{
		Symbol:        "BTCUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(0.001),
		Price:         decimal.NewFromFloat(50000),
		ClientOrderID: "",
	}
	// проверим отсутствие активных ордеров
	openedOrders, err := suite.api.GetFuturesOpenOrders(ctx, request.Symbol)
	require.Error(t, err)
	assert.ErrorContainsf(t, err, "no opened orders found", err.Error())
	assert.Len(t, openedOrders, 0)
	// создадим ордер
	createdOrder, err := suite.api.CreateFuturesOrder(ctx, request)
	require.NoError(t, err)
	// проверим созданный ордер
	var expectedOrderResponseAck domain.BybitOrderResponseAck
	assert.IsType(t, expectedOrderResponseAck, createdOrder)
	// пытаемся изменить ордер
	deltaQty := decimal.NewFromFloat(0.001)
	deltaPrice := decimal.NewFromFloat(100)
	newPrice := decimal.Sum(request.Price, deltaPrice)
	newQty := decimal.Sum(request.Quantity, deltaQty)
	err = suite.api.AmendFuturesOrder(
		ctx,
		"BTCUSDT",
		createdOrder.OrderID,
		createdOrder.ClientOrderID,
		newQty,
		newPrice,
	)

	require.NoError(t, err)
	// получаем открытые ордера
	openedOrders, err = suite.api.GetFuturesOpenOrders(ctx, request.Symbol)
	require.NoError(t, err)
	// проверяем изменения
	assert.Equal(t, newPrice.String(), openedOrders[0].Price.String())
	assert.Equal(t, newQty.String(), openedOrders[0].OrigQty.String())
	// отменим все ордера для теста CancelFutureOrders
	err = suite.api.CancelFuturesOrders(ctx, request.Symbol)
	require.NoError(t, err)
}

func TestGetOrderBook_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	response, err := suite.api.GetOrderBook(ctx, "BTCUSDT", 5)
	require.NoError(t, err)
	assert.Len(t, response.Bids, 5)
	assert.Len(t, response.Asks, 5)
}

func TestGetBalances_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	balances, err := suite.api.GetBalances(ctx)
	require.NoError(t, err)
	assert.Len(t, balances, 2)
}

func TestGetSpotOrderHistory_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	// проверяем наличие ордеров за последние 2 сек
	startTime := time.Now().Add(-2 * time.Second)
	orders, err := suite.api.GetSpotOrderHistory(ctx, startTime, time.Now())
	require.NoError(t, err)
	assert.Len(t, orders, 0)
	// размещаем ордер
	_, err = suite.api.CreateSpotOrder(ctx, domain.OrderRequest{
		Symbol:        "DOGEUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(1),
		Price:         decimal.NewFromFloat(0.05),
		ClientOrderID: "",
	})
	require.NoError(t, err)
	// отменяем ордер
	err = suite.api.CancelSpotOrders(ctx, "DOGEUSDT")
	require.NoError(t, err)
	// проверяем наличие истории ордеров за последние после добавления
	time.Sleep(1500 * time.Millisecond)
	orders, err = suite.api.GetSpotOrderHistory(ctx, startTime, time.Now())
	require.NoError(t, err)
	assert.Len(t, orders, 1)
}

func TestGetFuturesOrderHistory_Happy(t *testing.T) {
	ctx, suite := NewSuite(t)
	// проверяем наличие ордеров за последние 2 сек
	startTime := time.Now().Add(-2 * time.Second)
	orders, err := suite.api.GetFuturesOrderHistory(ctx, startTime, time.Now())
	require.NoError(t, err)
	assert.Len(t, orders, 0)
	// размещаем ордер
	_, err = suite.api.CreateFuturesOrder(ctx, domain.OrderRequest{
		Symbol:        "ETHUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(0.1),
		Price:         decimal.NewFromFloat(3500),
		ClientOrderID: "",
	})
	require.NoError(t, err)
	// отменяем ордер
	err = suite.api.CancelFuturesOrders(ctx, "ETHUSDT")
	require.NoError(t, err)
	// проверяем наличие истории ордеров за последние после добавления
	time.Sleep(1500 * time.Millisecond)
	orders, err = suite.api.GetFuturesOrderHistory(ctx, startTime, time.Now())
	require.NoError(t, err)
	assert.Len(t, orders, 1)
}

func TestCreateSpotOrder_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.CreateSpotOrder(ctx, domain.OrderRequest{
		Symbol:        "RUBTENGE",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(1),
		Price:         decimal.NewFromFloat(0.10),
		ClientOrderID: "",
	})
	require.Error(t, err)
}

func TestCreateFuturesOrder_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.CreateFuturesOrder(ctx, domain.OrderRequest{
		Symbol:        "TENGEUSDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   domain.TimeInForceGTC,
		Quantity:      decimal.NewFromFloat(1),
		Price:         decimal.NewFromFloat(0.10),
		ClientOrderID: "",
	})
	require.Error(t, err)
}

func TestAmendSpotOrder_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	err := suite.api.AmendSpotOrder(
		ctx,
		"TENGEUSDT",
		"123",
		"none",
		decimal.NewFromFloat(0.10),
		decimal.NewFromFloat(11000),
	)
	require.Error(t, err)
}

func TestAmendFuturesOrder_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	err := suite.api.AmendFuturesOrder(
		ctx,
		"TENGEUSDT",
		"123",
		"none",
		decimal.NewFromFloat(0.10),
		decimal.NewFromFloat(11000),
	)
	require.Error(t, err)
}

func TestCancelOrder_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.CancelOrder(ctx, "DOGEUSDT", "123", "none")
	require.Error(t, err)
}

func TestGetSpotOpenOrders_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.GetSpotOpenOrders(ctx, "")
	require.Error(t, err)
}

func TestGetFuturesOpenOrders_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.GetFuturesOpenOrders(ctx, "")
	require.Error(t, err)
}

func TestGetOrderBook_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.GetOrderBook(ctx, "", -1)
	require.Error(t, err)
}

func TestGetSpotOrderHistory_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.GetSpotOrderHistory(ctx, time.Now().Add(time.Hour), time.Now())
	require.Error(t, err)
}

func TestGetFuturesOrderHistory_Error(t *testing.T) {
	ctx, suite := NewSuite(t)
	_, err := suite.api.GetFuturesOrderHistory(ctx, time.Now().Add(time.Hour), time.Now())
	require.Error(t, err)
}
