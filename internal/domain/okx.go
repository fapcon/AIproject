package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type OKXOrderBook struct {
	Type      UpdateType
	Symbol    OKXInstrument
	Timestamp time.Time
	Asks      map[string]float64
	Bids      map[string]float64
	Checksum  int
}

type UpdateType string

const (
	Snapshot UpdateType = "snapshot"
	Update   UpdateType = "update"
	Unknown  UpdateType = "unknown"
)

type OKXOrderResponseAck struct {
	ClientOrderID ClientOrderID
	OrderID       OrderID
}

type OKXInstrument interface {
	Symbol() Symbol
	Exchange() Exchange
	isInstrument()
	isSpotInstrument() bool
	isFuturesInstrument() bool
}

type OKXSpotInstrument string // instrument in format BTC-USDT, ETH-USDT etc

func (i OKXSpotInstrument) Symbol() Symbol            { return Symbol(i) }
func (i OKXSpotInstrument) Exchange() Exchange        { return OKXExchange }
func (i OKXSpotInstrument) isInstrument()             {}
func (i OKXSpotInstrument) isSpotInstrument() bool    { return true }
func (i OKXSpotInstrument) isFuturesInstrument() bool { return false }

type OKXFuturesInstrument string // instrument in format BTC-USDT-SWAP, ETH-USDT-SWAP etc

func (i OKXFuturesInstrument) Symbol() Symbol            { return Symbol(i) }
func (i OKXFuturesInstrument) Exchange() Exchange        { return OKXExchange }
func (i OKXFuturesInstrument) isInstrument()             {}
func (i OKXFuturesInstrument) isSpotInstrument() bool    { return false }
func (i OKXFuturesInstrument) isFuturesInstrument() bool { return true }

type OKXUserStreamOrderUpdate struct {
	Instrument        OKXInstrument
	OrderID           OrderID
	TradeID           OrderID
	ClientOrderID     ClientOrderID
	RequestID         ClientOrderID // Needed to identify amend/move order
	AmendResult       bool
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
	FillTime          time.Time
	FeeAmount         decimal.Decimal
	FeeCurrency       Currency
}
