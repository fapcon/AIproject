package okx

import (
	"fmt"
	"strconv"
	"time"

	"studentgit.kata.academy/quant/torque/internal/domain"

	"github.com/shopspring/decimal"
)

const (
	unknown = "unknown"
)

func OrderBookSnapshotToDomain(response OrderBookResponse) (domain.OKXOrderBook, error) {
	if len(response.Data) == 0 {
		return domain.OKXOrderBook{}, fmt.Errorf("no data available")
	}

	data := response.Data[0]

	ts, err := strconv.ParseInt(data.TS, 10, 64)
	if err != nil {
		return domain.OKXOrderBook{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	timestamp := time.UnixMilli(ts)

	bids := make(map[string]float64)
	asks := make(map[string]float64)

	// Convert bids
	for _, bid := range data.Bids {
		price := bid[0]
		qty, errX := strconv.ParseFloat(bid[1], 64)
		if errX != nil {
			return domain.OKXOrderBook{}, fmt.Errorf("failed to parse bid quantity: %w", err)
		}
		bids[price] = qty
	}

	// Convert asks
	for _, ask := range data.Asks {
		price := ask[0]
		qty, errX := strconv.ParseFloat(ask[1], 64)
		if errX != nil {
			return domain.OKXOrderBook{}, fmt.Errorf("failed to parse ask quantity: %w", err)
		}
		asks[price] = qty
	}

	return domain.OKXOrderBook{
		Symbol:    domain.OKXSpotInstrument(""),
		Type:      domain.Snapshot,
		Timestamp: timestamp,
		Bids:      bids,
		Asks:      asks,
		Checksum:  0,
	}, nil
}

func BalancesToDomain(response BalancesResponse) []domain.Balance {
	balances := make([]domain.Balance, 0, len(response.Data))
	for _, obj := range response.Data {
		for _, balance := range obj.Details {
			available, err := decimal.NewFromString(balance.AvailBal)
			if err != nil {
				continue
			}
			if available.IsZero() {
				continue
			}
			balances = append(balances, domain.Balance{
				Symbol:    domain.Symbol(balance.Ccy),
				Available: available,
			})
		}
	}
	return balances
}

func OpenOrdersToDomain(response OpenOrdersResponse) []domain.OpenOrder {
	var openOrders []domain.OpenOrder
	for _, order := range response.Data {
		cr, _ := strconv.ParseInt(order.CTime, 10, 64)
		ut, _ := strconv.ParseInt(order.UTime, 10, 64)
		created := time.UnixMilli(cr)
		updatedAt := time.UnixMilli(ut)
		price, _ := decimal.NewFromString(order.Px)
		origQty, _ := decimal.NewFromString(order.Sz)
		executedQty, _ := decimal.NewFromString(order.AccFillSz)
		orderType, tif := OrderTypeToDomain(order.OrdType)
		orderID, _ := strconv.ParseInt(order.OrdID, 10, 64)
		openOrder := domain.OpenOrder{
			ClientOrderID: domain.ClientOrderID(order.ClOrdID),
			OrderID:       domain.OrderID(strconv.FormatInt(orderID, 10)),
			Side:          SideToDomain(order.Side),
			Instrument:    domain.Symbol(order.InstID),
			Type:          orderType,
			TimeInForce:   tif,
			OrigQty:       origQty,
			Price:         price,
			ExecutedQty:   executedQty,
			Created:       created,
			UpdatedAt:     updatedAt,
			//Status:        StatusToDomain(order.State),
		}
		openOrders = append(openOrders, openOrder)
	}
	return openOrders
}

func OrderHistoryToDomain(response []HistoricalOrder) []domain.Order {
	orders := make([]domain.Order, 0, len(response))

	for _, o := range response {
		filled, _ := decimal.NewFromString(o.AccFillSz)
		fee, _ := decimal.NewFromString(o.Fee)

		orderType, _ := OrderTypeToDomain(o.OrdType)

		price, _ := decimal.NewFromString(o.FillPx)
		size, _ := decimal.NewFromString(o.FillSz)
		fillTime, _ := strconv.ParseInt(o.FillTime, 10, 64)

		cr, _ := strconv.ParseInt(o.CTime, 10, 64)
		orderID, _ := strconv.ParseInt(o.OrdID, 10, 64)
		tradeID, _ := strconv.ParseInt(o.TradeID, 10, 64)
		order := domain.Order{
			ClientOrderRequestID:  "", // we do not have this info for now, we fill it in the app
			PistonOrderID:         "",
			ExchangeAccount:       "",
			ExchangeClientOrderID: domain.ClientOrderID(o.ClOrdID),
			OrderID:               domain.OrderID(strconv.FormatInt(orderID, 10)),
			TradeID:               domain.OrderID(strconv.FormatInt(tradeID, 10)),
			Exchange:              domain.OKXExchange,
			Side:                  SideToDomain(o.Side),
			Instrument:            domain.Symbol(o.InstID),
			Type:                  orderType,
			Size:                  size, // In quote currency for market orders
			Price:                 price,
			RemainingSize:         size.Sub(filled),
			Filled:                filled,
			Volume:                size.Mul(price),
			Fee:                   fee,
			FeeCurrency:           domain.Currency(o.FeeCcy),
			Created:               time.UnixMilli(cr),
			UpdatedAt:             time.UnixMilli(fillTime),
			Status:                StatusToDomain(o.State),
		}

		orders = append(orders, order)
	}

	return orders
}

func FillsHistoryToDomain(response []HistoricalFill) []domain.Order {
	orders := make([]domain.Order, 0, len(response))

	for _, o := range response {
		filled, _ := decimal.NewFromString(o.FillSz)
		fee, _ := decimal.NewFromString(o.Fee)

		price, _ := decimal.NewFromString(o.FillPx)
		size, _ := decimal.NewFromString(o.FillSz)
		fillTime, _ := strconv.ParseInt(o.FillTime, 10, 64)

		orderID, _ := strconv.ParseInt(o.OrdID, 10, 64)
		tradeID, _ := strconv.ParseInt(o.TradeID, 10, 64)
		order := domain.Order{
			ClientOrderRequestID:  "",
			PistonOrderID:         "",
			ExchangeAccount:       "",
			Type:                  "",
			UpdatedAt:             time.Time{},
			Status:                domain.OrderStatusFilled,
			ExchangeClientOrderID: domain.ClientOrderID(o.ClOrdID),
			OrderID:               domain.OrderID(strconv.FormatInt(orderID, 10)),
			TradeID:               domain.OrderID(strconv.FormatInt(tradeID, 10)),
			Exchange:              domain.OKXExchange,
			Side:                  SideToDomain(o.Side),
			Instrument:            domain.Symbol(o.InstID),
			Size:                  size, // In quote currency for market orders
			Price:                 price,
			RemainingSize:         size.Sub(filled),
			Filled:                filled,
			Volume:                size.Mul(price),
			Fee:                   fee,
			FeeCurrency:           domain.Currency(o.FeeCcy),
			Created:               time.UnixMilli(fillTime),
		}

		orders = append(orders, order)
	}

	return orders
}

func CreateOrderToDomain(response CreateOrderResponse) domain.OKXOrderResponseAck {
	orderID, err := strconv.ParseInt(response.Data[0].OrdID, 10, 64)
	if err != nil {
		orderID = 0
	}
	return domain.OKXOrderResponseAck{
		ClientOrderID: domain.ClientOrderID(response.Data[0].ClOrdID),
		OrderID:       domain.OrderID(strconv.FormatInt(orderID, 10)),
	}
}

func CancelOrderToDomain(response CancelOrderResponse) domain.OKXOrderResponseAck {
	orderID, err := strconv.ParseInt(response.Data[0].OrdID, 10, 64)
	if err != nil {
		orderID = 0
	}
	return domain.OKXOrderResponseAck{
		ClientOrderID: domain.ClientOrderID(response.Data[0].ClOrdID),
		OrderID:       domain.OrderID(strconv.FormatInt(orderID, 10)),
	}
}

func SideFromDomain(s domain.Side) string {
	switch s {
	case domain.SideBuy:
		return "buy"
	case domain.SideSell:
		return "sell"
	case domain.SideUnknown:
		return unknown
	default:
		return unknown
	}
}

func SideToDomain(s string) domain.Side {
	switch s {
	case "buy":
		return domain.SideBuy
	case "sell":
		return domain.SideSell
	default:
		return domain.SideUnknown
	}
}

func OrderTypeFromDomain(t domain.OrderType) string {
	switch t { //nolint:exhaustive // we only support this orders for now
	case domain.OrderTypeLimit:
		return "limit"
	case domain.OrderTypeMarket:
		return "market"
	case domain.OrderTypeUnknown:
		return unknown
	default:
		return unknown
	}
}

func OrderTypeToDomain(t string) (domain.OrderType, domain.TimeInForce) {
	switch t {
	case "limit":
		return domain.OrderTypeLimit, domain.TimeInForceGTC
	case "market":
		return domain.OrderTypeMarket, ""
	case "fok":
		return domain.OrderTypeLimit, domain.TimeInForceFOK
	case "ioc":
		return domain.OrderTypeLimit, domain.TimeInForceIOC
	default:
		return domain.OrderTypeUnknown, domain.TimeInForceUnknown
	}
}

func StatusToDomain(s string) domain.OrderStatus {
	switch s {
	case "canceled":
		return domain.OrderStatusCanceled
	case "live":
		return domain.OrderStatusCreated
	case "partially_filled":
		return domain.OrderStatusPartiallyFilled
	case "filled":
		return domain.OrderStatusFilled
	default:
		return domain.OrderStatusUnknown
	}
}
