package gateio

// https://www.gate.io/docs/developers/apiv4/ws/en/#spot-websocket-v4

type CommonResponse struct {
	Time    int64  `json:"time"`
	TimeMS  int64  `json:"time_ms"`
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Error   struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result any `json:"result"` // Channel-specific
}

/* https://www.gate.io/docs/developers/apiv4/ws/en/#limited-level-full-order-book-snapshot

Channel - spot.order_book
Update speed: 1000ms or 100ms
This channel does not require authentication.
*/

type MDRequest struct {
	Time    int64     `json:"time"`
	Channel string    `json:"channel"`
	Event   EventType `json:"event"`
	// Payload is currency pair, level, interval as ["BTC_USDT", "5", "100ms"]
	Payload [][3]string `json:"payload"`
}

type EventType string

const (
	EventTypeSubscribe   EventType = "subscribe"
	EventTypeUnsubscribe EventType = "unsubscribe"
	EventTypeUpdate      EventType = "update"
)

type OrderBookUpdate struct {
	Time    int64           `json:"time"`
	TimeMS  int64           `json:"time_ms"`
	Channel string          `json:"channel"`
	Event   EventType       `json:"event"`
	Result  OrderBookResult `json:"result"`
}

type OrderBookResult struct {
	TimeMS       int64      `json:"t"`
	LastUpdateID int64      `json:"lastUpdateId"`
	S            string     `json:"s"`
	Bids         [][]string `json:"bids"`
	// Bids [][]string `json:"b"` // for updates testing
	Asks [][]string `json:"asks"`
	// Asks [][]string `json:"a"` // for updates testing
}

/*
https://www.gate.io/docs/developers/apiv4/ws/en/#order-book-channel

spot.book_ticker
Pushes any update about the price and amount of best bid or ask price in realtime for subscribed currency pairs.

spot.order_book_update
Periodically notify order book changed levels which can be used to locally manage an order book.

spot.order_book
Periodically notify top bids and asks snapshot with limited levels.

*/

type OrderBookChannel string

const (
	SpotOrderBook = "spot.order_book"
	// UpdateOrderBook = "spot.order_book_update" // for updates testing.
)

type OrderBookLevel string
type OrderBookInterval string

const (
	Level50       OrderBookLevel    = "50"
	Interval100ms OrderBookInterval = "100ms"
)
