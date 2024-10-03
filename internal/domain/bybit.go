package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type BybitOrderBook struct {
	Type      UpdateTypeBybit
	Symbol    BybitInstrument
	Timestamp time.Time
	Asks      map[string]float64
	Bids      map[string]float64
}

type BybitInstrument interface {
	Symbol() Symbol
	Exchange() Exchange
	isInstrument()
	isSpotInstrument() bool
	isFuturesInstrument() bool
}

type BybitSpotInstrument string

func (i BybitSpotInstrument) Symbol() Symbol            { return Symbol(i) }
func (i BybitSpotInstrument) Exchange() Exchange        { return BybitExchange }
func (i BybitSpotInstrument) isInstrument()             {}
func (i BybitSpotInstrument) isSpotInstrument() bool    { return true }
func (i BybitSpotInstrument) isFuturesInstrument() bool { return false }

type UpdateTypeBybit string

const (
	SnapshotBybit UpdateTypeBybit = "snapshot"
	UnknownBybit  UpdateTypeBybit = "unknown"
	DeltaBybit    UpdateTypeBybit = "delta"
)

type BybitFuturesInstrument string

func (i BybitFuturesInstrument) Symbol() Symbol            { return Symbol(i) }
func (i BybitFuturesInstrument) Exchange() Exchange        { return BybitExchange }
func (i BybitFuturesInstrument) isInstrument()             {}
func (i BybitFuturesInstrument) isSpotInstrument() bool    { return false }
func (i BybitFuturesInstrument) isFuturesInstrument() bool { return true }

type BybitUserStreamOrderUpdate struct {
	Instrument BybitInstrument
	OrderID    OrderID
	//TradeID           OrderID
	ClientOrderID ClientOrderID
	//RequestID         ClientOrderID // Needed to identify amend/move order
	//AmendResult       bool
	Price             decimal.Decimal
	OrigQty           decimal.Decimal
	ExecutedQty       decimal.Decimal
	FillSize          decimal.Decimal
	FillPrice         decimal.Decimal
	Status            OrderStatus
	TimeInForce       TimeInForce
	Type              OrderType
	Side              Side
	TransactionTime   time.Time
	OrderCreationTime time.Time
	//FillTime          time.Time
	FeeAmount   decimal.Decimal
	FeeCurrency Currency
}

type BybitOrderResponseAck struct {
	OrderID       OrderID
	ClientOrderID ClientOrderID
}
