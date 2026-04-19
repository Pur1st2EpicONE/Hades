package impl

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const (
	income        = "income"  // income is the allowed value for income type.
	expense       = "expense" // expense is the allowed value for expense type.
	maxAmount     = 1e9       // maxAmount is the maximum allowed transaction amount (1,000,000,000).
	defaultSortBy = "date"    // defaultSortBy is the default field for sorting items.
)

// validateItem performs complete validation of an Item before creation or update.
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

// validateType ensures the type is either "income" or "expense".
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

// validateAmount ensures the amount is positive, non-zero, and below the maximum.
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

// validateDate ensures the date is within the allowed range (±1 year from now).
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

// validateCategory ensures the category is non-empty and between 3 and 100 characters.
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

// validateDescription ensures the description does not exceed 1000 characters.
func validateDescription(desc string) error {
	if len(desc) > 1000 {
		return errs.ErrDescriptionTooLong
	}
	return nil
}

// validateOptions validates and normalizes query options (sorting, grouping, filtering).
// It sets default values for missing optional fields.
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
