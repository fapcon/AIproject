package clients

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"studentgit.kata.academy/quant/torque/config"

	"github.com/google/go-querystring/query"
	"github.com/mailru/easyjson"
	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/adapters/clients/types/okx"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const (
	maxResultsPerPage = 100
)

type OKXApi struct {
	logger         logster.Logger
	apiKey         string
	apiSecret      string
	host           string
	apiPassphrase  string
	client         *http.Client
	timeout        time.Duration
	rateLimiterMap *RateLimiterMap
}

func NewOKXAPI(logger logster.Logger, config config.HTTPClient) *OKXApi {
	return &OKXApi{
		logger:         logger.WithField(f.APIClient, "okx_api_client"),
		client:         &http.Client{},
		apiKey:         config.APIKey,
		apiSecret:      config.APISecret,
		apiPassphrase:  config.APIPassphrase,
		host:           config.URL,
		timeout:        config.Timeout,
		rateLimiterMap: NewRateLimiterMap(),
	}
}

func (api *OKXApi) request(
	ctx context.Context, method string, path okx.Endpoint,
	body string, bodyMarshaller easyjson.Marshaler, unmarshaler easyjson.Unmarshaler,
	isPrivate bool,
) (interface{}, error) {
	startTime := time.Now()
	defer LogElapsed(api.logger, startTime, method, string(path))
	requestBody := []byte("")
	if method == http.MethodGet && body != "" {
		path = okx.Endpoint(fmt.Sprintf("%s?%s", path, body))
	} else if method == http.MethodPost && bodyMarshaller != nil {
		var err error
		requestBody, err = easyjson.Marshal(bodyMarshaller)
		if err != nil {
			api.logger.WithError(err).Errorf("Error marshaling JSON body")
			return nil, err
		}
	}
	fullPath := fmt.Sprintf("%s%s", api.host, path)

	ctxTimeout, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxTimeout, method, fullPath, bytes.NewBuffer(requestBody))
	if err != nil {
		api.logger.WithError(err).Errorf("Error preparing request")
		return nil, err
	}
	if isPrivate {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
		var signBody string
		if method == http.MethodGet {
			signBody = fmt.Sprintf("%v%v%v", timestamp, method, path)
		} else if method == http.MethodPost {
			signBody = fmt.Sprintf("%v%v%v%v", timestamp, method, path, string(requestBody))
		}
		sig := api.generateSignature(signBody)

		req.Header.Set("OK-ACCESS-KEY", api.apiKey)
		req.Header.Set("OK-ACCESS-SIGN", sig)
		req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Set("OK-ACCESS-PASSPHRASE", api.apiPassphrase)
		req.Header.Set("Content-Type", "application/json")
	}
	api.logger.WithField(f.URL, req.URL).Infof("Sending request")

	resp, err := api.client.Do(req) //nolint:bodyclose // body is closed in closure bellow
	if err != nil {
		api.logger.WithError(err).Errorf("Error sending request")
		return nil, err
	}
	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		api.logger.WithError(err).Errorf("Error reading response body")
		return nil, err
	}
	defer func(body io.ReadCloser) {
		errX := body.Close()
		if errX != nil {
			api.logger.WithError(err).Errorf("Error closing response body")
		}
	}(resp.Body)
	partialErrResp := &okx.BaseResponse{} //nolint:exhaustruct // OK
	decodeErr := easyjson.Unmarshal(bodyBytes, partialErrResp)
	if decodeErr != nil {
		api.logger.WithError(decodeErr).Errorf("Error decoding response")
		return nil, decodeErr
	}
	if partialErrResp.Code != "0" {
		api.logger.WithField(f.OKXError, partialErrResp).Errorf("Error in OKX response")
		return nil, fmt.Errorf("status code: %v, code: %v, message: %v, data: %+v",
			resp.StatusCode, partialErrResp.Code, partialErrResp.Msg, partialErrResp.Data)
	}
	err = easyjson.Unmarshal(bodyBytes, unmarshaler)
	if err != nil {
		api.logger.Errorf("Error handling response, %v", err)
	}
	return unmarshaler, err
}

func (api *OKXApi) CreateSpotOrder(
	ctx context.Context, request domain.OrderRequest,
) (domain.OKXOrderResponseAck, error) {
	if api.rateLimiterMap.Get(okx.Order).Allow() {
		response, err := api.createOrder(ctx, request, true)
		if err != nil {
			return domain.OKXOrderResponseAck{}, err
		}
		return okx.CreateOrderToDomain(*response), nil
	}
	return domain.OKXOrderResponseAck{}, domain.ErrPlaceOrderRateLimitExceed
}

