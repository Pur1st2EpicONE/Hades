package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Item struct {
	ID          int             `db:"id" json:"id"`
	Type        string          `db:"type" json:"type"`
	Amount      decimal.Decimal `db:"amount" json:"amount"`
	Date        time.Time       `db:"date" json:"date"`
	Category    string          `db:"category" json:"category"`
	Description string          `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}
