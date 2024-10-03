package gateio

import (
	"strconv"
	"time"

	"studentgit.kata.academy/quant/torque/internal/domain"
)

func OrderBookToDomain(book OrderBookUpdate) domain.GateIOOrderBook {
	data := book.Result

	bids := make(map[string]float64, len(book.Result.Bids))
	asks := make(map[string]float64, len(book.Result.Asks))

	// Convert bids
	for _, bid := range data.Bids {
		price := bid[0]
		qty, _ := strconv.ParseFloat(bid[1], 64)
		bids[price] = qty
	}

	// Convert asks
	for _, ask := range data.Asks {
		price := ask[0]
		qty, _ := strconv.ParseFloat(ask[1], 64)
		asks[price] = qty
	}

	return domain.GateIOOrderBook{
		Type:      UpdateTypeToDomain(book.Event),
		Symbol:    domain.GateIOSpotInstrument(book.Result.S),
		Timestamp: time.UnixMilli(book.Result.TimeMS),
		Bids:      bids,
		Asks:      asks,
	}
}

func UpdateTypeToDomain(updateType EventType) domain.UpdateType {
	switch updateType { //nolint:exhaustive // we only care about order book updates here
	case "update":
		return domain.Update
	default:
		return domain.Unknown
	}
}
