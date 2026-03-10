package postgres

import (
	"Hades/internal/models"
	"context"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) CreateItem(ctx context.Context, item models.Item) (int, error) {

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), `

	INSERT INTO items (type, amount, date, category, description, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`,

		item.Type, item.Amount, item.Date, item.Category, item.Description, item.CreatedAt)
	if err != nil {
		return 0, err
	}

	var result int
	if err := row.Scan(&result); err != nil {
		return 0, err
	}

	return result, nil

}
