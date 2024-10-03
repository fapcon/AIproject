package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const UniqueConstraintCode = "23505"

var (
	ErrNotFound         = errors.New("row not found")
	ErrUniqueConstraint = errors.New("unique constraint")
)

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == UniqueConstraintCode {
			return fmt.Errorf("%w: %w", ErrUniqueConstraint, err)
		}
	}

	return fmt.Errorf(format+": %w", append(args, err)...)
}
