package cache

import (
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/cache"
)

type InstrumentCache struct {
	ExchangeInstrument *cache.Cache[domain.LocalInstrument, domain.Instrument]
	LocalInstrument    *cache.Cache[domain.ExchangeInstrument, domain.Instrument]
}

func NewInstrumentCache() *InstrumentCache {
	return &InstrumentCache{
		ExchangeInstrument: cache.NewCache(func(instrument domain.Instrument) domain.LocalInstrument {
			return domain.LocalInstrument{
				Exchange: instrument.Exchange,
				Symbol:   instrument.LocalSymbol,
			}
		}),
		LocalInstrument: cache.NewCache(func(instrument domain.Instrument) domain.ExchangeInstrument {
			return domain.ExchangeInstrument{
				Exchange: instrument.Exchange,
				Symbol:   instrument.ExchangeSymbol,
			}
		}),
	}
}

func (i *InstrumentCache) Add(instrument domain.Instrument) {
	i.ExchangeInstrument.Add(instrument)
	i.LocalInstrument.Add(instrument)
}

func (i *InstrumentCache) GetLocalInstrument(instrument domain.ExchangeInstrument) (domain.Instrument, bool) {
	return i.LocalInstrument.Get(instrument)
}

func (i *InstrumentCache) GetExchangeInstrument(instrument domain.LocalInstrument) (domain.Instrument, bool) {
	return i.ExchangeInstrument.Get(instrument)
}
