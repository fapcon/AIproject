package clients

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/mailru/easyjson"
	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/clients/types/bybit"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const (
	receiveWindow = 5000
)

type BybitAPI struct {
	logger         logster.Logger
	apiKey         string
	apiSecret      string
	host           string
	client         *http.Client
	timeout        time.Duration
	rateLimiterMap *RateLimiterMap
}

func NewBybitAPI(logger logster.Logger, config config.HTTPClient) *BybitAPI {
	return &BybitAPI{
		logger:         logger.WithField(f.APIClient, "bybit_api_client"),
		client:         &http.Client{},
		apiKey:         config.APIKey,
		apiSecret:      config.APISecret,
		host:           config.URL,
		timeout:        config.Timeout,
		rateLimiterMap: NewRateLimiterMap(),
	}
}

func (api *BybitAPI) request(
	ctx context.Context, method string, path bybit.Endpoint,
	paramsQuery string, inBodyParams easyjson.Marshaler, output easyjson.Unmarshaler,
	isPrivate bool,
) (interface{}, error) {
	startTime := time.Now()
	logger := api.logger.WithField("op", "clients.BybitAPI.request")
	defer LogElapsed(logger, startTime, method, string(path))
	requestBody := []byte("")
	if method == http.MethodGet && paramsQuery != "" {
		path = bybit.Endpoint(fmt.Sprintf("%s?%s", path, paramsQuery))
	} else if method == http.MethodPost && inBodyParams != nil {
		var err error
		requestBody, err = easyjson.Marshal(inBodyParams)
		if err != nil {
			api.logger.WithError(err).Errorf("Eror marshaling json body")
			return nil, err
		}
	}
	fullPath := fmt.Sprintf("%s%s", api.host, path)
	ctxTimeout, cancel := context.WithTimeout(ctx, api.timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctxTimeout, method, fullPath, bytes.NewBuffer(requestBody))
	if err != nil {
		api.logger.WithError(err).Errorf("Error preparing request")
		return nil, err
	}
	if isPrivate {
		timestamp := time.Now().UnixMilli()
		var signBody string
		if method == http.MethodGet {
			signBody = fmt.Sprintf(
				"%d%s%d%s", timestamp, api.apiKey, receiveWindow, paramsQuery,
			)
		} else if method == http.MethodPost {
			signBody = fmt.Sprintf(
				"%d%s%d%s", timestamp, api.apiKey, receiveWindow, string(requestBody),
			)
		}
		api.setPrivateHeader(request, signBody, timestamp)
	}

	api.logger.WithField(f.URL, request.URL).Infof("Sending request")

	response, err := api.client.Do(request)
	if err != nil {
		api.logger.WithError(err).Errorf("Error sending response")
		return nil, err
	}

	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(response.Body)
	if err != nil {
		api.logger.WithError(err).Errorf("Error reading response body")
		return nil, err
	}

	defer func() {
		if errX := response.Body.Close(); errX != nil {
			api.logger.WithError(errX).Errorf("Error closing response body")
		}
	}()

	var partialErrResp bybit.CommonResponse
	decodeErr := easyjson.Unmarshal(bodyBytes, &partialErrResp)
	if decodeErr != nil {
		api.logger.WithError(decodeErr).Errorf("Error decoding response")
		return nil, decodeErr
	}
	if partialErrResp.Code != 0 {
		api.logger.WithField(f.BybitError, partialErrResp).Errorf("Error in Bybit response")
		return nil, fmt.Errorf(
			"code: %v, message: %v, extraInfo: %v",
			partialErrResp.Code, partialErrResp.Msg, partialErrResp.Info)
	}
	err = easyjson.Unmarshal(bodyBytes, output)
	if err != nil {
		api.logger.Errorf("Error handling response, %v", err)
	}
	return output, nil
}

