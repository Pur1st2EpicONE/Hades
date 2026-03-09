package postgres

import (
	"Hades/internal/models"
	"context"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) CreateItem(ctx context.Context, item models.Item) (models.Item, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `

	INSERT INTO items (type, amount, date, category, description, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, type, amount, date, category, description, created_at`,

		item.Type, item.Amount, item.Date, item.Category, item.Description, item.CreatedAt)
	if err != nil {
		return models.Item{}, err
	}

	var result models.Item
	if err := row.Scan(&result.ID, &result.Type, &result.Amount, &result.Date, &result.Category, &result.Description, &result.CreatedAt); err != nil {
		return models.Item{}, err
	}

	return result, nil

}
