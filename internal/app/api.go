package app

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/internal/domain/f"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type apiStorage interface {
	OKXSubscribeInstruments(context.Context, []domain.OKXInstrument) error
	OKXUnsubscribeInstruments(context.Context, []domain.OKXInstrument) error

	AddInstruments(context.Context, []domain.Instrument) error
	AddInstrumentDetails(ctx context.Context, instruments []domain.InstrumentDetails) error

	SubscribeInstrumentsRequest(context.Context, []domain.Instrument)
}

type API struct {
	logger  logster.Logger
	storage apiStorage
}

func NewAPI(logger logster.Logger, storage apiStorage) *API {
	return &API{
		logger:  logger.WithField(f.Module, "api"),
		storage: storage,
	}
}

func (a *API) SubscribeInstrumentsRequest(ctx context.Context, instruments []domain.Instrument) {
	a.storage.SubscribeInstrumentsRequest(ctx, instruments)
}

func (a *API) SubscribeInstruments(ctx context.Context, instruments []domain.Instrument) error {
	exchanges := make(map[domain.Exchange][]domain.Instrument)
	for _, instrument := range instruments {
		exchanges[instrument.Exchange] = append(exchanges[instrument.Exchange], instrument)
	}
	for exchange, subscriptionList := range exchanges {
		if exchange == domain.OKXExchange {
			err := a.storage.OKXSubscribeInstruments(ctx, slices.Map(subscriptionList,
				func(instrument domain.Instrument) domain.OKXInstrument {
					return domain.OKXSpotInstrument(instrument.ExchangeSymbol) // TODO Modify this for futures in future
				}),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *API) UnsubscribeInstruments(ctx context.Context, instruments []domain.Instrument) error {
	exchanges := make(map[domain.Exchange][]domain.Instrument)
	for _, instrument := range instruments {
		exchanges[instrument.Exchange] = append(exchanges[instrument.Exchange], instrument)
	}
	for exchange, subscriptionList := range exchanges {
		if exchange == domain.OKXExchange {
			err := a.storage.OKXUnsubscribeInstruments(ctx, slices.Map(subscriptionList,
				func(instrument domain.Instrument) domain.OKXInstrument {
					return domain.OKXSpotInstrument(instrument.ExchangeSymbol) // TODO Modify this for futures in future
				}),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *API) AddInstruments(ctx context.Context, instruments []domain.Instrument) error {
	a.logger.WithField(f.Instrument, instruments).Infof("adding instruments")
	return a.storage.AddInstruments(ctx, instruments)
}

func (a *API) AddInstrumentDetails(ctx context.Context, instrumentDetails []domain.InstrumentDetails) error {
	a.logger.WithField(f.InstrumentDetails, instrumentDetails).Infof("adding instruments details")
	return a.storage.AddInstrumentDetails(ctx, instrumentDetails)
}
