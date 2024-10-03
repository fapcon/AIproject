package piston

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
)

func MarketDataMessageToDomain(mdMsg *MarketdataMessage) domain.PistonSubscribeMarketDataMessage {
	var msg domain.PistonSubscribeMarketDataMessage
	switch msgType := mdMsg.GetMsgType(); msgType { //nolint:exhaustive // We don't need to process other events
	case MarketdataMsgType_SUBSCRIBE:
		subscribeMsg := mdMsg.GetSubscribeMessage()
		msg = domain.PistonSubscribeMarketDataMessage{
			Instruments: slices.Map(subscribeMsg.GetInstruments(), func(instrument string) domain.Symbol {
				return domain.Symbol(instrument)
			}),
		}
	default:
	}
	return msg
}

func BalanceToTradingMessage(b domain.PistonBalancesMessage) *TradingMessage {
	balances := make(map[string]float64)
	for _, balance := range b.Balances {
		availableFloat64, _ := balance.Available.Float64()
		balances[string(balance.Symbol)] = availableFloat64
	}
	return &TradingMessage{
		Type: TradingMessageType_ExchangeBalancesResponse,
		Data: &TradingMessage_Balances{
			Balances: &ExchangeBalances{
				Exchange:        string(b.Exchange),
				ExchangeAccount: string(b.ExchangeAccount),
				Balances:        balances,
				Timestamp:       b.Timestamp.UnixNano(),
			},
		},
	}
}

func OpenOrdersToTradingMessage(m domain.PistonOpenOrdersMessage) *TradingMessage {
	orders := m.Orders
	ordersMsg := make([]*OpenOrder, 0, len(orders))
	for _, order := range orders {
		price, _ := order.Price.Float64()
		remainingSize, _ := order.RemainingSize.Float64()
		size, _ := order.Size.Float64()

		ordersMsg = append(ordersMsg, &OpenOrder{
			PistonClientOrderID:   string(order.PistonID),
			ExchangeID:            string(order.ExchangeID),
			ExchangeClientOrderID: string(order.ExchangeClientOrderID),
			Exchange:              string(order.Exchange),
			Side:                  SideFromDomain(order.Side),
			Instrument:            string(order.Instrument),
			Type:                  OrderTypeFromDomain(order.OrderType),
			Size:                  size,
			Price:                 price,
			RemainingSize:         remainingSize,
			Created:               order.Created.UnixNano(),
			LastUpdated:           order.LastUpdated.UnixNano(),
		})
	}
	return &TradingMessage{
		Type: TradingMessageType_OpenOrdersResponse,
		Data: &TradingMessage_OpenOrders{
			OpenOrders: &OpenOrdersRequestResult{
				RequestID: m.RequestID,
				Orders:    ordersMsg,
			},
		},
	}
}

func CancelOrderTradingMessageToDomain(order *CancelOrder) domain.CancelOrder {
	return domain.CancelOrder{
		Instrument:            domain.Symbol(order.GetInstrument()),
		OrderRequestID:        domain.PistonOrderRequestID(order.GetOrderRequestID()),
		PistonID:              domain.PistonID(order.GetPistonID()),
		ExchangeID:            domain.OrderID(order.ExchangeID),
		ExchangeClientOrderID: domain.ClientOrderID(order.GetExchangeClientOrderID()),
		Side:                  domain.SideUnknown,
	}
}

func CreateOrderTradingMessageToDomain(order *AddOrder) domain.AddOrder {
	return domain.AddOrder{
		OrderRequestID: domain.PistonOrderRequestID(order.OrderRequestID),
		PistonID:       domain.PistonID(order.PistonID),
		Instrument:     domain.Symbol(order.Instrument),
		Side:           SideToDomain(order.Side),
		Price:          decimal.NewFromFloat(order.Price),
		Size:           decimal.NewFromFloat(order.Size),
		OrderType:      OrderTypeToDomain(order.Type),
		Created:        time.Unix(0, order.Created),
	}
}