func (api *OKXApi) CreateFuturesOrder(
	ctx context.Context, request domain.OrderRequest,
) (domain.OKXOrderResponseAck, error) {
	if api.rateLimiterMap.Get(okx.Order).Allow() {
		response, err := api.createOrder(ctx, request, false)
		if err != nil {
			return domain.OKXOrderResponseAck{}, err
		}
		return okx.CreateOrderToDomain(*response), nil
	}
	return domain.OKXOrderResponseAck{}, domain.ErrPlaceOrderRateLimitExceed
}

func (api *OKXApi) createOrder(
	ctx context.Context, request domain.OrderRequest, spot bool,
) (*okx.CreateOrderResponse, error) {
	var mode string
	if spot {
		mode = "cash"
	} else {
		mode = "cross"
	}
	//nolint:exhaustruct // TargetCurrency is filled bellow
	body := okx.RequestCreateOrder{
		Instrument:    string(request.Symbol),
		ClientOrderID: string(request.ClientOrderID),
		Side:          okx.SideFromDomain(request.Side),
		Type:          okx.OrderTypeFromDomain(request.Type),
		Quantity:      request.Quantity.String(),
		Price:         request.Price.String(),
		TdMode:        mode,
		// TIF is a part of Type in OKX
	}
	if request.Type == domain.OrderTypeMarket {
		body.TargetCurrency = "base_ccy"
	}
	response, err := api.request(ctx, http.MethodPost, okx.Order,
		"", body, &okx.CreateOrderResponse{}, true) //nolint:exhaustruct // OK
	if err != nil {
		return &okx.CreateOrderResponse{}, err //nolint:exhaustruct // OK
	}
	return response.(*okx.CreateOrderResponse), nil
}

func (api *OKXApi) CancelOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) (domain.OKXOrderResponseAck, error) {
	if api.rateLimiterMap.Get(okx.CancelOrder).Allow() {
		response, err := api.cancelOrder(ctx, string(symbol), string(orderID), string(clientOrderID))
		if err != nil {
			return domain.OKXOrderResponseAck{}, err
		}
		return okx.CancelOrderToDomain(*response), nil
	}
	return domain.OKXOrderResponseAck{}, domain.ErrCancelOrderRateLimitExceed
}

func (api *OKXApi) cancelOrder(
	ctx context.Context, symbol string, orderID, clientOrderID string,
) (*okx.CancelOrderResponse, error) {
	body := okx.RequestCancelOrder{
		Instrument:    symbol,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
	}
	response, err := api.request(ctx, http.MethodPost, okx.CancelOrder,
		"", body, &okx.CancelOrderResponse{}, true) //nolint:exhaustruct // OK
	if err != nil {
		return &okx.CancelOrderResponse{}, err //nolint:exhaustruct // OK
	}
	return response.(*okx.CancelOrderResponse), nil
}

func (api *OKXApi) AmendOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	if api.rateLimiterMap.Get(okx.AmendOrder).Allow() {
		return api.amendOrder(ctx, string(symbol), string(orderID), string(clientOrderID), newSize, newPrice)
	}
	return domain.ErrAmendOrderRateLimitExceed
}

func (api *OKXApi) amendOrder(
	ctx context.Context, symbol string, orderID, clientOrderID string, newSize, newPrice decimal.Decimal,
) error {
	body := okx.RequestAmendOrder{
		Instrument:    symbol,
		OrderID:       orderID,
		ClientOrderID: clientOrderID,
		NewPrice:      newPrice.String(),
		NewSize:       newSize.String(),
		RequestID:     clientOrderID,
	}
	response, err := api.request(ctx, http.MethodPost, okx.AmendOrder,
		"", body, &okx.AmendOrderResponse{}, true) //nolint:exhaustruct // OK
	api.logger.WithField(f.AmendOrder, response).Infof("Amend Order Response")
	if err != nil {
		return err
	}
	return nil
}

func (api *OKXApi) GetOrder(
	ctx context.Context, symbol domain.OKXSpotInstrument, orderID domain.OrderID, origClientOrderID domain.ClientOrderID,
) (domain.OpenOrder, error) {
	response, err := api.getOrder(ctx, string(symbol), string(orderID), string(origClientOrderID))
	if err != nil {
		return domain.OpenOrder{}, err
	}
	result := okx.OpenOrdersToDomain(*response)
	if len(result) == 0 {
		return domain.OpenOrder{}, nil //nolint:exhaustruct // OK
	}
	return result[0], nil
}

