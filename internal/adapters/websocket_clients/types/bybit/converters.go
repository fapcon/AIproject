package bybit

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
)

func OrderBookToDomain(book OrderBookUpdate) domain.BybitOrderBook {
	data := book.DataMD

	bids := make(map[string]float64, len(data.Bids))
	asks := make(map[string]float64, len(data.Asks))

	// Convert bids
	for _, bid := range data.Bids {
		price := bid[0]
		qty, _ := strconv.ParseFloat(bid[1], 64)
		bids[price] = qty
	}

	// Convert asks
	for _, ask := range data.Asks {
		price := ask[0]
		qty, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return domain.BybitOrderBook{} //exhaustruct:ignore
		}
		asks[price] = qty
	}

	return domain.BybitOrderBook{
		Type:      UpdateTypeToDomain(book.Type),
		Symbol:    domain.BybitSpotInstrument(book.DataMD.Symbol),
		Timestamp: time.UnixMilli(book.TS),
		Bids:      bids,
		Asks:      asks,
	}
}

func UpdateTypeToDomain(updateType OpType) domain.UpdateTypeBybit {
	switch updateType { //nolint:exhaustive // we only care about order book updates here
	case "snapshot":
		return domain.SnapshotBybit
	case "delta":
		return domain.DeltaBybit
	default:
		return domain.UnknownBybit
	}
}

func OrderToDomain(o OrderUpdate) []domain.BybitUserStreamOrderUpdate {
	if len(o.Data) == 0 {
		return []domain.BybitUserStreamOrderUpdate{} //exhaustruct:ignore
	}

	orders := make([]domain.BybitUserStreamOrderUpdate, 0, len(o.Data))

	for _, data := range o.Data {
		var executedQty, feeAmount, fillPrice, fillSize decimal.Decimal
		price, _ := decimal.NewFromString(data.Price)
		origQty, _ := decimal.NewFromString(data.Qty)

		createdTime, _ := strconv.ParseInt(data.CreatedTime, 10, 64)
		updatedTime, _ := strconv.ParseInt(data.UpdatedTime, 10, 64)

		orderStatus := StatusToDomain(data.OrderStatus)

		//nolint:exhaustive // we only need to add some data for fill struct
		switch orderStatus {
		case domain.OrderStatusFilled, domain.OrderStatusPartiallyFilled:
			fillPrice = price
			fillSize, _ = decimal.NewFromString(data.Qty)
			feeAmount, _ = decimal.NewFromString(data.CumExecFee)
			executedQty, _ = decimal.NewFromString(data.CumExecQty)
		}

		orders = append(orders, domain.BybitUserStreamOrderUpdate{
			Instrument:        domain.BybitSpotInstrument(data.Symbol),
			OrderID:           domain.OrderID(data.OrderID),
			ClientOrderID:     domain.ClientOrderID(data.OrderLinkID),
			Price:             price,
			OrigQty:           origQty,
			ExecutedQty:       executedQty,
			FillSize:          fillSize,
			FillPrice:         fillPrice,
			Status:            orderStatus,
			TimeInForce:       OrderTimeInForceToDomain(data.TimeInForce),
			Type:              OrderTypeToDomain(data.OrderType),
			Side:              SideToDomain(data.Side),
			TransactionTime:   time.UnixMilli(updatedTime),
			OrderCreationTime: time.UnixMilli(createdTime),
			FeeAmount:         feeAmount,
			FeeCurrency:       domain.Currency(data.FeeCurrency),
		})
	}
	return orders
}

func SideToDomain(s string) domain.Side {
	switch s {
	case "Buy":
		return domain.SideBuy
	case "Sell":
		return domain.SideSell
	default:
		return domain.SideUnknown
	}
}

func OrderTypeToDomain(t string) domain.OrderType {
	switch t {
	case "Limit":
		return domain.OrderTypeLimit
	case "Market":
		return domain.OrderTypeMarket
	default:
		return domain.OrderTypeUnknown
	}
}

func OrderTimeInForceToDomain(t string) domain.TimeInForce {
	switch t {
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

func StatusToDomain(s string) domain.OrderStatus {
	switch s {
	case "New":
		return domain.OrderStatus(s)
	case "PartiallyFilled":
		return domain.OrderStatusPartiallyFilled
	case "Rejected":
		return domain.OrderStatusRejected
	case "Filled":
		return domain.OrderStatusFilled
	case "Cancelled":
		return domain.OrderStatusCanceled
	default:
		return domain.OrderStatusUnknown
	}
}
