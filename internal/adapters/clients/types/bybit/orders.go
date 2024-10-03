package bybit

//go:generate easyjson -all orders.go

type CreateOrderParams struct {
	Category    string      `json:"category"`
	Symbol      string      `json:"symbol"`
	Side        string      `json:"side"`
	OrderType   string      `json:"orderType"`
	Qty         string      `json:"qty"`
	Price       string      `json:"price"`
	TimeInForce TimeInForce `json:"timeInForce"`
}

type CreateOrderResponse struct {
	CommonResponse
	Result struct {
		OrderID     string `json:"orderId"`
		OrderLinkID string `json:"orderLinkId"`
	} `json:"result"`
}

type AmendOrderParams struct {
	Category    string `json:"category"`
	Symbol      string `json:"symbol"`
	OrderID     string `json:"orderId"`
	OrderLinkID string `json:"orderLinkId"`
	Qty         string `json:"qty"`
	Price       string `json:"price"`
}

type AmendOrderResponse struct {
	CommonResponse
	Result struct {
		OrderID     string `json:"orderId"`
		OrderLinkID string `json:"orderLinkId"`
	} `json:"result"`
}

type CancelOrderParams struct {
	Category    string `json:"category"`
	Symbol      string `json:"symbol"`
	OrderLinkID string `json:"orderLinkId"`
	OrderID     string `json:"orderId"`
}

type CancelOrderResponse struct {
	CommonResponse
	Result struct {
		OrderID     string `json:"orderId"`
		OrderLinkID string `json:"orderLinkId"`
	} `json:"result"`
}

type GetOpenOrdersResponse struct {
	CommonResponse
	Result struct {
		List []struct {
			OrderID     string `json:"orderId"`
			OrderLinkID string `json:"orderLinkId"`
			Qty         string `json:"qty"`
			Price       string `json:"price"`
			CumExecQty  string `json:"cumExecQty"`
			Side        string `json:"side"`
			Symbol      string `json:"symbol"`
			OrderType   string `json:"orderType"`
			TimeInForce string `json:"timeInForce"`
			CreatedTime string `json:"createdTime"`
			UpdatedTime string `json:"updatedTime"`
		} `json:"list"`
		NextPageCursor string `json:"nextPageCursor"`
		Category       string `json:"category"`
	} `json:"result"`
}

type CancelAllOrdersParams struct {
	Category InstrumentType `json:"category"`
	Symbol   string         `json:"symbol"`
}

type CancelAllOrdersResponse struct {
	CommonResponse
	Result struct {
		List []struct {
			OrderID     string `json:"orderId"`
			OrderLinkID string `json:"orderLinkId"`
		} `json:"list"`
		Success string `json:"success"`
	} `json:"result"`
}

type OrderHistoryResponse struct {
	CommonResponse
	Result struct {
		List []struct {
			OrderID      string `json:"orderId"`
			OrderLinkID  string `json:"orderLinkId"`
			BlockTradeID string `json:"blockTradeId"`
			Symbol       string `json:"symbol"`
			Price        string `json:"price"`
			Qty          string `json:"qty"`
			Side         string `json:"side"`
			PlaceType    string `json:"placeType"`
			OrderType    string `json:"orderType"`
			OrderStatus  string `json:"orderStatus"`
			AvgPrice     string `json:"avgPrice"`
			LeavesQty    string `json:"leavesQty"`
			CumExecQty   string `json:"cumExecQty"`
			CumExecFee   string `json:"cumExecFee"`
			CreatedTime  string `json:"createdTime"`
			UpdatedTime  string `json:"updatedTime"`
		} `json:"list"`
		NextPageCursor string `json:"nextPageCursor"`
		Category       string `json:"category"`
	} `json:"result"`
}

type InstrumentType string

const (
	SPOT    InstrumentType = "spot"
	LINEAR  InstrumentType = "linear"
	INVERSE InstrumentType = "inverse"
	OPTION  InstrumentType = "option"
)

func (i InstrumentType) GetString() string {
	return string(i)
}

type TimeInForce string

const (
	TimeInForceGTC     TimeInForce = "GTC"
	TimeInForceIOC     TimeInForce = "IOC"
	TimeInForceFOK     TimeInForce = "FOK"
	TimeInForceUnknown TimeInForce = "UNKNOWN"
)