func MoveOrderTradingMessageToDomain(order *MoveOrder) domain.MoveOrder {
	return domain.MoveOrder{
		OrderRequestID:        domain.PistonOrderRequestID(order.OrderRequestID),
		PistonID:              domain.PistonID(order.PistonID),
		ExchangeID:            domain.OrderID(order.ExchangeID),
		Exchange:              domain.Exchange(order.Exchange),
		ExchangeClientOrderID: domain.ClientOrderID(order.ExchangeClientOrderID),
		Instrument:            domain.Symbol(order.Instrument),
		Side:                  SideToDomain(order.Side),
		Price:                 decimal.NewFromFloat(order.Price),
		Size:                  decimal.NewFromFloat(order.Size),
		OrderType:             OrderTypeToDomain(order.Type),
		Created:               time.Unix(0, order.Created),
	}
}

func OrderBookToMarketDataMessage(orderBook domain.OrderBook) *MarketdataMessage {
	// TODO here we should have ready to send book, this should be done in BL
	//nolint:gomnd // Not ok, will fix
	topAsks := orderBook.TopLevels(false, 25) // false for asks
	//nolint:gomnd // Not ok, will fix
	topBids := orderBook.TopLevels(true, 25) // true for bids

	asks := convertTopLevelsToProto(topAsks)
	bids := convertTopLevelsToProto(topBids)

	pbOrderBook := &OrderBook{
		Instrument:     string(orderBook.Instrument),
		Exchange:       string(orderBook.Exchange),
		Timestamp:      time.Now().UnixNano(),
		LocalTimestamp: time.Now().UnixNano(),
		Bids:           bids,
		Asks:           asks,
	}

	// Create a pb.OrderBookSnapshot
	orderBookSnapshot := &OrderBookSnapshot{
		Snapshot: pbOrderBook,
	}
	bookMessage := &BookMessage{
		Book: orderBookSnapshot,
	}

	marketDataMessage := &MarketdataMessage{
		MsgType: MarketdataMsgType_BOOK,
		Message: &MarketdataMessage_BookMessage{BookMessage: bookMessage},
	}
	return marketDataMessage
}

func convertTopLevelsToProto(levels []domain.PriceLevel) []*LevelInfo {
	protoLevels := make([]*LevelInfo, len(levels))
	for i, level := range levels {
		price, _ := strconv.ParseFloat(level.Price, 64)
		protoLevels[i] = &LevelInfo{
			Price: price,
			Size:  level.Quantity,
		}
	}
	return protoLevels
}

func SideFromDomain(side domain.Side) OrderSide {
	switch side {
	case domain.SideBuy:
		return OrderSide_BID
	case domain.SideSell:
		return OrderSide_ASK
	case domain.SideUnknown:
		return OrderSide_UNKNOWN
	default:
		return OrderSide_UNKNOWN
	}
}

func SideToDomain(side OrderSide) domain.Side {
	switch side {
	case OrderSide_BID:
		return domain.SideBuy
	case OrderSide_ASK:
		return domain.SideSell
	case OrderSide_UNKNOWN:
		return domain.SideUnknown
	default:
		return domain.SideUnknown
	}
}

func OrderTypeFromDomain(orderType domain.OrderType) OrderType {
	switch orderType {
	case domain.OrderTypeLimit:
		return OrderType_LIMIT
	case domain.OrderTypeMarket:
		return OrderType_MARKET
	case domain.OrderTypeLimitMaker:
		return OrderType_POST_ONLY
	case domain.OrderTypeLimitIOC:
		return OrderType_LIMIT_IOC
	case domain.OrderTypeLimitFOK:
		return OrderType_LIMIT_FOK
	case domain.OrderTypeUnknown:
		return OrderType_UNKNOWN_TYPE
	default:
		return OrderType_UNKNOWN_TYPE
	}
}

func OrderTypeToDomain(orderType OrderType) domain.OrderType {
	//nolint:exhaustive // We do not support all of the order types here
	switch orderType {
	case OrderType_LIMIT:
		return domain.OrderTypeLimit
	case OrderType_MARKET:
		return domain.OrderTypeMarket
	case OrderType_POST_ONLY:
		return domain.OrderTypeLimitMaker
	case OrderType_LIMIT_FOK:
		return domain.OrderTypeLimitFOK
	case OrderType_LIMIT_IOC:
		return domain.OrderTypeLimitIOC
	default:
		return domain.OrderTypeUnknown
	}
}

