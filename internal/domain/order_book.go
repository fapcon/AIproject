package domain

import (
	"sort"
	"strconv"
)

type OrderBook struct {
	Exchange   Exchange
	Instrument Symbol
	Asks       map[string]float64
	Bids       map[string]float64
	Checksum   int
}

// UpdatePriceLevels updates multiple price levels in the order book.
func (ob *OrderBook) UpdatePriceLevels(isBid bool, levels map[float64]float64) {
	targetMap := ob.Asks
	if isBid {
		targetMap = ob.Bids
	}

	for price, qty := range levels {
		pricePrice := strconv.FormatFloat(price, 'f', -1, 64)

		if qty == 0 {
			delete(targetMap, pricePrice)
		} else {
			targetMap[pricePrice] = qty
		}
	}
}

// UpdatePriceLevelsString updates multiple price levels in the order book.
func (ob *OrderBook) UpdatePriceLevelsString(isBid bool, levels map[string]float64) {
	targetMap := ob.Asks
	if isBid {
		targetMap = ob.Bids
	}
	for price, qty := range levels {
		if qty == 0 {
			delete(targetMap, price)
		} else {
			targetMap[price] = qty
		}
	}
}

func getSortedLevels(levels map[string]float64, isBid bool) []PriceLevel {
	var sortedLevels []PriceLevel
	for k, v := range levels {
		sortedLevels = append(sortedLevels, PriceLevel{Price: k, Quantity: v})
	}

	// Sorting logic: For bids, in descending order; for asks, in ascending order
	sort.Slice(sortedLevels, func(i, j int) bool {
		priceI, _ := strconv.ParseFloat(sortedLevels[i].Price, 64)
		priceJ, _ := strconv.ParseFloat(sortedLevels[j].Price, 64)

		if isBid {
			return priceI > priceJ
		}
		return priceI < priceJ
	})

	return sortedLevels
}

type PriceLevel struct {
	Price    string
	Quantity float64
}

// TopLevels returns the top X levels of bids or asks.
func (ob *OrderBook) TopLevels(isBid bool, x int) []PriceLevel {
	var levels map[string]float64
	if isBid {
		levels = ob.Bids
	} else {
		levels = ob.Asks
	}

	sortedLevels := getSortedLevels(levels, isBid)

	// If x is greater than the number of levels, return all levels
	if x > len(sortedLevels) {
		x = len(sortedLevels)
	}

	return sortedLevels[:x]
}
