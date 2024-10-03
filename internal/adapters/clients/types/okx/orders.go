package okx

type OpenOrdersResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		AccFillSz         string `json:"accFillSz"`
		AvgPx             string `json:"avgPx"`
		CTime             string `json:"cTime"`
		Category          string `json:"category"`
		Ccy               string `json:"ccy"`
		ClOrdID           string `json:"clOrdId"`
		Fee               string `json:"fee"`
		FeeCcy            string `json:"feeCcy"`
		FillPx            string `json:"fillPx"`
		FillSz            string `json:"fillSz"`
		FillTime          string `json:"fillTime"`
		InstID            string `json:"instId"`
		InstType          string `json:"instType"`
		Lever             string `json:"lever"`
		OrdID             string `json:"ordId"`
		OrdType           string `json:"ordType"`
		Pnl               string `json:"pnl"`
		PosSide           string `json:"posSide"`
		Px                string `json:"px"`
		PxUsd             string `json:"pxUsd"`
		PxVol             string `json:"pxVol"`
		PxType            string `json:"pxType"`
		Rebate            string `json:"rebate"`
		RebateCcy         string `json:"rebateCcy"`
		Side              string `json:"side"`
		AttachAlgoClOrdID string `json:"attachAlgoClOrdId"`
		SlOrdPx           string `json:"slOrdPx"`
		SlTriggerPx       string `json:"slTriggerPx"`
		SlTriggerPxType   string `json:"slTriggerPxType"`
		State             string `json:"state"`
		StpID             string `json:"stpId"`
		StpMode           string `json:"stpMode"`
		Sz                string `json:"sz"`
		Tag               string `json:"tag"`
		TgtCcy            string `json:"tgtCcy"`
		TdMode            string `json:"tdMode"`
		Source            string `json:"source"`
		TpOrdPx           string `json:"tpOrdPx"`
		TpTriggerPx       string `json:"tpTriggerPx"`
		TpTriggerPxType   string `json:"tpTriggerPxType"`
		TradeID           string `json:"tradeId"`
		ReduceOnly        string `json:"reduceOnly"`
		QuickMgnType      string `json:"quickMgnType"`
		AlgoClOrdID       string `json:"algoClOrdId"`
		AlgoID            string `json:"algoId"`
		UTime             string `json:"uTime"`
	} `json:"data"`
}

type OrderHistoryResponse struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Data []HistoricalOrder `json:"data"`
}

type HistoricalOrder struct {
	InstType           string `json:"instType"`
	InstID             string `json:"instId"`
	Ccy                string `json:"ccy"`
	OrdID              string `json:"ordId"`
	ClOrdID            string `json:"clOrdId"`
	Tag                string `json:"tag"`
	Px                 string `json:"px"`
	PxUsd              string `json:"pxUsd"`
	PxVol              string `json:"pxVol"`
	PxType             string `json:"pxType"`
	Sz                 string `json:"sz"`
	OrdType            string `json:"ordType"`
	Side               string `json:"side"`
	PosSide            string `json:"posSide"`
	TdMode             string `json:"tdMode"`
	AccFillSz          string `json:"accFillSz"`
	FillPx             string `json:"fillPx"`
	TradeID            string `json:"tradeId"`
	FillSz             string `json:"fillSz"`
	FillTime           string `json:"fillTime"`
	State              string `json:"state"`
	AvgPx              string `json:"avgPx"`
	Lever              string `json:"lever"`
	AttachAlgoClOrdID  string `json:"attachAlgoClOrdId"`
	TpTriggerPx        string `json:"tpTriggerPx"`
	TpTriggerPxType    string `json:"tpTriggerPxType"`
	TpOrdPx            string `json:"tpOrdPx"`
	SlTriggerPx        string `json:"slTriggerPx"`
	SlTriggerPxType    string `json:"slTriggerPxType"`
	SlOrdPx            string `json:"slOrdPx"`
	StpID              string `json:"stpId"`
	StpMode            string `json:"stpMode"`
	FeeCcy             string `json:"feeCcy"`
	Fee                string `json:"fee"`
	RebateCcy          string `json:"rebateCcy"`
	Source             string `json:"source"`
	Rebate             string `json:"rebate"`
	TgtCcy             string `json:"tgtCcy"`
	Pnl                string `json:"pnl"`
	Category           string `json:"category"`
	ReduceOnly         string `json:"reduceOnly"`
	CancelSource       string `json:"cancelSource"`
	CancelSourceReason string `json:"cancelSourceReason"`
	AlgoClOrdID        string `json:"algoClOrdId"`
	AlgoID             string `json:"algoId"`
	UTime              string `json:"uTime"`
	CTime              string `json:"cTime"`
}

type InstrumentType string

const (
	SPOT    InstrumentType = "SPOT"
	MARGIN  InstrumentType = "MARGIN"
	SWAP    InstrumentType = "SWAP"
	FUTURES InstrumentType = "FUTURES"
	OPTION  InstrumentType = "OPTION"
)

type BaseResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []interface{} `json:"data"`
}

type CreateOrderResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ClOrdID string `json:"clOrdId"`
		OrdID   string `json:"ordId"`
		Tag     string `json:"tag"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

type CancelOrderResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ClOrdID string `json:"clOrdId"`
		OrdID   string `json:"ordId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

type AmendOrderResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ClOrdID string `json:"clOrdId"`
		OrdID   string `json:"ordId"`
		ReqID   string `json:"reqId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

type CancelAllOrdersAfterMessage struct {
	Code string               `json:"code"`
	Msg  string               `json:"msg"`
	Data []CancelAllAfterData `json:"data"`
}

type CancelAllAfterData struct {
	TriggerTime string `json:"triggerTime"`
	TS          string `json:"ts"`
}

type FillsHistoryResponse struct {
	Code string           `json:"code"`
	Msg  string           `json:"msg"`
	Data []HistoricalFill `json:"data"`
}

type HistoricalFill struct {
	InstType    string `json:"instType"`
	InstID      string `json:"instId"`
	TradeID     string `json:"tradeId"`
	OrdID       string `json:"ordId"`
	ClOrdID     string `json:"clOrdId"`
	BillID      string `json:"billId"`
	Tag         string `json:"tag"`
	FillPx      string `json:"fillPx"`
	FillSz      string `json:"fillSz"`
	FillIdxPx   string `json:"fillIdxPx"`
	FillPnl     string `json:"fillPnl"`
	FillPxVol   string `json:"fillPxVol"`
	FillPxUsd   string `json:"fillPxUsd"`
	FillMarkVol string `json:"fillMarkVol"`
	FillFwdPx   string `json:"fillFwdPx"`
	FillMarkPx  string `json:"fillMarkPx"`
	Side        string `json:"side"`
	PosSide     string `json:"posSide"`
	ExecType    string `json:"execType"`
	FeeCcy      string `json:"feeCcy"`
	Fee         string `json:"fee"`
	TS          string `json:"ts"`
	FillTime    string `json:"fillTime"`
}