func OrderAddedToTradingMessage(order domain.OrderAdded) *TradingMessage {
	size, _ := order.Size.Float64()
	price, _ := order.Price.Float64()
	return &TradingMessage{
		Type: TradingMessageType_OrderAddedResponse,
		Data: &TradingMessage_OrderAdded{
			OrderAdded: &OrderAdded{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				Exchange:              string(order.Exchange),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
				Instrument:            string(order.Instrument),
				Side:                  SideFromDomain(order.Side),
				Type:                  OrderTypeFromDomain(order.OrderType),
				Size:                  size,
				Price:                 price,
				LastUpdated:           order.LastUpdated.UnixNano(),
			},
		},
	}
}

func OrderAddRejectToTradingMessage(order domain.OrderAddReject) *TradingMessage {
	return &TradingMessage{
		Type: TradingMessageType_OrderAddRejectedResponse,
		Data: &TradingMessage_OrderAddRejected{
			OrderAddRejected: &OrderAddRejected{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
			},
		},
	}
}

func OrderMovedToTradingMessage(order domain.OrderMoved) *TradingMessage {
	size, _ := order.Size.Float64()
	price, _ := order.Price.Float64()
	return &TradingMessage{
		Type: TradingMessageType_OrderMovedResponse,
		Data: &TradingMessage_OrderMoved{
			OrderMoved: &OrderMoved{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				Exchange:              string(order.Exchange),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
				Size:                  size,
				Price:                 price,
				Created:               order.Created.UnixNano(),
			},
		},
	}
}

func OrderMoveRejectToTradingMessage(order domain.OrderMoveReject) *TradingMessage {
	return &TradingMessage{
		Type: TradingMessageType_OrderMoveRejectedResponse,
		Data: &TradingMessage_OrderMoveRejected{
			OrderMoveRejected: &OrderMoveRejected{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
			},
		},
	}
}

func OrderCancelledToTradingMessage(order domain.OrderCancelled) *TradingMessage {
	return &TradingMessage{
		Type: TradingMessageType_OrderCancelledResponse,
		Data: &TradingMessage_OrderCancelled{
			OrderCancelled: &OrderCancelled{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
			},
		},
	}
}

func OrderCancelRejectToTradingMessage(order domain.OrderCancelReject) *TradingMessage {
	return &TradingMessage{
		Type: TradingMessageType_OrderCancelRejectedResponse,
		Data: &TradingMessage_OrderCancelRejected{
			OrderCancelRejected: &OrderCancelRejected{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
			},
		},
	}
}

func OrderFilledToTradingMessage(order domain.OrderFilled) *TradingMessage {
	size, _ := order.Size.Float64()
	price, _ := order.Price.Float64()
	fee, _ := order.Fee.Float64()
	return &TradingMessage{
		Type: TradingMessageType_OrderFilledResponse,
		Data: &TradingMessage_OrderFilled{
			OrderFilled: &OrderFilled{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				Exchange:              string(order.Exchange),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
				TradeID:               string(order.TradeID),
				Instrument:            string(order.Instrument),
				Side:                  SideFromDomain(order.Side),
				Size:                  size,
				Price:                 price,
				Fee:                   fee,
				FeeCurrency:           string(order.FeeCurrency),
				Timestamp:             order.Timestamp.UnixNano(),
				Market:                "",
			},
		},
	}
}

func OrderExecutedToTradingMessage(order domain.OrderExecuted) *TradingMessage {
	return &TradingMessage{
		Type: TradingMessageType_OrderExecutedResponse,
		Data: &TradingMessage_OrderExecuted{
			OrderExecuted: &OrderExecuted{
				OrderRequestID:        string(order.OrderRequestID),
				PistonID:              string(order.PistonID),
				ExchangeID:            string(order.ExchangeID),
				ExchangeClientOrderID: string(order.ExchangeClientOrderID),
			},
		},
	}
}

func InstrumentDetailsToTradingMessage(details []domain.InstrumentDetails) *TradingMessage {
	return &TradingMessage{
		Type: TradingMessageType_InstrumentsDetailsResponse,
		Data: &TradingMessage_InstrumentsDetailsResponse{
			InstrumentsDetailsResponse: &InstrumentsDetails{
				InstrumentsDetails: slices.Map(details, func(d domain.InstrumentDetails) *InstrumentDetails {
					lot, _ := d.MinLot.Float64()
					return &InstrumentDetails{
						Instrument: string(d.LocalSymbol),
						MinLot:     lot,
					}
				}),
			},
		},
	}
}