func (api *BybitAPI) createOrder(
	ctx context.Context, request domain.OrderRequest, spot bool,
) (*bybit.CreateOrderResponse, error) {
	var category bybit.InstrumentType
	logger := api.logger.WithField("op", "clients.BybitAPI.createOrder")
	switch spot {
	case true:
		category = bybit.SPOT
	case false:
		category = bybit.LINEAR
	}
	inBodyParams := bybit.CreateOrderParams{
		Category:    category.GetString(),
		Symbol:      bybit.SymbolFromDomain(request.Symbol),
		Side:        bybit.SideFromDomain(request.Side),
		OrderType:   bybit.OrderTypeFromDomain(request.Type),
		Qty:         request.Quantity.String(),
		Price:       request.Price.String(),
		TimeInForce: bybit.TimeInForceFromDomain(request.TimeInForce),
	}

	var output bybit.CreateOrderResponse
	data, err := api.request(
		ctx, http.MethodPost, bybit.CreateOrder, "", inBodyParams, &output, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error creating order")
		return nil, err
	}
	return data.(*bybit.CreateOrderResponse), nil
}

func (api *BybitAPI) CreateSpotOrder(
	ctx context.Context, request domain.OrderRequest,
) (domain.BybitOrderResponseAck, error) {
	// Only difference between CreateFuturesOrder and CreateSpotOrder is in parameter "spot" here
	response, err := api.createOrder(ctx, request, true)
	if err != nil {
		return domain.BybitOrderResponseAck{}, err
	}
	return domain.BybitOrderResponseAck{
		OrderID:       domain.OrderID(response.Result.OrderID),
		ClientOrderID: domain.ClientOrderID(response.Result.OrderLinkID),
	}, nil
}

func (api *BybitAPI) CreateFuturesOrder(
	ctx context.Context, request domain.OrderRequest,
) (domain.BybitOrderResponseAck, error) {
	// Only difference between CreateFuturesOrder and CreateSpotOrder is in parameter "spot" here
	response, err := api.createOrder(ctx, request, false)
	if err != nil {
		return domain.BybitOrderResponseAck{}, err
	}
	return domain.BybitOrderResponseAck{
		OrderID:       domain.OrderID(response.Result.OrderID),
		ClientOrderID: domain.ClientOrderID(response.Result.OrderLinkID),
	}, nil
}

func (api *BybitAPI) amendOrder(
	ctx context.Context,
	instrument bybit.InstrumentType,
	symbol string,
	orderID, clientOrderID string,
	newSize, newPrice decimal.Decimal,
) error {
	logger := api.logger.WithField("op", "clients.BybitAPI.amendOrder")
	inBodyParams := bybit.AmendOrderParams{
		Category:    instrument.GetString(),
		Symbol:      symbol,
		OrderID:     orderID,
		OrderLinkID: clientOrderID,
		Qty:         newSize.String(),
		Price:       newPrice.String(),
	}
	var response bybit.AmendOrderResponse
	_, err := api.request(
		ctx, http.MethodPost, bybit.AmendOrder, "", &inBodyParams, &response, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error amending order")
		return err
	}
	return nil
}

func (api *BybitAPI) AmendSpotOrder(
	ctx context.Context,
	symbol domain.Symbol,
	orderID domain.OrderID,
	clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	return api.amendOrder(
		ctx, bybit.SPOT, string(symbol), string(orderID), string(clientOrderID), newSize, newPrice,
	)
}

func (api *BybitAPI) AmendFuturesOrder(
	ctx context.Context,
	symbol domain.Symbol,
	orderID domain.OrderID,
	clientOrderID domain.ClientOrderID,
	newSize, newPrice decimal.Decimal,
) error {
	return api.amendOrder(
		ctx, bybit.LINEAR, bybit.SymbolFromDomain(symbol), string(orderID), string(clientOrderID), newSize, newPrice,
	)
}

func (api *BybitAPI) cancelOrder(
	ctx context.Context, symbol string, orderID, clientOrderID string,
) (*bybit.CancelOrderResponse, error) {
	logger := api.logger.WithField("op", "clients.BybitAPI.cancelOrder")
	inBodyParams := bybit.CancelOrderParams{
		Category:    bybit.SPOT.GetString(),
		Symbol:      symbol,
		OrderLinkID: clientOrderID,
		OrderID:     orderID,
	}
	var output bybit.CancelOrderResponse
	response, err := api.request(
		ctx, http.MethodPost, bybit.CancelOrder, "", &inBodyParams, &output, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error cancelling order")
		return nil, err
	}
	return response.(*bybit.CancelOrderResponse), nil
}

