package api

import (
	"github.com/shopspring/decimal"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
)

func (x *Instruments) ToDomain() []domain.Instrument {
	return slices.Map(x.GetInstruments(), func(instrument *Instrument) domain.Instrument {
		return domain.Instrument{
			LocalSymbol:    domain.Symbol(instrument.GetLocalSymbol()),
			ExchangeSymbol: domain.Symbol(instrument.GetExchangeSymbol()),
			Exchange:       instrument.GetExchange().ToDomain(),
		}
	})
}

func (x *InstrumentDetails) ToDomain() []domain.InstrumentDetails {
	return slices.Map(x.GetInstrumentDetails(), func(instrument *InstrumentDetail) domain.InstrumentDetails {
		return domain.InstrumentDetails{
			LocalSymbol:    domain.Symbol(instrument.GetLocalSymbol()),
			ExchangeSymbol: domain.Symbol(instrument.GetExchangeSymbol()),
			Exchange:       instrument.GetExchange().ToDomain(),
			PricePrecision: int(instrument.GetPricePrecision()),
			SizePrecision:  int(instrument.GetSizePrecision()),
			MinLot:         decimal.NewFromFloat(float64(instrument.MinLot)),
		}
	})
}

func (x Exchange) ToDomain() domain.Exchange {
	return map[Exchange]domain.Exchange{
		Exchange_ExchangeUnknown: "",
		Exchange_ExchangeOKX:     domain.OKXExchange,
	}[x]
}
