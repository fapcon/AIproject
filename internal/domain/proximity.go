package domain

import (
	"time"
)

type PistonOrderID string
type PistonClientOrderRequestID string

type PistonID string
type PistonOrderRequestID string

type PistonIDKey struct {
	OrderRequestID PistonOrderRequestID
	PistonID       PistonID
}

func (k *PistonIDKey) IsEmpty() bool {
	return k.PistonID == "" && k.OrderRequestID == ""
}

type PistonSubscribeMarketDataMessage struct {
	Instruments []Symbol
}

type PistonBalancesRequestMessage struct{}
type PistonInstrumentsDetailsRequestMessage struct{}

type PistonOpenOrderRequestMessage struct {
	RequestID int64
}

type PistonOrderHistoryRequestMessage struct {
	RequestID int64
	StartTime time.Time
	EndTime   time.Time
}

type PistonBalancesMessage struct {
	Exchange        Exchange
	ExchangeAccount ExchangeAccount
	Balances        []Balance
	Timestamp       time.Time
}

type PistonOrderHistoryMessage struct {
	RequestID int64
	FromTime  time.Time
	ToTime    time.Time
	Orders    []PistonOrderHistory
}

type PistonOpenOrdersMessage struct {
	RequestID int64
	Orders    []PistonOpenOrder
}
