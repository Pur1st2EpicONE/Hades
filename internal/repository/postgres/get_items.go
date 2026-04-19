package postgres

import (
	"Hades/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

// GetItems retrieves a list of items from the database matching the given filter options.
// It builds a query with WHERE conditions and ORDER BY clause based on the options.
// The operation uses a retry strategy defined in the configuration.
func (s Storage) GetItems(ctx context.Context, options models.Options) ([]models.Item, error) {

	query, args := buildQuery(options)

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Amount,
			&item.Date,
			&item.Category,
			&item.Description,
			&item.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return items, nil

}
