package clients_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/clients"
	"studentgit.kata.academy/quant/torque/internal/adapters/clients/types/okx"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

func setup() (clients.OKXApi, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case string(okx.Order):
			if r.Method == http.MethodPost {
				createSpotOrder(w, r)
			}
			if r.Method == http.MethodGet {
				getOrder(w, r)
			}
		case string(okx.CancelOrder):
			cancelOrder(w, r)
		case string(okx.AmendOrder):
			amendOrder(w, r)
		case string(okx.OrderHistory):
			if r.URL.Query().Get("instType") == string(okx.SWAP) {
				getFuturesOrderHistory(w, r)
			}
			if r.URL.Query().Get("instType") == string(okx.SPOT) {
				getSpotOrderHistory(w, r)
			}
		case string(okx.OpenOrders):
			if r.URL.Query().Get("instType") == string(okx.SWAP) {
				getFuturesOpenOrders(w, r)
			}
			if r.URL.Query().Get("instType") == string(okx.SPOT) {
				getSpotOpenOrders(w, r)
			}
		case string(okx.Balances):
			getBalances(w, r)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	cfg := config.HTTPClient{
		APIKey:        "test_api_key",
		APISecret:     "test_api_secret",
		APIPassphrase: "test_api_passphrase",
		URL:           server.URL,
		Timeout:       1 * time.Second,
	}
	api := *clients.NewOKXAPI(logster.Discard, cfg)

	return api, server
}

