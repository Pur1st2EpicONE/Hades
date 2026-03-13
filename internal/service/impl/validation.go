package impl

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	income        = "income"
	expense       = "expense"
	maxAmount     = 1e9
	defaultSortBy = "date"
)

func validateItem(item models.Item) error {

	if err := validateType(item.Type); err != nil {
		return err
	}

	if err := validateAmount(item.Amount); err != nil {
		return err
	}

	if err := validateDate(item.Date); err != nil {
		return err
	}

	if err := validateCategory(item.Category); err != nil {
		return err
	}

	if err := validateDescription(item.Description); err != nil {
		return err
	}

	return nil

}

func validateType(t string) error {

	t = strings.TrimSpace(t)

	if t == "" {
		return errs.ErrMissingType
	}

	if t != income && t != expense {
		return errs.ErrInvalidType
	}

	return nil

}

func validateAmount(amount decimal.Decimal) error {

	if amount.LessThan(decimal.Zero) {
		return errs.ErrNegativeAmount
	}

	if amount.IsZero() {
		return errs.ErrZeroAmount
	}

	if amount.GreaterThan(decimal.NewFromInt(maxAmount)) {
		return errs.ErrAmountTooLarge
	}

	return nil

}

func validateDate(d time.Time) error {

	now := time.Now().UTC()

	if d.Before(now.AddDate(-1, 0, 0)) {
		return errs.ErrDateTooOld
	}

	if d.After(now.AddDate(1, 0, 0)) {
		return errs.ErrDateTooFar
	}

	return nil

}

func validateCategory(category string) error {

	category = strings.TrimSpace(category)

	if category == "" {
		return errs.ErrMissingCategory
	}

	if len(category) < 3 {
		return errs.ErrCategoryTooShort
	}

	if len(category) > 100 {
		return errs.ErrCategoryTooLong
	}

	return nil

}

func validateDescription(desc string) error {
	if len(desc) > 1000 {
		return errs.ErrDescriptionTooLong
	}
	return nil
}

func validateOptions(options *models.Options) error {

	if options.Type != "" && options.Type != income && options.Type != expense {
		return errs.ErrInvalidType
	}

	if options.Sort != "" && options.Sort != "ASC" && options.Sort != "DESC" {
		return errs.ErrInvalidSortOrder
	}

	switch options.SortBy {
	case "", "date", "amount", "type", "category":
	default:
		return errs.ErrInvalidSortBy
	}

	if options.Sort == "" {
		options.Sort = "DESC"
	}

	if options.SortBy == "" {
		options.SortBy = defaultSortBy
	}

	switch options.GroupBy {
	case "", "day", "week", "category":
	default:
		return errs.ErrInvalidGroupBy
	}

	return nil

}
