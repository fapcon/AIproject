package okx

type OrderUpdate struct {
	SubscribeResponse
	Arg struct {
		Channel  string `json:"channel"`
		InstType string `json:"instType"`
		UID      string `json:"uid"`
	} `json:"arg"`
	Data []struct {
		AccFillSz         string `json:"accFillSz"`
		AlgoClOrdID       string `json:"algoClOrdId"`
		AlgoID            string `json:"algoId"`
		AmendResult       string `json:"amendResult"`
		AmendSource       string `json:"amendSource"`
		AttachAlgoClOrdID string `json:"attachAlgoClOrdId"`
		AvgPx             string `json:"avgPx"`
		CTime             string `json:"cTime"`
		CancelSource      string `json:"cancelSource"`
		Category          string `json:"category"`
		Ccy               string `json:"ccy"`
		ClOrdID           string `json:"clOrdId"`
		Code              string `json:"code"`
		ExecType          string `json:"execType"`
		Fee               string `json:"fee"`
		FeeCcy            string `json:"feeCcy"`
		FillFee           string `json:"fillFee"`
		FillFeeCcy        string `json:"fillFeeCcy"`
		FillFwdPx         string `json:"fillFwdPx"`
		FillMarkPx        string `json:"fillMarkPx"`
		FillMarkVol       string `json:"fillMarkVol"`
		FillNotionalUsd   string `json:"fillNotionalUsd"`
		FillPnl           string `json:"fillPnl"`
		FillPx            string `json:"fillPx"`
		FillPxUsd         string `json:"fillPxUsd"`
		FillPxVol         string `json:"fillPxVol"`
		FillSz            string `json:"fillSz"`
		FillTime          string `json:"fillTime"`
		InstID            string `json:"instId"`
		InstType          string `json:"instType"`
		Lever             string `json:"lever"`
		Msg               string `json:"msg"`
		NotionalUsd       string `json:"notionalUsd"`
		OrdID             string `json:"ordId"`
		OrdType           string `json:"ordType"`
		Pnl               string `json:"pnl"`
		PosSide           string `json:"posSide"`
		Px                string `json:"px"`
		PxType            string `json:"pxType"`
		PxUsd             string `json:"pxUsd"`
		PxVol             string `json:"pxVol"`
		QuickMgnType      string `json:"quickMgnType"`
		Rebate            string `json:"rebate"`
		RebateCcy         string `json:"rebateCcy"`
		ReduceOnly        string `json:"reduceOnly"`
		ReqID             string `json:"reqId"`
		Side              string `json:"side"`
		SlOrdPx           string `json:"slOrdPx"`
		SlTriggerPx       string `json:"slTriggerPx"`
		SlTriggerPxType   string `json:"slTriggerPxType"`
		Source            string `json:"source"`
		State             string `json:"state"`
		StpID             string `json:"stpId"`
		StpMode           string `json:"stpMode"`
		Sz                string `json:"sz"`
		Tag               string `json:"tag"`
		TdMode            string `json:"tdMode"`
		TgtCcy            string `json:"tgtCcy"`
		TpOrdPx           string `json:"tpOrdPx"`
		TpTriggerPx       string `json:"tpTriggerPx"`
		TpTriggerPxType   string `json:"tpTriggerPxType"`
		TradeID           string `json:"tradeId"`
		UTime             string `json:"uTime"`
	} `json:"data"`
}

type InstrumentType string

const (
	SPOT    InstrumentType = "SPOT"
	MARGIN  InstrumentType = "MARGIN"
	SWAP    InstrumentType = "SWAP"
	FUTURES InstrumentType = "FUTURES"
	OPTION  InstrumentType = "OPTION"
)

type SubscribeResponse struct {
	Event EventType `json:"event"`
	Code  string    `json:"code"`
	Msg   string    `json:"msg"`
}

type SubscriptionArgs struct {
	Channel  string         `json:"channel"`
	InstType InstrumentType `json:"instType"`
}

type SubscriptionMessage struct {
	Op   EventType          `json:"op"`
	Args []SubscriptionArgs `json:"args"`
}

type LoginArgs struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  int64  `json:"timestamp"`
	Sign       string `json:"sign"`
}

type LoginMessage struct {
	Op   EventType   `json:"op"`
	Args []LoginArgs `json:"args"`
}

type OrderAck struct {
	ID      string    `json:"id"`
	Op      EventType `json:"op"`
	Event   EventType `json:"event"` // needed only for login message and should be ignored for all other messages
	Data    []Data    `json:"data"`
	Code    string    `json:"code"`
	Msg     string    `json:"msg"`
	InTime  string    `json:"inTime"`
	OutTime string    `json:"outTime"`
}

type Data struct {
	ClOrdID string `json:"clOrdId"`
	OrdID   string `json:"ordId"`
	Tag     string `json:"tag"`
	ReqID   string `json:"reqId"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

type CancelOrder struct {
	ID   string       `json:"id"`
	Op   EventType    `json:"op"`
	Args []CancelArgs `json:"args"`
}

type CancelArgs struct {
	InstID  string `json:"instId"`
	ClOrdID string `json:"clOrdId,omitempty"`
	OrdID   string `json:"ordId,omitempty"`
}

type ModifyOrder struct {
	ID   string       `json:"id"`
	Op   EventType    `json:"op"`
	Args []ModifyArgs `json:"args"`
}

type ModifyArgs struct {
	Instrument    string `json:"instId"`
	OrderID       string `json:"ordId,omitempty"`
	ClientOrderID string `json:"clOrdId,omitempty"`
	NewPrice      string `json:"newPx"`
	NewSize       string `json:"newSz"`
	RequestID     string `json:"reqId"`
}

type CreateOrder struct {
	ID   string            `json:"id"`
	Op   EventType         `json:"op"`
	Args []CreateOrderArgs `json:"args"`
}

type CreateOrderArgs struct {
	Instrument     string `json:"instId"`
	ClientOrderID  string `json:"clOrdId"`
	Side           string `json:"side"`
	Type           string `json:"ordType"`
	Quantity       string `json:"sz"`
	Price          string `json:"px,omitempty"`
	TdMode         string `json:"tdMode"`
	TargetCurrency string `json:"tgtCcy,omitempty"`
}