func TestGetBalances(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()
	response, err := api.GetBalances(context.Background())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func getBalances(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v5/account/balance" {
		http.Error(w, "URL path is not as expected", http.StatusBadRequest)
		return
	}

	response := `{
    "code": "0",
    "data": [
		{
		"adjEq": "55415.624719833286",
		"borrowFroz": "0",
		"details": 
		[
		{
		"availBal": "4834.317093622894",
		"availEq": "4834.3170936228935",
		"borrowFroz": "0",
		"cashBal": "4850.435693622894",
		"ccy": "USDT",
		"crossLiab": "0",
		"disEq": "4991.542013297616",
		"eq": "4992.890093622894",
		"eqUsd": "4991.542013297616",
		"fixedBal": "0",
		"frozenBal": "158.573",
		"imr": "",
		"interest": "0",
		"isoEq": "0",
		"isoLiab": "0",
		"isoUpl": "0",
		"liab": "0",
		"maxLoan": "0",
		"mgnRatio": "",
		"mmr": "",
		"notionalLever": "",
		"ordFrozen": "0",
		"rewardBal": "0",
		"spotInUseAmt": "",
		"spotIsoBal": "0",
		"stgyEq": "150",
		"twap": "0",
		"uTime": "1705449605015",
		"upl": "-7.545600000000006",
		"uplLiab": "0"
		}
		],
		"imr": "8.57068529",
		"isoEq": "0",
		"mgnRatio": "143682.59776662575",
		"mmr": "0.3428274116",
		"notionalUsd": "85.7068529",
		"ordFroz": "0",
		"totalEq": "55837.43556134779",
		"uTime": "1705474164160",
		"upl": "-7.543562688000006"
        }
    ],
    "msg": ""
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestGetSpotOpenOrders(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()
	response, err := api.GetSpotOpenOrders(context.Background())
	// Assertions
	v := response[0]
	assert.Equal(t, domain.OrderID("301835739059335168"), v.OrderID)
	assert.Equal(t, domain.Side("BUY"), v.Side)
	assert.Equal(t, domain.Symbol("BTC-USDT"), v.Instrument)
	floatPrice, _ := v.Price.Float64()
	assert.Equal(t, 59200.0, floatPrice)
	floatSize, _ := v.OrigQty.Float64()
	assert.Equal(t, 1.0, floatSize)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func getSpotOpenOrders(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v5/trade/orders-pending" {
		http.Error(w, "URL path is not as expected", http.StatusBadRequest)
		return
	}
	if r.URL.RawQuery != "instType=SPOT" {
		http.Error(w, "URL query parameters are not as expected", http.StatusBadRequest)
		return
	}

	response := `{
    "code": "0",
    "msg": "",
    "data": [
        {
		"accFillSz": "0",
		"avgPx": "",
		"cTime": "1618235248028",
		"category": "normal",
		"ccy": "",
		"clOrdId": "",
		"fee": "0",
		"feeCcy": "BTC",
		"fillPx": "",
		"fillSz": "0",
		"fillTime": "",
		"instId": "BTC-USDT",
		"instType": "SPOT",
		"lever": "5.6",
		"ordId": "301835739059335168",
		"ordType": "limit",
		"pnl": "0",
		"posSide": "net",
		"px": "59200",
		"pxUsd":"",
		"pxVol":"",
		"pxType":"",
		"rebate": "0",
		"rebateCcy": "USDT",
		"side": "buy",
		"attachAlgoClOrdId": "",
		"slOrdPx": "",
		"slTriggerPx": "",
		"slTriggerPxType": "last",
		"attachAlgoOrds": [],
		"state": "live",
		"stpId": "",
		"stpMode": "",
		"sz": "1",
		"tag": "",
		"tgtCcy": "",
		"tdMode": "cross",
		"source":"",
		"tpOrdPx": "",
		"tpTriggerPx": "",
		"tpTriggerPxType": "last",
		"tradeId": "",
		"reduceOnly": "false",
		"quickMgnType": "",
		"algoClOrdId": "",
		"algoId": "",
		"uTime": "1618235248028",
		"isTpLimit": "false"
        }
    ]
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestGetFuturesOpenOrders(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()
	response, err := api.GetFuturesOpenOrders(context.Background())

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func getFuturesOpenOrders(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v5/trade/orders-pending" {
		http.Error(w, "URL path is not as expected", http.StatusBadRequest)
		return
	}
	if r.URL.RawQuery != "instType=SWAP" {
		http.Error(w, "URL query parameters are not as expected", http.StatusBadRequest)
		return
	}

	response := `{
    "code": "0",
    "msg": "",
    "data": [
        {
		"accFillSz": "0",
		"avgPx": "",
		"cTime": "1618235248028",
		"category": "normal",
		"ccy": "",
		"clOrdId": "",
		"fee": "0",
		"feeCcy": "BTC",
		"fillPx": "",
		"fillSz": "0",
		"fillTime": "",
		"instId": "BTC-USDT",
		"instType": "SPOT",
		"lever": "5.6",
		"ordId": "301835739059335168",
		"ordType": "limit",
		"pnl": "0",
		"posSide": "net",
		"px": "59200",
		"pxUsd":"",
		"pxVol":"",
		"pxType":"",
		"rebate": "0",
		"rebateCcy": "USDT",
		"side": "buy",
		"attachAlgoClOrdId": "",
		"slOrdPx": "",
		"slTriggerPx": "",
		"slTriggerPxType": "last",
		"attachAlgoOrds": [],
		"state": "live",
		"stpId": "",
		"stpMode": "",
		"sz": "1",
		"tag": "",
		"tgtCcy": "",
		"tdMode": "cross",
		"source":"",
		"tpOrdPx": "",
		"tpTriggerPx": "",
		"tpTriggerPxType": "last",
		"tradeId": "",
		"reduceOnly": "false",
		"quickMgnType": "",
		"algoClOrdId": "",
		"algoId": "",
		"uTime": "1618235248028",
		"isTpLimit": "false"
		}
    ]
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestGetSpotOrderHistory(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()

	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2021, 6, 2, 0, 0, 0, 0, time.UTC)
	response, err := api.GetSpotOrderHistory(context.Background(), startTime, endTime)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func getSpotOrderHistory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v5/trade/orders-history" {
		http.Error(w, "URL path is not as expected", http.StatusBadRequest)
		return
	}
	if r.URL.RawQuery != "begin=1622505600000&end=1622592000000&instType=SPOT&limit=100" {
		http.Error(w, "URL query parameters are not as expected", http.StatusBadRequest)
		return
	}

	response := `{
    "code": "0",
    "data": [
        {
		"accFillSz": "0.00192834",
		"algoClOrdId": "",
		"algoId": "",
		"attachAlgoClOrdId": "",
		"attachAlgoOrds": [],
		"avgPx": "51858",
		"cTime": "1708587373361",
		"cancelSource": "",
		"cancelSourceReason": "",
		"category": "normal",
		"ccy": "",
		"clOrdId": "",
		"fee": "-0.00000192834",
		"feeCcy": "BTC",
		"fillPx": "51858",
		"fillSz": "0.00192834",
		"fillTime": "1708587373361",
		"instId": "BTC-USDT",
		"instType": "SPOT",
		"lever": "",
		"linkedAlgoOrd": {
		"algoId": ""
		},
		"ordId": "680800019749904384",
		"ordType": "market",
		"pnl": "0",
		"posSide": "",
		"px": "",
		"pxType": "",
		"pxUsd": "",
		"pxVol": "",
		"quickMgnType": "",
		"rebate": "0",
		"rebateCcy": "USDT",
		"reduceOnly": "false",
		"side": "buy",
		"slOrdPx": "",
		"slTriggerPx": "",
		"slTriggerPxType": "",
		"source": "",
		"state": "filled",
		"stpId": "",
		"stpMode": "",
		"sz": "100",
		"tag": "",
		"tdMode": "cash",
		"tgtCcy": "quote_ccy",
		"tpOrdPx": "",
		"tpTriggerPx": "",
		"tpTriggerPxType": "",
		"tradeId": "744876980",
		"uTime": "1708587373362",
		"isTpLimit": "false"
		}
    	],
    	"msg": ""
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestGetFuturesOrderHistory(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()
	startTime := time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2021, 6, 2, 0, 0, 0, 0, time.UTC)
	response, err := api.GetFuturesOrderHistory(context.Background(), startTime, endTime)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func getFuturesOrderHistory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v5/trade/orders-history" {
		http.Error(w, "URL path is not as expected", http.StatusBadRequest)
		return
	}
	if r.URL.RawQuery != "begin=1622505600000&end=1622592000000&instType=SWAP&limit=100" {
		http.Error(w, "URL query parameters are not as expected", http.StatusBadRequest)
		return
	}

	response := `{
		"code": "0",
		"msg": "",
		"data": [
		{
    	"instId": "BTC-USDT",
		"ordId": "123456",
		"clOrdId": "b15",
		"side": "buy",
		"ordType": "limit",
		"px": "2.15",
		"sz": "2",
		"state": "filled"
    	}
		]
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestGetOrder(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()
	response, err := api.GetOrder(context.Background(), "BTC-USDT", "123456", "b15")

	expectedResponse := domain.OpenOrder{
		ClientOrderID: "15",
		OrderID:       "680800019749904384",
		Side:          "BUY",
		Instrument:    "BTC-USDT",
		Type:          "MARKET",
		TimeInForce:   "",
		OrigQty:       decimal.New(100, 0),
		Price:         decimal.Decimal{},
		ExecutedQty:   decimal.New(192834, -8),
		Created:       time.Date(2024, 02, 22, 7, 36, 13, 361000000, time.UTC),
		UpdatedAt:     time.Date(2024, 02, 22, 7, 36, 13, 362000000, time.UTC),
	}
	assert.Equal(t, expectedResponse.ClientOrderID, response.ClientOrderID)
	assert.Equal(t, expectedResponse.OrderID, response.OrderID)
	assert.Equal(t, expectedResponse.Side, response.Side)
	assert.Equal(t, expectedResponse.Instrument, response.Instrument)
	assert.Equal(t, expectedResponse.Type, response.Type)
	assert.Equal(t, expectedResponse.TimeInForce, response.TimeInForce)
	assert.Equal(t, expectedResponse.OrigQty, response.OrigQty)
	assert.Equal(t, expectedResponse.Price, response.Price)
	assert.Equal(t, expectedResponse.ExecutedQty, response.ExecutedQty)
	assert.Equal(t, expectedResponse.Created.Unix(), response.Created.Unix())
	assert.Equal(t, expectedResponse.UpdatedAt.Unix(), response.UpdatedAt.Unix())
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v5/trade/order" {
		http.Error(w, "URL path is not as expected", http.StatusBadRequest)
		return
	}
	if r.URL.RawQuery != "clOrdId=b15&instId=BTC-USDT&ordId=123456" {
		http.Error(w, "URL query parameters are not as expected", http.StatusBadRequest)
		return
	}
	response := `{
		"code": "0",
		"data": [
		{
		"accFillSz": "0.00192834",
		"algoClOrdId": "",
		"algoId": "",
		"attachAlgoClOrdId": "",
		"attachAlgoOrds": [],
		"avgPx": "51858",
		"cTime": "1708587373361",
		"cancelSource": "",
		"cancelSourceReason": "",
		"category": "normal",
		"ccy": "",
		"clOrdId": "15",
		"fee": "-0.00000192834",
		"feeCcy": "BTC",
		"fillPx": "51858",
		"fillSz": "0.00192834",
		"fillTime": "1708587373361",
		"instId": "BTC-USDT",
		"instType": "SPOT",
		"isTpLimit": "false",
		"lever": "",
		"linkedAlgoOrd": {
		"algoId": ""
		},
		"ordId": "680800019749904384",
		"ordType": "market",
		"pnl": "0",
		"posSide": "net",
		"px": "",
		"pxType": "",
		"pxUsd": "",
		"pxVol": "",
		"quickMgnType": "",
		"rebate": "0",
		"rebateCcy": "USDT",
		"reduceOnly": "false",
		"side": "buy",
		"slOrdPx": "",
		"slTriggerPx": "",
		"slTriggerPxType": "",
		"source": "",
		"state": "filled",
		"stpId": "",
		"stpMode": "",
		"sz": "100",
		"tag": "",
		"tdMode": "cash",
		"tgtCcy": "quote_ccy",
		"tpOrdPx": "",
		"tpTriggerPx": "",
		"tpTriggerPxType": "",
		"tradeId": "744876980",
		"uTime": "1708587373362"
		}
		],
		"msg": ""
	}`

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestCreateSpotOrder(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()
	request := domain.OrderRequest{
		Symbol:        "BTC-USDT",
		Side:          domain.SideBuy,
		Type:          domain.OrderTypeLimit,
		TimeInForce:   "",
		Quantity:      decimal.NewFromInt(2),
		Price:         decimal.NewFromFloat(2.15),
		ClientOrderID: "b15",
	}

	response, err := api.CreateSpotOrder(context.Background(), request)
	// Assertions
	assert.Equal(t, domain.OrderID("312269865356374016"), response.OrderID)
	assert.Equal(t, domain.ClientOrderID("oktswap6"), response.ClientOrderID)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func createSpotOrder(w http.ResponseWriter, r *http.Request) {
	expectedRequest := okx.RequestCreateOrder{
		Instrument:     "BTC-USDT",
		ClientOrderID:  "b15",
		Side:           "buy",
		Type:           "limit",
		Quantity:       "2",
		Price:          "2.15",
		TdMode:         "cash",
		TargetCurrency: "",
	}

	var request okx.RequestCreateOrder
	_ = json.NewDecoder(r.Body).Decode(&request)

	if request != expectedRequest {
		http.Error(w, "wrong request", http.StatusBadRequest)
		return
	}
	response := `{
		"code": "0",
		"msg": "",
		"data": [
    	{
		"clOrdId": "oktswap6",
		"ordId": "312269865356374016",
		"tag": "",
		"sCode": "0",
		"sMsg": ""
    	}
		],
		"inTime": "1695190491421339",
		"outTime": "1695190491423240"
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}

func TestCancelOrder(t *testing.T) {
	t.Parallel()
	api, server := setup()
	defer server.Close()

	expectedResponse := domain.OKXOrderResponseAck{
		ClientOrderID: "b15",
		OrderID:       "123456",
	}
	response, err := api.CancelOrder(context.Background(), "BTC-USDT", "123456", "b15")

	// Assertions
	assert.Equal(t, expectedResponse.ClientOrderID, response.ClientOrderID)
	assert.Equal(t, expectedResponse.OrderID, response.OrderID)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	expectedRequest := okx.RequestCancelOrder{
		Instrument:    "BTC-USDT",
		OrderID:       "123456",
		ClientOrderID: "b15",
	}
	var request okx.RequestCancelOrder
	_ = json.NewDecoder(r.Body).Decode(&request)

	if request != expectedRequest {
		http.Error(w, "wrong request", http.StatusBadRequest)
		return
	}

	response := `{
		"code": "0",
		"msg": "",
		"data": [
    	{
      	"clOrdId": "b15",
      	"ordId": "123456",
      	"sCode": "0",
      	"sMsg": ""
    	}
		]
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}
func TestAmendOrder(t *testing.T) {
	t.Parallel()
	api, server := setup()

	defer server.Close()

	err := api.AmendOrder(context.Background(),
		"BTC-USDT",
		"123456",
		"b15",
		decimal.NewFromInt(3),
		decimal.NewFromFloat(2.20))

	assert.NoError(t, err)
}

func amendOrder(w http.ResponseWriter, r *http.Request) {
	expectedRequest := okx.RequestAmendOrder{
		Instrument:    "BTC-USDT",
		OrderID:       "123456",
		ClientOrderID: "b15",
		NewPrice:      "2.2",
		NewSize:       "3",
		RequestID:     "b15",
	}

	var request okx.RequestAmendOrder
	_ = json.NewDecoder(r.Body).Decode(&request)

	if request != expectedRequest {
		http.Error(w, "wrong request", http.StatusBadRequest)
		return
	}

	response := `{
    	"code":"0",
    	"msg":"",
    	"data":[
		{
		"clOrdId":"",
		"ordId":"12344",
		"reqId":"b12344",
		"sCode":"0",
		"sMsg":""
        }
    	],
    	"inTime": "1695190491421339",
    	"outTime": "1695190491423240"
	}`
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}