// CancelOrder cancels only spot order.
func (api *BybitAPI) CancelOrder(
	ctx context.Context, symbol domain.Symbol, orderID domain.OrderID, clientOrderID domain.ClientOrderID,
) (domain.BybitOrderResponseAck, error) {
	response, err := api.cancelOrder(ctx, string(symbol), string(orderID), string(clientOrderID))
	if err != nil {
		return domain.BybitOrderResponseAck{}, err
	}
	return domain.BybitOrderResponseAck{
		OrderID:       domain.OrderID(response.Result.OrderID),
		ClientOrderID: domain.ClientOrderID(response.Result.OrderLinkID),
	}, nil
}

func (api *BybitAPI) getOpenOrders(
	ctx context.Context, instrumentType bybit.InstrumentType, symbol string,
) (*bybit.GetOpenOrdersResponse, error) {
	if instrumentType == "" || symbol == "" {
		return nil, bybit.ErrInvalidParameters
	}
	logger := api.logger.WithField("op", "clients.BybitAPI.getOpenOrders")
	params := fmt.Sprintf("category=%v&symbol=%v", instrumentType.GetString(), symbol)
	var output bybit.GetOpenOrdersResponse
	response, err := api.request(
		ctx, http.MethodGet, bybit.GetOpenedOrders, params, nil, &output, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error getting open orders")
		return nil, err
	} else if len(response.(*bybit.GetOpenOrdersResponse).Result.List) == 0 {
		return nil, bybit.ErrNoOpenedOrdersFound
	}
	return response.(*bybit.GetOpenOrdersResponse), nil
}

func (api *BybitAPI) GetSpotOpenOrders(ctx context.Context, symbol domain.Symbol) ([]domain.OpenOrder, error) {
	response, err := api.getOpenOrders(ctx, bybit.SPOT, string(symbol))
	if err != nil {
		return nil, err
	}
	return bybit.OpenOrdersToDomain(*response), nil
}

func (api *BybitAPI) GetFuturesOpenOrders(
	ctx context.Context,
	symbol domain.Symbol,
) ([]domain.OpenOrder, error) {
	response, err := api.getOpenOrders(ctx, bybit.LINEAR, string(symbol))
	if err != nil {
		return nil, err
	}
	return bybit.OpenOrdersToDomain(*response), nil
}

func (api *BybitAPI) cancelAllOrders(
	ctx context.Context,
	symbol string,
	instrument bybit.InstrumentType,
) error {
	logger := api.logger.WithField("op", "clients.BybitAPI.cancelAllOrders")
	inBodyParams := &bybit.CancelAllOrdersParams{
		Category: instrument,
		Symbol:   symbol,
	}
	var output bybit.CancelAllOrdersResponse
	_, err := api.request(
		ctx, http.MethodPost, bybit.CancelAllOrders, "", inBodyParams, &output, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error cancelling all orders")
		return err
	}
	return nil
}

func (api *BybitAPI) CancelSpotOrders(
	ctx context.Context,
	symbol domain.Symbol,
) error {
	err := api.cancelAllOrders(ctx, bybit.SymbolFromDomain(symbol), bybit.SPOT)
	if err != nil {
		return err
	}
	return nil
}

func (api *BybitAPI) CancelFuturesOrders(
	ctx context.Context,
	symbol domain.Symbol,
) error {
	err := api.cancelAllOrders(ctx, bybit.SymbolFromDomain(symbol), bybit.LINEAR)
	if err != nil {
		return err
	}
	return nil
}

func (api *BybitAPI) GetBalances(ctx context.Context) ([]domain.Balance, error) {
	response, err := api.getBalances(ctx)
	if err != nil {
		return nil, err
	}
	return bybit.BalancesToDomain(*response), nil
}

