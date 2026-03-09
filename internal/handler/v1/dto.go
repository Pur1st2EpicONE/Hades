package v1

import (
	"github.com/shopspring/decimal"
)

type CreateItemDTO struct {
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Date        string          `json:"date"`
	Category    string          `json:"category"`
	Description string          `json:"description,omitempty"`
}

type CreateItemResponseDTO struct {
	ID          int
	Type        string
	Amount      decimal.Decimal
	Date        string
	Category    string
	Description string
}
