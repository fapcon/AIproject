package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	"studentgit.kata.academy/quant/torque/internal/adapters/postgres/queries"
	"studentgit.kata.academy/quant/torque/pkg/postgres"
)

type dbOrTx interface {
	queries.DBTX
	Begin(context.Context) (pgx.Tx, error)
}

type Database struct {
	q  *queries.Queries
	db dbOrTx
}

func NewDatabase(db *postgres.DB) (*Database, error) {
	err := db.ValidateQueries([]string{
		queries.SelectAllNameResolvers,
		queries.InsertNameResolvers,
		queries.InsertInstrumentsDetails,
		queries.SelectAllInstrumentDetails,
	})
	if err != nil {
		return nil, err
	}

	pool := db.Pool()
	return &Database{
		q:  queries.New(pool),
		db: pool,
	}, nil
}

func (db *Database) NewTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := db.db.Begin(ctx)
	return tx, postgres.Wrapf(postgres.BeginTxError(err), "Database.NewTx()")
}

func (db *Database) WithTx(tx pgx.Tx) *Database {
	return &Database{
		q:  db.q.WithTx(tx),
		db: tx,
	}
}
