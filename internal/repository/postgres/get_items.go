package postgres

import (
	"Hades/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) GetItems(ctx context.Context, options models.Options) ([]models.Item, error) {

	query, args := buildQuery(options)

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

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

func buildQuery(options models.Options) (string, []any) {

	query := `

	SELECT id, type, amount, date, category, description, created_at
	FROM items
	WHERE TRUE`

	args := []any{}
	argIndex := 1

	if options.Type != "" {
		query += fmt.Sprintf(` AND type = $%d`, argIndex)
		args = append(args, options.Type)
		argIndex++
	}

	if !options.From.IsZero() {
		query += fmt.Sprintf(` AND date >= $%d`, argIndex)
		args = append(args, options.From)
		argIndex++
	}

	if !options.To.IsZero() {
		query += fmt.Sprintf(` AND date <= $%d`, argIndex)
		args = append(args, options.To)
		argIndex++
	}

	if options.Category != "" {
		query += fmt.Sprintf(` AND category = $%d`, argIndex)
		args = append(args, options.Category)
		argIndex++
	}

	if options.Sort != "" {
		query += fmt.Sprintf(` ORDER BY date %s`, options.Sort)
	}

	return query, args

}
