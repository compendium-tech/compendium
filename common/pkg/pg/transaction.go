package pg

import (
	"database/sql"
	"errors"

	errorutils "github.com/compendium-tech/compendium/common/pkg/error"
	"github.com/ztrue/tracerr"
)

func DeferRollback(finalErr *error, tx *sql.Tx) {
	errorutils.DeferTry(finalErr, func() error { return rollbackIfNotFinished(tx) })
}

func rollbackIfNotFinished(tx *sql.Tx) error {
	err := tx.Rollback()

	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		return tracerr.Wrap(err)
	}

	return nil
}
