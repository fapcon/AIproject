package bybit

import (
	"errors"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
)

func SideFromDomain(s domain.Side) string {
	switch s {
	case domain.SideBuy:
		return "Buy"
	case domain.SideSell:
		return "Sell"
	case domain.SideUnknown:
		return "Unknown"
	default:
		return ""
	}
}

func TimeInForceFromDomain(t domain.TimeInForce) TimeInForce {
	switch t {
	case domain.TimeInForceGTC:
		return TimeInForceGTC
	case domain.TimeInForceFOK:
		return TimeInForceFOK
	case domain.TimeInForceIOC:
		return TimeInForceIOC
	case domain.TimeInForceUnknown:
		return TimeInForceUnknown
	default:
		return ""
	}
}

func SymbolFromDomain(s domain.Symbol) string {
	return string(s)
}

func OrderTypeFromDomain(t domain.OrderType) string {
	switch t {
	case domain.OrderTypeLimit:
		return "limit"
	case domain.OrderTypeMarket:
		return "market"
	case domain.OrderTypeUnknown:
		return "unknown"
	case domain.OrderTypeLimitFOK:
		return "limit_fok"
	case domain.OrderTypeLimitIOC:
		return "limit_ioc"
	case domain.OrderTypeLimitMaker:
		return "limit_maker"
	default:
		return ""
	}
}

func OpenOrdersToDomain(response GetOpenOrdersResponse) []domain.OpenOrder {
	output := make([]domain.OpenOrder, 0, len(response.Result.List))

	for _, order := range response.Result.List {
		origQty, _ := decimal.NewFromString(order.Qty)
		price, _ := decimal.NewFromString(order.Price)
		executedQty, _ := decimal.NewFromString(order.CumExecQty)
		created, _ := timeFromString(order.CreatedTime)
		updatedAt, _ := timeFromString(order.CreatedTime)
		output = append(output, domain.OpenOrder{
			ClientOrderID: domain.ClientOrderID(order.OrderLinkID),
			OrderID:       domain.OrderID(order.OrderID),
			Side:          sideToDomain(order.Side),
			Instrument:    domain.Symbol(order.Symbol),
			Type:          orderTypeToDomain(order.OrderType),
			TimeInForce:   timeInForceToDomain(order.TimeInForce),
			OrigQty:       origQty,
			Price:         price,
			ExecutedQty:   executedQty,
			Created:       created,
			UpdatedAt:     updatedAt,
		})
	}
	return output
}

func sideToDomain(s string) domain.Side {
	switch s {
	case "Buy":
		return domain.SideBuy
	case "Sell":
		return domain.SideSell
	default:
		return domain.SideUnknown
	}
}

func orderTypeToDomain(s string) domain.OrderType {
	switch s {
	case "Limit":
		return domain.OrderTypeLimit
	case "Market":
		return domain.OrderTypeMarket
	default:
		return domain.OrderTypeUnknown
	}
}

func timeInForceToDomain(s string) domain.TimeInForce {
	switch s {
	case "GTC":
		return domain.TimeInForceGTC
	case "IOC":
		return domain.TimeInForceIOC
	case "FOK":
		return domain.TimeInForceFOK
	default:
		return domain.TimeInForceUnknown
	}
}

func BalancesToDomain(response GetBalancesResponse) []domain.Balance {
	balances := make([]domain.Balance, 0, len(response.Result.List))
	for _, balance := range response.Result.List[0].Coin {
		total, _ := decimal.NewFromString(balance.WalletBalance)
		balances = append(balances, domain.Balance{
			Symbol:    domain.Symbol(balance.Coin),
			Available: total,
		})
	}
	return balances
}

func OrderHistoryToDomain(response OrderHistoryResponse) []domain.Order {
	output := make([]domain.Order, 0, len(response.Result.List))
	for _, order := range response.Result.List {
		avgPrice, _ := strconv.ParseFloat(order.AvgPrice, 64)
		filledQty, _ := strconv.ParseFloat(order.CumExecQty, 64)

		size, _ := decimal.NewFromString(order.Qty)
		remainingSize, _ := decimal.NewFromString(order.LeavesQty)
		filled, _ := decimal.NewFromString(order.CumExecQty)
		price, _ := decimal.NewFromString(order.Price)
		volume := decimal.NewFromFloat(avgPrice * filledQty)
		fee, _ := decimal.NewFromString(order.CumExecFee)
		createdAt, _ := timeFromString(order.CreatedTime)
		updatedAt, _ := timeFromString(order.UpdatedTime)

		output = append(output, domain.Order{
			ClientOrderRequestID:  "",
			PistonOrderID:         "",
			ExchangeAccount:       "",
			ExchangeClientOrderID: domain.ClientOrderID(order.OrderLinkID),
			OrderID:               domain.OrderID(order.OrderID),
			TradeID:               domain.OrderID(order.BlockTradeID),
			Exchange:              domain.BybitExchange,
			Side:                  sideToDomain(order.Side),
			Instrument:            domain.Symbol(order.Symbol),
			Type:                  orderTypeToDomain(order.OrderType),
			Size:                  size,
			Price:                 price,
			RemainingSize:         remainingSize,
			Filled:                filled,
			Volume:                volume,
			Fee:                   fee,
			FeeCurrency:           "",
			Created:               createdAt,
			UpdatedAt:             updatedAt,
			Status:                domain.OrderStatus(order.OrderStatus),
		})
	}
	return output
}

func OrderBookToDomain(in *OrderBookResponse) (domain.OrderBook, error) {
	if len(in.Result.A) == 0 || len(in.Result.B) == 0 {
		return domain.OrderBook{}, errors.New("empty orderbook")
	}
	outAks := make(map[string]float64, len(in.Result.A))
	outBids := make(map[string]float64, len(in.Result.B))

	for _, ask := range in.Result.A {
		price := ask[0]
		qty, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return domain.OrderBook{}, errors.New("failed to parse ask qty")
		}
		outAks[price] = qty
	}
	for _, bid := range in.Result.B {
		price := bid[0]
		qty, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return domain.OrderBook{}, errors.New("failed to parse bid qty")
		}
		outBids[price] = qty
	}
	return domain.OrderBook{
		Exchange:   domain.BybitExchange,
		Instrument: domain.Symbol(in.Result.S),
		Asks:       outAks,
		Bids:       outBids,
		Checksum:   0,
	}, nil
}

func timeFromString(in string) (time.Time, error) {
	u, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(u), nil
}
