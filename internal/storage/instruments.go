package storage

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/adapters/kafka"

	"studentgit.kata.academy/quant/torque/internal/domain"
)

func (s *Storage) GetLocalInstrument(_ context.Context, i domain.ExchangeInstrument) (domain.LocalInstrument, bool) {
	s.mustHave(CapabilityNameResolver)
	v, ok := s.cache.InstrumentCache.GetLocalInstrument(i)
	if ok {
		return domain.LocalInstrument{
			Symbol:   v.LocalSymbol,
			Exchange: v.Exchange,
		}, true
	}
	return domain.LocalInstrument{}, false //nolint:exhaustruct // OK
}

func (s *Storage) GetExchangeInstrument(_ context.Context, i domain.LocalInstrument) (domain.ExchangeInstrument, bool) {
	s.mustHave(CapabilityNameResolver)
	v, ok := s.cache.InstrumentCache.GetExchangeInstrument(i)
	if ok {
		return domain.ExchangeInstrument{
			Symbol:   v.ExchangeSymbol,
			Exchange: v.Exchange,
		}, true
	}
	return domain.ExchangeInstrument{}, false //nolint:exhaustruct // OK
}

func (s *Storage) AddInstruments(ctx context.Context, instruments []domain.Instrument) error {
	s.mustHave(CapabilityNameResolver)
	err := s.postgres.InsertNameResolvers(ctx, instruments)
	if err != nil {
		return err
	}
	for _, i := range instruments {
		s.cache.InstrumentCache.Add(i)
	}
	return nil
}

func (s *Storage) GetAllInstrumentsDetailsByExchange(
	_ context.Context, exchange domain.Exchange,
) []domain.InstrumentDetails {
	s.mustHave(CapabilityInstrumentDetails)
	return s.cache.InstrumentDetails.Filter(func(detail domain.InstrumentDetails) bool {
		return detail.Exchange == exchange
	})
}

func (s *Storage) GetInstrumentsDetails(_ context.Context, i domain.LocalInstrument) (domain.InstrumentDetails, bool) {
	s.mustHave(CapabilityInstrumentDetails)
	return s.cache.InstrumentDetails.Get(i)
}

func (s *Storage) AddInstrumentDetails(ctx context.Context, instruments []domain.InstrumentDetails) error {
	s.mustHave(CapabilityInstrumentDetails)
	err := s.postgres.InsertInstrumentDetails(ctx, instruments)
	if err != nil {
		return err
	}
	s.cache.InstrumentDetails.Add(instruments...)
	return nil
}

func (b *Builder) WithInstrumentDetails() *Builder {
	b.buildCap(CapabilityInstrumentDetails)
	b.runCtx(func(s *Storage, ctx context.Context) {
		instrumentDetails, err := s.postgres.SelectAllInstrumentDetails(ctx)
		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Errorf("Storage.WithInstrumentDetails() error")
			return
		}
		s.cache.InstrumentDetails.Add(instrumentDetails...)
	})
	return b
}

func (b *Builder) WithNameResolver() *Builder {
	b.buildCap(CapabilityNameResolver)
	b.runCtx(func(s *Storage, ctx context.Context) {
		nameResolvers, err := s.postgres.SelectAllNameResolvers(ctx)
		if err != nil {
			s.logger.WithContext(ctx).WithError(err).Errorf("Storage.WithNameResolver() error")
			return
		}
		for _, nameResolver := range nameResolvers {
			s.cache.InstrumentCache.Add(nameResolver)
		}
	})
	return b
}

func (s *Storage) SubscribeInstrumentsRequest(_ context.Context, instruments []domain.Instrument) {
	s.mustHave(CapabilitySubscribeInstrumentsRequestsProducer)
	s.kafka.Producers.InternalEvents.SendSubscribeInstrumentsRequested(instruments)
}

func (b *Builder) WithSubscribeInstrumentsRequestProducer() *Builder {
	b.buildCap(CapabilitySubscribeInstrumentsRequestsProducer)
	b.buildKafka((*kafka.Kafka).RegisterInternalEventsProducer)
	return b
}
