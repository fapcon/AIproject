package domain

import (
	"time"
)

type GateIOOrderBook struct {
	Type      UpdateType
	Symbol    GateIOInstrument
	Timestamp time.Time
	Asks      map[string]float64
	Bids      map[string]float64
}

type GateIOInstrument interface {
	Symbol() Symbol
	Exchange() Exchange
	isInstrument()
	isSpotInstrument() bool
	isFuturesInstrument() bool
}

type GateIOSpotInstrument string // instrument in format BTC_USDT, ETH_USDT etc

func (i GateIOSpotInstrument) Symbol() Symbol            { return Symbol(i) }
func (i GateIOSpotInstrument) Exchange() Exchange        { return GateIOExchange }
func (i GateIOSpotInstrument) isInstrument()             {}
func (i GateIOSpotInstrument) isSpotInstrument() bool    { return true }
func (i GateIOSpotInstrument) isFuturesInstrument() bool { return false }
