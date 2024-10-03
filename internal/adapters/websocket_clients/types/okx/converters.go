package okx

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
)

const unknown = "unknown"

func OrderBookToDomain(book OrderBookUpdate) domain.OKXOrderBook {
	data := book.Data[0]

	bids := make(map[string]float64)
	asks := make(map[string]float64)

	// Convert bids
	for _, bid := range data.Bids {
		price := bid[0]
		qty, err := strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return domain.OKXOrderBook{} //exhaustruct:ignore
		}
		bids[price] = qty
	}

	// Convert asks
	for _, ask := range data.Asks {
		price := ask[0]
		qty, err := strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return domain.OKXOrderBook{} //exhaustruct:ignore
		}
		asks[price] = qty
	}

	ts, _ := strconv.ParseInt(data.TS, 10, 64)

	return domain.OKXOrderBook{
		Type:      UpdateTypeToDomain(book.Action),
		Symbol:    domain.OKXSpotInstrument(book.Arg.InstID),
		Timestamp: time.UnixMilli(ts),
		Bids:      bids,
		Asks:      asks,
		Checksum:  data.Checksum,
	}
}

func UpdateTypeToDomain(updateType EventType) domain.UpdateType {
	switch updateType { //nolint:exhaustive // we only care about order book updates here
	case "snapshot":
		return domain.Snapshot
	case "update":
		return domain.Update
	default:
		return domain.Unknown
	}
}

func OrderToDomain(o OrderUpdate) domain.OKXUserStreamOrderUpdate {
	if len(o.Data) == 0 {
		return domain.OKXUserStreamOrderUpdate{} //exhaustruct:ignore
	}

	data := o.Data[0] // TODO Check

	price, _ := decimal.NewFromString(data.Px)
	origQty, _ := decimal.NewFromString(data.Sz)
	executedQty, _ := decimal.NewFromString(data.AccFillSz)
	feeAmount, _ := decimal.NewFromString(data.FillFee)

	fillPrice, _ := decimal.NewFromString(data.FillPx)
	fillSize, _ := decimal.NewFromString(data.FillSz)

	orderType, timeInForce := OrderTypeToDomain(data.OrdType)
	orderID, _ := strconv.ParseInt(data.OrdID, 10, 64)
	tradeID, _ := strconv.ParseInt(data.TradeID, 10, 64)

	createdTime, _ := strconv.ParseInt(data.CTime, 10, 64)
	updatedTime, _ := strconv.ParseInt(data.UTime, 10, 64)
	fillTime, _ := strconv.ParseInt(data.FillTime, 10, 64)

	orderStatus := StatusToDomain(data.State)

	switch orderStatus { //nolint:exhaustive // we only need to add some data for fill struct
	case domain.OrderStatusFilled, domain.OrderStatusPartiallyFilled:
		price, _ = decimal.NewFromString(data.FillPx)
	}
	amendResult := true
	if data.AmendResult == "-1" {
		amendResult = false
	}

	return domain.OKXUserStreamOrderUpdate{
		Instrument:        domain.OKXSpotInstrument(data.InstID),
		OrderID:           domain.OrderID(strconv.FormatInt(orderID, 10)),
		TradeID:           domain.OrderID(strconv.FormatInt(tradeID, 10)),
		ClientOrderID:     domain.ClientOrderID(data.ClOrdID),
		RequestID:         domain.ClientOrderID(data.ReqID),
		AmendResult:       amendResult,
		Price:             price,
		OrigQty:           origQty,
		ExecutedQty:       executedQty,
		FillSize:          fillSize,
		FillPrice:         fillPrice,
		Status:            orderStatus,
		TimeInForce:       timeInForce,
		Type:              orderType,
		Side:              SideToDomain(data.Side),
		TransactionTime:   time.UnixMilli(updatedTime),
		OrderCreationTime: time.UnixMilli(createdTime),
		FillTime:          time.UnixMilli(fillTime),
		FeeAmount:         feeAmount,
		FeeCurrency:       domain.Currency(data.FeeCcy),
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

func OrderTypeFromDomain(t domain.OrderType) string {
	switch t {
	case domain.OrderTypeLimit:
		return "limit"
	case domain.OrderTypeMarket:
		return "market"
	case domain.OrderTypeLimitMaker,
		domain.OrderTypeLimitFOK,
		domain.OrderTypeLimitIOC,
		domain.OrderTypeUnknown:
		return unknown
	default:
		return unknown
	}
}

func OrderAckToDomain(ack OrderAck) domain.OKXOrderResponseAck {
	if len(ack.Data) > 0 {
		orderID, _ := strconv.ParseInt(ack.Data[0].OrdID, 10, 64)
		return domain.OKXOrderResponseAck{
			ClientOrderID: domain.ClientOrderID(ack.Data[0].ClOrdID),
			OrderID:       domain.OrderID(strconv.FormatInt(orderID, 10)),
		}
	}
	return domain.OKXOrderResponseAck{} //nolint:exhaustruct // OK
}
