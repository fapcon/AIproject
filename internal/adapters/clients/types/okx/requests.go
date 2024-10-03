package okx

type RequestOrderBook struct {
	Instrument string `url:"instId"`
	Limit      int    `url:"sz"`
}

type RequestOpenOrders struct {
	InstrumentType string `url:"instType"`
}

type RequestGetOrder struct {
	Instrument    string `url:"instId"`
	OrderID       string `url:"ordId,omitempty"`
	ClientOrderID string `url:"clOrdId,omitempty"`
}

type RequestCancelOrder struct {
	Instrument    string `json:"instId"`
	OrderID       string `json:"ordId,omitempty"`
	ClientOrderID string `json:"clOrdId,omitempty"`
}

type RequestAmendOrder struct {
	Instrument    string `json:"instId"`
	OrderID       string `json:"ordId,omitempty"`
	ClientOrderID string `json:"clOrdId,omitempty"`
	NewPrice      string `json:"newPx"`
	NewSize       string `json:"newSz"`
	RequestID     string `json:"reqId"`
}

type RequestOrderHistory struct {
	InstrumentType string `url:"instType"`
	StartTime      string `url:"begin"`
	EndTime        string `url:"end"`
	Limit          string `url:"limit"`
	After          string `url:"after,omitempty"`
}

type RequestFillsHistory struct {
	InstrumentType string `url:"instType"`
	StartTime      string `url:"begin"`
	EndTime        string `url:"end"`
	Limit          string `url:"limit"`
	After          string `url:"after,omitempty"`
}

type RequestCreateOrder struct {
	Instrument     string `json:"instId"`
	ClientOrderID  string `json:"clOrdId"`
	Side           string `json:"side"`
	Type           string `json:"ordType"`
	Quantity       string `json:"sz"`
	Price          string `json:"px,omitempty"`
	TdMode         string `json:"tdMode"`
	TargetCurrency string `json:"tgtCcy,omitempty"`
}

type RequestCancelAllAfter struct {
	Timeout string `json:"timeOut"`
}