func (api *OKXApi) getOrder(
	ctx context.Context, symbol string, orderID, origClientOrderID string,
) (*okx.OpenOrdersResponse, error) {
	body := okx.RequestGetOrder{
		Instrument:    symbol,
		OrderID:       orderID,
		ClientOrderID: origClientOrderID,
	}
	values, _ := query.Values(body)
	response, err := api.request(ctx, http.MethodGet, okx.Order,
		values.Encode(), nil, &okx.OpenOrdersResponse{}, true) //nolint:exhaustruct // OK
	if err != nil {
		return &okx.OpenOrdersResponse{}, err //nolint:exhaustruct // OK
	}
	return response.(*okx.OpenOrdersResponse), nil
}

func (api *OKXApi) GetSpotOpenOrders(ctx context.Context) ([]domain.OpenOrder, error) {
	response, err := api.getOpenOrders(ctx, okx.SPOT)
	if err != nil {
		return []domain.OpenOrder{}, err
	}
	return okx.OpenOrdersToDomain(*response), nil
}

func (api *OKXApi) GetFuturesOpenOrders(ctx context.Context) ([]domain.OpenOrder, error) {
	response, err := api.getOpenOrders(ctx, okx.SWAP)
	if err != nil {
		return []domain.OpenOrder{}, err
	}
	return okx.OpenOrdersToDomain(*response), nil
}

func (api *OKXApi) getOrderHistory(
	ctx context.Context, instrumentType okx.InstrumentType, startTime, endTime time.Time,
) ([]okx.HistoricalOrder, error) {
	allOrders := make([]okx.HistoricalOrder, 0, maxResultsPerPage)
	var after string // keep track of the last order ID for pagination
	for {
		body := okx.RequestOrderHistory{
			InstrumentType: string(instrumentType),
			StartTime:      strconv.FormatInt(startTime.UnixNano()/int64(time.Millisecond), 10),
			EndTime:        strconv.FormatInt(endTime.UnixNano()/int64(time.Millisecond), 10),
			Limit:          strconv.Itoa(maxResultsPerPage),
			After:          after,
		}
		values, _ := query.Values(body)

		response, err := api.request(ctx, http.MethodGet, okx.OrderHistory,
			values.Encode(), nil, &okx.OrderHistoryResponse{}, true) //nolint:exhaustruct // OK
		if err != nil {
			return nil, err
		}
		orderHistoryResponse, ok := response.(*okx.OrderHistoryResponse)
		if !ok {
			return nil, errors.New("failed to convert response to OrderHistoryResponse type")
		}
		for _, order := range orderHistoryResponse.Data {
			if okx.StatusToDomain(order.State) == domain.OrderStatusFilled ||
				okx.StatusToDomain(order.State) == domain.OrderStatusPartiallyFilled {
				allOrders = append(allOrders, order)
			}
		}

		// Exit loop if less than the max number of results were returned, indicating no more pages
		if len(orderHistoryResponse.Data) < maxResultsPerPage {
			break
		}
		// For next iteration, set the 'after' to the last order ID from the current response
		after = orderHistoryResponse.Data[len(orderHistoryResponse.Data)-1].OrdID
	}
	return allOrders, nil
}

func (api *OKXApi) GetSpotOrderHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	orders, err := api.getOrderHistory(ctx, okx.SPOT, startTime, endTime)
	if err != nil {
		return []domain.Order{}, err
	}
	return okx.OrderHistoryToDomain(orders), nil
}

func (api *OKXApi) GetFuturesOrderHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	orders, err := api.getOrderHistory(ctx, okx.SWAP, startTime, endTime)
	if err != nil {
		return []domain.Order{}, err
	}
	return okx.OrderHistoryToDomain(orders), nil
}

func (api *OKXApi) getOpenOrders(
	ctx context.Context, instrumentType okx.InstrumentType,
) (*okx.OpenOrdersResponse, error) {
	body := okx.RequestOpenOrders{InstrumentType: string(instrumentType)}
	values, _ := query.Values(body)
	response, err := api.request(ctx, http.MethodGet, okx.OpenOrders,
		values.Encode(), nil, &okx.OpenOrdersResponse{}, true) //nolint:exhaustruct // OK
	if err != nil {
		return &okx.OpenOrdersResponse{}, err //nolint:exhaustruct // OK
	}
	api.logger.Infof("Open orders response: %+v", response)
	return response.(*okx.OpenOrdersResponse), nil
}

func (api *OKXApi) GetSpotFillsHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	orders, err := api.getFillsHistory(ctx, okx.SPOT, startTime, endTime)
	if err != nil {
		return []domain.Order{}, err
	}
	return okx.FillsHistoryToDomain(orders), nil
}

func (api *OKXApi) GetFuturesFillsHistory(ctx context.Context, startTime, endTime time.Time) ([]domain.Order, error) {
	orders, err := api.getFillsHistory(ctx, okx.FUTURES, startTime, endTime)
	if err != nil {
		return []domain.Order{}, err
	}
	return okx.FillsHistoryToDomain(orders), nil
}

