package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Item struct {
	ID          int             `json:"id"`
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Date        time.Time       `json:"date"`
	Category    string          `json:"category"`
	Description string          `json:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type Options struct {
	Category     string
	Type         string
	From         time.Time
	To           time.Time
	Sort         string
	SortBy       string
	GroupBy      string
	ExportFormat string
}

type Analytics struct {
	Count        int             `json:"count"`
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	Balance      decimal.Decimal `json:"balance"`
	AvgAmount    decimal.Decimal `json:"avg_amount"`
	Median       decimal.Decimal `json:"median"`
	Percentile90 decimal.Decimal `json:"percentile_90"`
}

type GroupedAnalytics struct {
	GroupKey     string          `json:"group_key"`
	Count        int             `json:"count"`
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	Balance      decimal.Decimal `json:"balance"`
	AvgAmount    decimal.Decimal `json:"avg_amount"`
}
