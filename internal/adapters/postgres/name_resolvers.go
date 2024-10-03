package postgres

import (
	"context"

	"studentgit.kata.academy/quant/torque/internal/adapters/postgres/queries"
	"studentgit.kata.academy/quant/torque/internal/domain"
	"studentgit.kata.academy/quant/torque/pkg/fp/slices"
)

func (db *Database) SelectAllNameResolvers(ctx context.Context) ([]domain.Instrument, error) {
	result, err := db.q.SelectAllNameResolvers(ctx)
	if err != nil {
		return nil, err
	}
	return slices.Map(result, func(row queries.QuantNameResolver) domain.Instrument {
		return domain.Instrument{
			Exchange:       domain.Exchange(row.Exchange),
			LocalSymbol:    domain.Symbol(row.LocalSymbol),
			ExchangeSymbol: domain.Symbol(row.ExchangeSymbol),
		}
	}), nil
}

func (db *Database) InsertNameResolvers(ctx context.Context, instrument []domain.Instrument) error {
	err := queries.BatchExec(db.q.InsertNameResolvers(
		ctx, slices.Map(instrument, func(instrument domain.Instrument,
		) queries.InsertNameResolversParams {
			return queries.InsertNameResolversParams{
				Exchange:       string(instrument.Exchange),
				LocalSymbol:    string(instrument.LocalSymbol),
				ExchangeSymbol: string(instrument.ExchangeSymbol),
			}
		})))
	if err != nil {
		return err
	}
	return nil
}