func (api *OKXApi) getFillsHistory(
	ctx context.Context, instrumentType okx.InstrumentType, startTime, endTime time.Time,
) ([]okx.HistoricalFill, error) {
	allOrders := make([]okx.HistoricalFill, 0, maxResultsPerPage)
	var after string // keep track of the last order ID for pagination
	for {
		body := okx.RequestFillsHistory{
			InstrumentType: string(instrumentType),
			StartTime:      strconv.FormatInt(startTime.UnixNano()/int64(time.Millisecond), 10),
			EndTime:        strconv.FormatInt(endTime.UnixNano()/int64(time.Millisecond), 10),
			Limit:          strconv.Itoa(maxResultsPerPage),
			After:          after,
		}
		values, _ := query.Values(body)

		response, err := api.request(ctx, http.MethodGet, okx.FillsHistory,
			values.Encode(), nil, &okx.FillsHistoryResponse{}, true) //nolint:exhaustruct // OK
		if err != nil {
			return nil, err
		}
		fillsHistoryResponse, ok := response.(*okx.FillsHistoryResponse)
		if !ok {
			return nil, errors.New("failed to convert response to OrderHistoryResponse type")
		}
		allOrders = append(allOrders, fillsHistoryResponse.Data...)

		// Exit loop if less than the max number of results were returned, indicating no more pages
		if len(fillsHistoryResponse.Data) < maxResultsPerPage {
			break
		}
		// For next iteration, set the 'after' to the last order ID from the current response
		after = fillsHistoryResponse.Data[len(fillsHistoryResponse.Data)-1].OrdID
	}
	return allOrders, nil
}

func (api *OKXApi) GetBalances(ctx context.Context) ([]domain.Balance, error) {
	if api.rateLimiterMap.Get(okx.Balances).Allow() {
		response, err := api.getBalances(ctx)
		if err != nil {
			return []domain.Balance{}, err
		}
		return okx.BalancesToDomain(*response), nil
	}
	return []domain.Balance{}, domain.ErrBalancesRateLimitExceed
}

func (api *OKXApi) getBalances(ctx context.Context) (*okx.BalancesResponse, error) {
	response, err := api.request(ctx, http.MethodGet, okx.Balances,
		"", nil, &okx.BalancesResponse{}, true) //nolint:exhaustruct // OK
	if err != nil {
		return &okx.BalancesResponse{}, err //nolint:exhaustruct // OK
	}
	return response.(*okx.BalancesResponse), nil
}

func (api *OKXApi) getOrderBook(ctx context.Context, symbol string, limit int) (*okx.OrderBookResponse, error) {
	body := okx.RequestOrderBook{
		Instrument: symbol,
		Limit:      limit,
	}
	values, _ := query.Values(body)
	response, err := api.request(ctx, http.MethodGet, okx.OrderBook,
		values.Encode(), nil, &okx.OrderBookResponse{}, false) //nolint:exhaustruct // OK
	if err != nil {
		return &okx.OrderBookResponse{}, err //nolint:exhaustruct // OK
	}
	return response.(*okx.OrderBookResponse), nil
}

func (api *OKXApi) GetOrderBook(
	ctx context.Context, symbol domain.OKXSpotInstrument, limit int,
) (domain.OKXOrderBook, error) {
	response, err := api.getOrderBook(ctx, string(symbol), limit)
	if err != nil {
		return domain.OKXOrderBook{}, err
	}
	book, err := okx.OrderBookSnapshotToDomain(*response)
	if err != nil {
		return domain.OKXOrderBook{}, err
	}
	book.Symbol = symbol
	return book, nil
}

func (api *OKXApi) DeadMensSwitch(ctx context.Context, duration time.Duration) (time.Time, error) {
	body := okx.RequestCancelAllAfter{
		Timeout: strconv.Itoa(int(duration.Seconds())),
	}
	response, err := api.request(
		ctx, http.MethodPost, okx.CancelAllAfter, "", body, &okx.CancelAllOrdersAfterMessage{}, true, //nolint:exhaustruct,lll // OK
	)
	if err != nil {
		return time.Time{}, err
	}
	cancel, ok := response.(*okx.CancelAllOrdersAfterMessage)
	if !ok {
		return time.Time{}, errors.New("failed to convert response to CancelAllOrdersAfterMessage type")
	}
	parseInt, err := strconv.ParseInt(cancel.Data[0].TriggerTime, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(parseInt), nil
}

func (api *OKXApi) generateSignature(data string) string {
	mac := hmac.New(sha256.New, []byte(api.apiSecret))
	mac.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
