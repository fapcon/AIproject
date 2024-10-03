package bybit

//go:generate easyjson -all orders.go

type AuthMessage struct {
	Op   OpType `json:"op"`
	Args []any  `json:"args"`
}

type Topic string

const (
	OrderTopic Topic = "order"
)

type SubscriptionMessage struct {
	Op   OpType  `json:"op"`
	Args []Topic `json:"args"`
}

type HeartbeatMessage struct {
	Op string `json:"op"`
}

type Response struct {
	Success bool   `json:"success"`
	RetMsg  string `json:"ret_msg"`
	Op      OpType `json:"op"`
	ConnID  string `json:"conn_id"`
}

type OrderUpdate struct {
	Response
	ID           string `json:"id"`
	Topic        Topic  `json:"topic"`
	CreationTime int64  `json:"creationTime"`
	Data         []Data `json:"data"`
}

type Data struct {
	Symbol             string `json:"symbol"`
	OrderID            string `json:"orderId"`
	Side               string `json:"side"`
	OrderType          string `json:"orderType"`
	CancelType         string `json:"cancelType"`
	Price              string `json:"price"`
	Qty                string `json:"qty"`
	OrderIv            string `json:"orderIv"`
	TimeInForce        string `json:"timeInForce"`
	OrderStatus        string `json:"orderStatus"`
	OrderLinkID        string `json:"orderLinkId"`
	LastPriceOnCreated string `json:"lastPriceOnCreated"`
	ReduceOnly         bool   `json:"reduceOnly"`
	LeavesQty          string `json:"leavesQty"`
	LeavesValue        string `json:"leavesValue"`
	CumExecQty         string `json:"cumExecQty"`
	CumExecValue       string `json:"cumExecValue"`
	AvgPrice           string `json:"avgPrice"`
	BlockTradeID       string `json:"blockTradeId"`
	PositionIdx        int64  `json:"positionIdx"`
	CumExecFee         string `json:"cumExecFee"`
	CreatedTime        string `json:"createdTime"`
	UpdatedTime        string `json:"updatedTime"`
	RejectReason       string `json:"rejectReason"`
	StopOrderType      string `json:"stopOrderType"`
	TpslMode           string `json:"tpslMode"`
	TriggerPrice       string `json:"triggerPrice"`
	TakeProfit         string `json:"takeProfit"`
	StopLoss           string `json:"stopLoss"`
	TpTriggerBy        string `json:"tpTriggerBy"`
	SlTriggerBy        string `json:"slTriggerBy"`
	TpLimitPrice       string `json:"tpLimitPrice"`
	SlLimitPrice       string `json:"slLimitPrice"`
	TriggerDirection   int64  `json:"triggerDirection"`
	TriggerBy          string `json:"triggerBy"`
	CloseOnTrigger     bool   `json:"closeOnTrigger"`
	Category           string `json:"category"`
	PlaceType          string `json:"placeType"`
	SMPType            string `json:"smpType"`
	SMPGroup           int64  `json:"smpGroup"`
	SMPOrderID         string `json:"smpOrderId"`
	FeeCurrency        string `json:"feeCurrency"`
}
