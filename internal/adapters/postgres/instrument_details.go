package postgres

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/adapters/postgres/queries"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
)

func (db *Database) SelectAllInstrumentDetails(ctx context.Context) ([]domain.InstrumentDetails, error) {
	result, err := db.q.SelectAllInstrumentDetails(ctx)
	if err != nil {
		return nil, err
	}
	return slices.Map(result, func(row queries.QuantInstrumentDetail) domain.InstrumentDetails {
		return domain.InstrumentDetails{
			Exchange:       domain.Exchange(row.Exchange),
			LocalSymbol:    domain.Symbol(row.LocalSymbol),
			ExchangeSymbol: domain.Symbol(row.ExchangeSymbol),
			PricePrecision: int(row.PricePrecision),
			SizePrecision:  int(row.SizePrecision),
			MinLot:         row.MinLot,
		}
	}), nil
}

func (db *Database) InsertInstrumentDetails(ctx context.Context, instrument []domain.InstrumentDetails) error {
	err := queries.BatchExec(db.q.InsertInstrumentsDetails(
		ctx, slices.Map(instrument, func(instrument domain.InstrumentDetails,
		) queries.InsertInstrumentsDetailsParams {
			return queries.InsertInstrumentsDetailsParams{
				Exchange:       string(instrument.Exchange),
				LocalSymbol:    string(instrument.LocalSymbol),
				ExchangeSymbol: string(instrument.ExchangeSymbol),
				PricePrecision: int32(instrument.PricePrecision),
				SizePrecision:  int32(instrument.SizePrecision),
				MinLot:         instrument.MinLot,
			}
		})))
	if err != nil {
		return err
	}
	return nil
}
