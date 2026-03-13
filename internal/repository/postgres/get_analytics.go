package postgres

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"context"
	"fmt"

	"github.com/wb-go/wbf/retry"
)

func (s Storage) GetAnalytics(ctx context.Context, options models.Options) (any, error) {

	where, args := buildWhere(options)

	if options.GroupBy != "" {
		return s.grouped(ctx, where, args, options.GroupBy)
	}
	return s.ungrouped(ctx, where, args)

}

func (s Storage) grouped(ctx context.Context, where string, args []any, groupBy string) ([]models.GroupedAnalytics, error) {

	group := ""
	switch groupBy {
	case "day":
		group = "DATE(date)"
	case "week":
		group = "DATE_TRUNC('week', date)"
	case "category":
		group = "category"
	default:
		return nil, errs.ErrInvalidGroupBy
	}

	query := `
		
	SELECT ` + group + `,                
	COUNT(*) as count,
	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
	COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0) as balance,
	COALESCE(AVG(amount), 0) as avg_amount
	FROM items` + where + `
	GROUP BY 1                               
	ORDER BY 1 DESC`

	rows, err := s.db.QueryWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute grouped query: %w", err)
	}
	defer rows.Close()

	var result []models.GroupedAnalytics
	for rows.Next() {
		var g models.GroupedAnalytics
		if err := rows.Scan(
			&g.GroupKey,
			&g.Count,
			&g.TotalIncome,
			&g.TotalExpense,
			&g.Balance,
			&g.AvgAmount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan grouped row: %w", err)
		}
		result = append(result, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil

}

func (s Storage) ungrouped(ctx context.Context, where string, args []any) (models.Analytics, error) {

	query := `
	
	SELECT 
	COUNT(*) as count,
	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
	COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0) as balance,
	COALESCE(AVG(amount), 0) as avg_amount,
	COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) as median,
	COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) as percentile_90
	FROM items` + where

	row, err := s.db.QueryRowWithRetry(ctx, retry.Strategy(s.config.QueryRetryStrategy), query, args...)
	if err != nil {
		return models.Analytics{}, fmt.Errorf("failed to execute analytics query: %w", err)
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
		return models.Analytics{}, fmt.Errorf("failed to scan analytics row: %w", err)
	}

	return analytics, nil

}
