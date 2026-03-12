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

type UpdateItemDTO struct {
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Date        string          `json:"date"`
	Category    string          `json:"category"`
	Description string          `json:"description,omitempty"`
}
