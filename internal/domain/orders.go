package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ClientOrderRequestID  PistonClientOrderRequestID // Needed for Piston
	PistonOrderID         PistonOrderID              // Needed for Piston
	ExchangeClientOrderID ClientOrderID              // Needed to cancel orders
	OrderID               OrderID                    // Needed for bookkeeping
	TradeID               OrderID
	Exchange              Exchange
	ExchangeAccount       ExchangeAccount
	Side                  Side
	Instrument            Symbol
	Type                  OrderType
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	RemainingSize         decimal.Decimal
	Filled                decimal.Decimal
	Volume                decimal.Decimal
	Fee                   decimal.Decimal
	FeeCurrency           Currency
	Created               time.Time
	UpdatedAt             time.Time
	Status                OrderStatus
}

type FixFill struct {
	ExchangeClientOrderID ClientOrderID
	OrderID               OrderID
	TradeID               OrderID
	Side                  Side
	Instrument            Symbol
	Price                 decimal.Decimal
	RemainingSize         decimal.Decimal
	FilledSize            decimal.Decimal
	Fee                   decimal.Decimal
	FeeCurrency           Currency
	Timestamp             time.Time
	Status                OrderStatus
}

// OpenOrder TODO Use this struct instead of ORDER to give Piston open orders
type OpenOrder struct {
	ClientOrderID ClientOrderID
	OrderID       OrderID
	Side          Side
	Instrument    Symbol
	Type          OrderType
	TimeInForce   TimeInForce
	OrigQty       decimal.Decimal
	Price         decimal.Decimal
	ExecutedQty   decimal.Decimal
	Created       time.Time
	UpdatedAt     time.Time
}

type OrderRequest struct {
	Symbol        Symbol
	Side          Side
	Type          OrderType
	TimeInForce   TimeInForce
	Quantity      decimal.Decimal
	Price         decimal.Decimal
	ClientOrderID ClientOrderID
}

type Side string

const (
	SideSell    Side = "SELL"
	SideBuy     Side = "BUY"
	SideUnknown Side = "UNKNOWN"
)

type OrderType string

const (
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeLimitMaker OrderType = "LIMIT_MAKER"
	OrderTypeLimitFOK   OrderType = "LIMIT_FOK"
	OrderTypeLimitIOC   OrderType = "LIMIT_IOC"
	OrderTypeUnknown    OrderType = "UNKNOWN"
)

type OrderStatus string

const (
	OrderStatusCanceled        OrderStatus = "Canceled"
	OrderStatusCreated         OrderStatus = "Created"
	OrderStatusPending         OrderStatus = "Pending"
	OrderStatusFilled          OrderStatus = "Filled"
	OrderStatusPartiallyFilled OrderStatus = "PartiallyFilled"
	OrderStatusRejected        OrderStatus = "Rejected"
	OrderStatusRejectCancel    OrderStatus = "RejectCancel"
	OrderStatusExpired         OrderStatus = "Expired"
	OrderStatusExecuted        OrderStatus = "Executed"
	OrderStatusUnknown         OrderStatus = "Unknown"
)

type TimeInForce string

const (
	TimeInForceGTC     TimeInForce = "GTC"
	TimeInForceIOC     TimeInForce = "IOC"
	TimeInForceFOK     TimeInForce = "FOK"
	TimeInForceUnknown TimeInForce = "UNKNOWN"
)

type AddOrder struct {
	OrderRequestID PistonOrderRequestID
	PistonID       PistonID
	Instrument     Symbol
	Side           Side
	OrderType      OrderType
	Size           decimal.Decimal
	Price          decimal.Decimal
	Created        time.Time
}

type OrderAdded struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	Exchange              Exchange
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
	Instrument            Symbol
	Side                  Side
	OrderType             OrderType
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	LastUpdated           time.Time
}

type FixOrderEvent struct {
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
}

type FixOrderRejectEvent struct {
	ClientOrderID   ClientOrderID
	ExchangeOrderID OrderID
	Error           string
}

type OrderAddReject struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	ExchangeClientOrderID ClientOrderID
}

type MoveOrder struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	Exchange              Exchange
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
	Instrument            Symbol
	Side                  Side
	OrderType             OrderType
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	Created               time.Time
}
type OrderMoved struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	Exchange              Exchange
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	Created               time.Time
}

type OrderMoveReject struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	ExchangeClientOrderID ClientOrderID
	ExchangeID            OrderID
}

type CancelOrder struct {
	Instrument            Symbol
	Side                  Side
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	ExchangeClientOrderID ClientOrderID
	ExchangeID            OrderID
}

type OrderCancelled struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	ExchangeClientOrderID ClientOrderID
	ExchangeID            OrderID
}
type OrderCancelReject struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	ExchangeClientOrderID ClientOrderID
	ExchangeID            OrderID
}

type OrderFilled struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	Exchange              Exchange
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
	TradeID               OrderID
	Instrument            Symbol
	Side                  Side
	OrderType             OrderType
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	Fee                   decimal.Decimal
	FeeCurrency           Currency
	Timestamp             time.Time
}

type OrderExecuted struct {
	OrderRequestID        PistonOrderRequestID
	PistonID              PistonID
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
}

type CancelAll struct{}

type PistonOpenOrder struct {
	PistonID              PistonID
	Exchange              Exchange
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
	Instrument            Symbol
	Side                  Side
	OrderType             OrderType
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	RemainingSize         decimal.Decimal
	Created               time.Time
	LastUpdated           time.Time
}

type ActiveCommand interface {
	GetID() ClientOrderID
}

type AddCmd struct {
	PistonIDKey     PistonIDKey
	ExchangeOrderID ClientOrderID
	Order           OrderAdded
}

func (cmd *AddCmd) GetID() ClientOrderID {
	return cmd.ExchangeOrderID
}

type MoveID struct {
	NewID    ClientOrderID
	CancelID ClientOrderID
}

type AddMoveCmd struct {
	MoveID
	PistonIDKey PistonIDKey
	OrderAdded  OrderAdded
	Processed   bool
}

func (cmd *AddMoveCmd) GetID() ClientOrderID {
	return cmd.NewID
}

type CancelMoveCmd struct {
	MoveID
	PistonIDKey PistonIDKey
	Processed   bool
}

func (cmd *CancelMoveCmd) GetID() ClientOrderID {
	return cmd.CancelID
}

type PistonOrderHistory struct {
	Exchange              Exchange
	ExchangeID            OrderID
	ExchangeClientOrderID ClientOrderID
	ExchangeAccount       ExchangeAccount
	TradeID               OrderID
	Instrument            Symbol
	Side                  Side
	Size                  decimal.Decimal
	Price                 decimal.Decimal
	Volume                decimal.Decimal
	Fee                   decimal.Decimal
	FeeCurrency           Currency
	Timestamp             time.Time
}

type NewOrderRequest struct {
	ClientOrderID   ClientOrderID
	LocalInstrument LocalInstrument
	Price           decimal.Decimal
	Qty             decimal.Decimal
	Side            Side
	Type            OrderType
	TimeInForce     TimeInForce
}
