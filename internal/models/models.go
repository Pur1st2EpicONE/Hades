// Package models defines the core data structures used throughout the Hades service.
package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Item represents a financial transaction record (income or expense).
type Item struct {
	ID          int             `json:"id"`                    // ID is the unique identifier of the item.
	Type        string          `json:"type"`                  // Type indicates whether it is "income" or "expense".
	Amount      decimal.Decimal `json:"amount"`                // Amount is the monetary value of the transaction.
	Date        time.Time       `json:"date"`                  // Date is the transaction date.
	Category    string          `json:"category"`              // Category groups similar transactions (e.g., "food", "salary").
	Description string          `json:"description,omitempty"` // Description is an optional free‑text field.
	CreatedAt   time.Time       `json:"created_at"`            // CreatedAt is the timestamp when the record was created.
}

// Options holds query parameters for filtering, sorting, grouping, and exporting items or analytics.
type Options struct {
	Category     string    // Category filters by category name.
	Type         string    // Type filters by "income" or "expense".
	From         time.Time // From is the lower bound for the Date field.
	To           time.Time // To is the upper bound for the Date field.
	Sort         string    // Sort defines the order direction: "ASC" or "DESC".
	SortBy       string    // SortBy specifies the field to sort on (e.g., "date", "amount").
	GroupBy      string    // GroupBy aggregates results by "day", "week", or "category".
	ExportFormat string    // ExportFormat can be "csv" to trigger CSV export; otherwise JSON.
}

// Analytics represents aggregated statistics for a set of items.
type Analytics struct {
	Count        int             `json:"count"`         // Count is the total number of items.
	TotalIncome  decimal.Decimal `json:"total_income"`  // TotalIncome is the sum of amounts for income items.
	TotalExpense decimal.Decimal `json:"total_expense"` // TotalExpense is the sum of amounts for expense items.
	Balance      decimal.Decimal `json:"balance"`       // Balance is TotalIncome minus TotalExpense.
	AvgAmount    decimal.Decimal `json:"avg_amount"`    // AvgAmount is the average amount across all items.
	Median       decimal.Decimal `json:"median"`        // Median is the median amount.
	Percentile90 decimal.Decimal `json:"percentile_90"` // Percentile90 is the 90th percentile amount.
}

// GroupedAnalytics holds aggregated statistics per group (e.g., per day, per category).
type GroupedAnalytics struct {
	GroupKey     string          `json:"group_key"`     // GroupKey is the name of the group (e.g., "2025-01-15" or "food").
	Count        int             `json:"count"`         // Count is the number of items in the group.
	TotalIncome  decimal.Decimal `json:"total_income"`  // TotalIncome is the sum of income amounts in the group.
	TotalExpense decimal.Decimal `json:"total_expense"` // TotalExpense is the sum of expense amounts in the group.
	Balance      decimal.Decimal `json:"balance"`       // Balance is TotalIncome minus TotalExpense in the group.
	AvgAmount    decimal.Decimal `json:"avg_amount"`    // AvgAmount is the average amount in the group.
}

const StatusDeleted = "deleted" // StatusDeleted indicates that an item has been deleted.
