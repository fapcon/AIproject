package queries

import "go.uber.org/multierr"

type batchResults interface {
	Close() error
	Exec(func(int, error))
}

func BatchExec(results batchResults) error {
	defer results.Close()
	var err error
	results.Exec(func(_ int, e error) {
		if e != nil && e.Error() != "batch already closed" {
			err = multierr.Combine(e, results.Close())
		}
	})
	return err
}
