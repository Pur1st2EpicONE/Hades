package v1

import (
	"github.com/shopspring/decimal"
)

// CreateItemDTO represents the JSON body for creating a new item.
type CreateItemDTO struct {
	Type        string          `json:"type"`                  // Type of the item (e.g., income, expense)
	Amount      decimal.Decimal `json:"amount"`                // Amount as a decimal value
	Date        string          `json:"date"`                  // Date in RFC3339 or YYYY-MM-DD format
	Category    string          `json:"category"`              // Category name
	Description string          `json:"description,omitempty"` // Optional description
}

// UpdateItemDTO represents the JSON body for updating an existing item.
type UpdateItemDTO struct {
	Type        string          `json:"type"`                  // Type of the item
	Amount      decimal.Decimal `json:"amount"`                // Amount
	Date        string          `json:"date"`                  // Date
	Category    string          `json:"category"`              // Category
	Description string          `json:"description,omitempty"` // Optional description
}
