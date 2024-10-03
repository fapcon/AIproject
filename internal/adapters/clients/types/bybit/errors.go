package bybit

import "errors"

var (
	ErrNoBalancesFound     = errors.New("no balances found")
	ErrNoOrderBookFound    = errors.New("no order book found")
	ErrInvalidParameters   = errors.New("invalid parameters")
	ErrNoOpenedOrdersFound = errors.New("no opened orders found")
)
