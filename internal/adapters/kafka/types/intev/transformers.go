package intev

import (
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
)

func SubscribeInstrumentsRequestedFromDomain(
	d []domain.Instrument,
) *SubscribeInstruments {
	return &SubscribeInstruments{
		Instruments: slices.Map(
			d, func(d domain.Instrument) *Instrument {
				return &Instrument{
					Exchange:       ExchangeFromDomain(d.Exchange),
					LocalSymbol:    string(d.LocalSymbol),
					ExchangeSymbol: string(d.ExchangeSymbol),
				}
			},
		),
	}
}

func ExchangeFromDomain(d domain.Exchange) Exchange {
	return map[domain.Exchange]Exchange{
		domain.OKXExchange:     Exchange_ExchangeOKX,
		domain.GateIOExchange:  Exchange_ExchangeGateIO,
		domain.BybitExchange:   Exchange_ExchangeBybit,
		domain.UnknownExchange: Exchange_ExchangeUnknown,
	}[d]
}

func (x AckableNewOrderRequested) ToDomain() domain.Ackable[domain.NewOrderRequest] {
	return domain.NewAckable(x.NewOrderRequested.ToDomain(), x.Ack)
}

func (x *NewOrderRequested) ToDomain() domain.NewOrderRequest {
	return domain.NewOrderRequest{
		ClientOrderID:   domain.ClientOrderID(x.GetClientOrderId()),
		LocalInstrument: x.GetLocalInstrument().ToDomain(),
		Price:           x.GetPrice().AsDecimal(),
		Qty:             x.GetQty().AsDecimal(),
		Side:            x.GetOrderSide().ToDomain(),
		Type:            x.GetType().ToDomain(),
		TimeInForce:     x.GetTimeInForce().ToDomain(),
	}
}

func (ot OrderType) ToDomain() domain.OrderType {
	switch ot {
	case OrderType_ORDER_TYPE_MARKET:
		return domain.OrderTypeMarket
	case OrderType_ORDER_TYPE_LIMIT:
		return domain.OrderTypeLimit
	case OrderType_ORDER_TYPE_UNSPECIFIED:
		return domain.OrderTypeUnknown
	default:
		return domain.OrderTypeUnknown
	}
}

func (os OrderSide) ToDomain() domain.Side {
	switch os {
	case OrderSide_ORDER_SIDE_BUY:
		return domain.SideBuy
	case OrderSide_ORDER_SIDE_SELL:
		return domain.SideSell
	case OrderSide_ORDER_SIDE_UNSPECIFIED:
		return domain.SideUnknown
	default:
		return domain.SideUnknown
	}
}

func (tif TimeInForce) ToDomain() domain.TimeInForce {
	switch tif {
	case TimeInForce_TIME_IN_FORCE_GTC:
		return domain.TimeInForceGTC
	case TimeInForce_TIME_IN_FORCE_IOC:
		return domain.TimeInForceIOC
	case TimeInForce_TIME_IN_FORCE_FOK:
		return domain.TimeInForceFOK
	case TimeInForce_TIME_IN_FORCE_UNSPECIFIED:
		return domain.TimeInForceUnknown
	default:
		return domain.TimeInForceUnknown
	}
}

func (li *LocalInstrument) ToDomain() domain.LocalInstrument {
	return domain.LocalInstrument{
		Exchange: li.Exchange.ToDomain(),
		Symbol:   domain.Symbol(li.Symbol),
	}
}

func (e Exchange) ToDomain() domain.Exchange {
	switch e {
	case Exchange_ExchangeOKX:
		return domain.OKXExchange
	case Exchange_ExchangeBybit:
		return domain.BybitExchange
	case Exchange_ExchangeGateIO:
		return domain.GateIOExchange
	case Exchange_ExchangeUnknown:
		return domain.UnknownExchange
	default:
		return domain.UnknownExchange
	}
}
