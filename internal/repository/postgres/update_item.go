package postgres

import (
	"Hades/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) UpdateItem(ctx context.Context, itemID int, updatedItem models.Item) (models.Item, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `

	UPDATE items
	SET type = $1, amount = $2, date = $3, category = $4, description = $5
	WHERE id = $6
	RETURNING id, type, amount, date, category, description`,

		updatedItem.Type, updatedItem.Amount, updatedItem.Date, updatedItem.Category, updatedItem.Description, itemID)
	if err != nil {
		return models.Item{}, fmt.Errorf("failed to execute query: %w", err)
	}

	var result models.Item
	if err := row.Scan(
		&result.ID,
		&result.Type,
		&result.Amount,
		&result.Date,
		&result.Category,
		&result.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, err
		}
		return models.Item{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return result, nil

}
