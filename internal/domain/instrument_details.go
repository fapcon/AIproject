package domain

import "github.com/shopspring/decimal"

type InstrumentDetails struct {
	Exchange       Exchange
	LocalSymbol    Symbol
	ExchangeSymbol Symbol
	PricePrecision int
	SizePrecision  int
	MinLot         decimal.Decimal
}
