package domain

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Exchange string

type Symbol string

type Currency string

type Level struct {
	Price decimal.Decimal
	Qty   decimal.Decimal
}

type LocalInstrument struct {
	Exchange Exchange
	Symbol   Symbol
}

type ExchangeInstrument struct {
	Exchange Exchange
	Symbol   Symbol
}

const (
	OKXExchange     Exchange = "OKX"
	GateIOExchange  Exchange = "GateIO"
	BybitExchange   Exchange = "Bybit"
	UnknownExchange Exchange = "Unknown"
)

type Instrument struct {
	Exchange       Exchange
	LocalSymbol    Symbol
	ExchangeSymbol Symbol
}

type ExchangeAccount string
type OrderID string
type ClientOrderID string

var (
	ErrOrderRateLimitExceed       = errors.New("order_rate_limit_error")
	ErrPlaceOrderRateLimitExceed  = errors.New("okx place order rate limit exceed")
	ErrCancelOrderRateLimitExceed = errors.New("okx cancel order rate limit exceed")
	ErrAmendOrderRateLimitExceed  = errors.New("okx amend order rate limit exceed")
	ErrBalancesRateLimitExceed    = errors.New("okx balances rate limit exceed")
)

const (
	FakeOrder = "fakeOrder"
)

const (
	Precision = 0.000001
)
