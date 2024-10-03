package bybit

//go:generate easyjson -all orderbook.go
type OrderBookResponse struct {
	CommonResponse
	Result struct {
		S   string     `json:"s"`   // Symbol name
		B   [][]string `json:"b"`   // Bid, buyer. Sort by price desc, > b[0] - bid price, b[1] - size
		A   [][]string `json:"a"`   // Ask, seller. Order by price asc, > a[0] - ask price, a[1] - size
		TS  int64      `json:"ts"`  // Timestamp
		U   int        `json:"u"`   // Update ID
		Seq int64      `json:"seq"` // Cross sequence
	} `json:"result"`
}

type OrderBookRequest struct {
	Category InstrumentType `url:"category"`
	Symbol   string         `url:"symbol"`
	Limit    int            `url:"limit,omitempty"`
}
