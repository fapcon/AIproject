package cache

import (
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/cache"
)

type InMemoryStorage struct {
	InstrumentCache     *InstrumentCache
	OrderBook           *cache.Deletable[domain.LocalInstrument, domain.OrderBook]
	InstrumentDetails   *cache.Deletable[domain.LocalInstrument, domain.InstrumentDetails]
	PistonIDCache       *InverseCache[domain.PistonIDKey, domain.ClientOrderID]
	ActiveCommandsCache *cache.Deletable[domain.ClientOrderID, domain.ActiveCommand]
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		InstrumentCache: NewInstrumentCache(),
		OrderBook: cache.NewDeletable(func(b domain.OrderBook) domain.LocalInstrument {
			return domain.LocalInstrument{
				Exchange: b.Exchange,
				Symbol:   b.Instrument,
			}
		}),
		InstrumentDetails: cache.NewDeletable(func(b domain.InstrumentDetails) domain.LocalInstrument {
			return domain.LocalInstrument{
				Exchange: b.Exchange,
				Symbol:   b.LocalSymbol,
			}
		}),
		PistonIDCache: NewInverseCache[domain.PistonIDKey, domain.ClientOrderID](),
		ActiveCommandsCache: cache.NewDeletable(func(cmd domain.ActiveCommand) domain.ClientOrderID {
			return cmd.GetID()
		}),
	}
}
