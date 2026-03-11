package postgres

import (
	"Hades/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) GetAnalytics(ctx context.Context, options models.Options) (models.Analytics, error) {

	condition, args := buildCondition(options)

	query := `

	SELECT 
	COUNT(*) as count,
	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
	COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0) as balance,
	COALESCE(AVG(amount), 0) as avg_amount,
	COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) as median,
	COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) as percentile_90
	FROM items` + condition

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), query, args...)
	if err != nil {
		return models.Analytics{}, fmt.Errorf("failed to execute query: %w", err)
	}

	var analytics models.Analytics
	if err := row.Scan(
		&analytics.Count,
		&analytics.TotalIncome,
		&analytics.TotalExpense,
		&analytics.Balance,
		&analytics.AvgAmount,
		&analytics.Median,
		&analytics.Percentile90,
	); err != nil {
		return models.Analytics{}, fmt.Errorf("failed to scan row: %w", err)
	}

	return analytics, nil

}
