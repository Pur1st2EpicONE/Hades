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
	From     time.Time
	To       time.Time
	Category string
	Type     string
	Sort     string
}
