package okx

type MDRequest struct {
	Event EventType `json:"op"`
	Arg   []Arg     `json:"args"`
}

type Arg struct {
	Channel string `json:"channel"`
	InstID  string `json:"instId"`
}

type MDResponse struct {
	Event EventType `json:"event"`
	Code  string    `json:"code"`
	Msg   string    `json:"msg"`
	Arg   Arg       `json:"arg"`
}

type OrderBookUpdate struct {
	Arg    Arg       `json:"arg"`
	Action EventType `json:"action"`
	Data   []struct {
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
		TS        string     `json:"ts"`
		Checksum  int        `json:"checksum"`
		PrevSeqID int        `json:"prevSeqId"`
		SeqID     int        `json:"seqId"`
	} `json:"data"`
}

type EventType string

const (
	EventTypeLogin       EventType = "login"
	EventTypeSubscribe   EventType = "subscribe"
	EventTypeUnsubscribe EventType = "unsubscribe"
	EventTypeError       EventType = "error"
	EventTypeUpdate      EventType = "update"
	EventTypeSnapshot    EventType = "snapshot"
	EventTypeAmend       EventType = "amend-order"
	EventTypeCancel      EventType = "cancel-order"
	EventTypeOrder       EventType = "order"
)

const (
	OrdersChannel = "orders"
	BooksChannel  = "books"
)