func (api *BybitAPI) getBalances(ctx context.Context) (*bybit.GetBalancesResponse, error) {
	logger := api.logger.WithField("op", "clients.BybitAPI.getBalances")

	params := fmt.Sprintf("accountType=%s", bybit.UnifiedAccount)
	var output bybit.GetBalancesResponse
	response, err := api.request(
		ctx, http.MethodGet, bybit.GetBalances, params, nil, &output, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error getting balances")
		return nil, err
	}
	return response.(*bybit.GetBalancesResponse), nil
}

func (api *BybitAPI) GetSpotOrderHistory(
	ctx context.Context,
	startTime, endTime time.Time,
) ([]domain.Order, error) {
	response, err := api.getOrderHistory(ctx, bybit.SPOT, startTime, endTime)
	if err != nil {
		return nil, err
	}
	return bybit.OrderHistoryToDomain(*response), nil
}

func (api *BybitAPI) GetFuturesOrderHistory(
	ctx context.Context,
	startTime, endTime time.Time,
) ([]domain.Order, error) {
	response, err := api.getOrderHistory(ctx, bybit.LINEAR, startTime, endTime)
	if err != nil {
		return nil, err
	}
	return bybit.OrderHistoryToDomain(*response), nil
}

func (api *BybitAPI) getOrderHistory(
	ctx context.Context, instrumentType bybit.InstrumentType, startTime, endTime time.Time,
) (*bybit.OrderHistoryResponse, error) {
	logger := api.logger.WithField("op", "clients.BybitAPI.getOrderHistory")
	params := fmt.Sprintf(
		"category=%v&startTime=%v&endTime=%v", instrumentType, startTime.UnixMilli(), endTime.UnixMilli(),
	)
	var output bybit.OrderHistoryResponse
	response, err := api.request(
		ctx, http.MethodGet, bybit.GetOrderHistory, params, nil, &output, true,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error getting order history")
		return nil, err
	}
	return response.(*bybit.OrderHistoryResponse), nil
}

func (api *BybitAPI) GetOrderBook(
	ctx context.Context, symbol domain.Symbol, limit int,
) (domain.OrderBook, error) {
	response, err := api.getOrderBook(ctx, true, string(symbol), limit)
	if err != nil {
		return domain.OrderBook{}, err
	} else if len(response.Result.A) == 0 || len(response.Result.B) == 0 {
		return domain.OrderBook{}, bybit.ErrNoOrderBookFound
	}
	return bybit.OrderBookToDomain(response)
}

func (api *BybitAPI) getOrderBook(
	ctx context.Context,
	isSpot bool,
	symbol string,
	limit int,
) (*bybit.OrderBookResponse, error) {
	var instrument bybit.InstrumentType
	if isSpot {
		instrument = bybit.SPOT
	} else {
		instrument = bybit.LINEAR
	}
	logger := api.logger.WithField("op", "clients.BybitAPI.getOrderBook")
	request := bybit.OrderBookRequest{
		Category: instrument,
		Symbol:   symbol,
		Limit:    limit,
	}
	params, err := query.Values(request)
	if err != nil {
		logger.WithError(err).Errorf("Error getting order book")
		return nil, err
	}
	var output bybit.OrderBookResponse
	response, err := api.request(
		ctx, http.MethodGet, bybit.GetOrderBook, params.Encode(), nil, &output, false,
	)
	if err != nil {
		logger.WithError(err).Errorf("Error getting order book")
		return nil, err
	}
	return response.(*bybit.OrderBookResponse), nil
}

func (api *BybitAPI) generateSignature(data string) string {
	mac := hmac.New(sha256.New, []byte(api.apiSecret))
	mac.Write([]byte(data))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func (api *BybitAPI) setPrivateHeader(req *http.Request, signatureBody string, timestamp int64) {
	signature := api.generateSignature(signatureBody)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-API-KEY", api.apiKey)
	req.Header.Set("X-BAPI-TIMESTAMP", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-BAPI-RECV-WINDOW", fmt.Sprintf("%d", receiveWindow))
	req.Header.Set("Content-Type", "application/json")
}
