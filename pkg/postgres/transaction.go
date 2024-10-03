package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Tx interface {
	pgx.Tx
	AfterCommit(func())
}

type txWithHooks struct {
	pgx.Tx
	hooks []func()
}

func (t *txWithHooks) AfterCommit(f func()) {
	t.hooks = append(t.hooks, f)
}

type transactioner interface {
	NewTx(context.Context) (pgx.Tx, error)
}

func DoTx(ctx context.Context, db transactioner, f func(tx Tx) error) error {
	tx, err := db.NewTx(ctx)
	if err != nil {
		return err
	}

	withHooks := &txWithHooks{
		Tx:    tx,
		hooks: nil,
	}

	if err = f(withHooks); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			if pgconn.Timeout(err) {
				return err
			}
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	for _, hook := range withHooks.hooks {
		hook()
	}

	return nil
}

func DoTxBatchExec(ctx context.Context, db transactioner, batch *pgx.Batch) error {
	return DoTx(ctx, db, func(tx Tx) error {
		return DoBatchExec(ctx, tx, batch)
	})
}

func DoBatchExec(ctx context.Context, tx pgx.Tx, batch *pgx.Batch) error {
	var err error
	br := tx.SendBatch(ctx, batch)
	defer func() {
		if closeErr := br.Close(); closeErr != nil {
			err = Wrapf(closeErr, "can't close batch (originalError %s)", err)
		}
	}()

	for i := 0; i < batch.Len(); i++ {
		if _, err = br.Exec(); err != nil {
			return Wrapf(err, "can't do tx: can't exec")
		}
	}

	return nil
}

func BeginTxError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("can't begin transaction: %w", err)
}
