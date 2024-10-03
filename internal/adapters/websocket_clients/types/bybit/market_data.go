package bybit

type MDRequest struct {
	Op   OpType   `json:"op"`
	Args []string `json:"args"`
}

type MDResponse struct {
	Success bool   `json:"success"`
	RetMsg  string `json:"ret_msg"`
	ConnID  string `json:"conn_id"`
	Op      OpType `json:"op"`
}

type OrderBookUpdate struct {
	Topic  string `json:"topic"`
	Type   OpType `json:"type"`
	TS     int64  `json:"ts"`
	DataMD DataMD `json:"data"`
	CTS    int64  `json:"cts"`
}

type DataMD struct {
	Symbol string     `json:"s"`
	Bids   [][]string `json:"b"`
	Asks   [][]string `json:"a"`
	Update int64      `json:"u"`
	Seq    int64      `json:"seq"`
}

type OpType string

const (
	OpTypeAuth              OpType = "auth"
	OpTypeSubscribe         OpType = "subscribe"
	OpTypeUnsubscribe       OpType = "unsubscribe"
	OpTypeError             OpType = "error"
	OpTypeDelta             OpType = "delta"
	OpTypeSnapshot          OpType = "snapshot"
	OpTypeRequestHeartbeat  OpType = "ping"
	OpTypeResponseHeartbeat OpType = "pong"
)

type OrderBookLevel string

const (
	// push frequency: 10ms.
	Level1 OrderBookLevel = "1"
	// push frequency: 20ms.
	Level50 OrderBookLevel = "50"
)
