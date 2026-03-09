package postgres

import (
	"Hades/internal/errs"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) DeleteItem(ctx context.Context, itemID int) error {

	row, err := s.db.ExecWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `
	
	DELETE FROM items
	WHERE id = $1`,

		itemID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rows, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of affected rows: %w", err)
	}

	if rows == 0 {
		return errs.ErrItemNotFound
	}

	return nil

}
