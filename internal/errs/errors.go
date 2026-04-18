// Package errs defines reusable error variables for the Hades application.
package errs

import "errors"

var (
	ErrInvalidJSON = errors.New("invalid JSON format") // invalid JSON format
	ErrInvalidID   = errors.New("invalid identifier")  // invalid identifier
	ErrMissingType = errors.New("missing type field")  // missing type field
	ErrInvalidType = errors.New("invalid type value")  // invalid type value

	ErrUnsupportedType = errors.New("unsupported resource type for CSV export") // unsupported resource type for CSV export
	ErrFailedCSV       = errors.New("failed to generate CSV file")              // failed to generate CSV file

	ErrZeroAmount     = errors.New("amount cannot be zero")                // amount cannot be zero
	ErrNegativeAmount = errors.New("amount cannot be negative")            // amount cannot be negative
	ErrAmountTooLarge = errors.New("amount exceeds maximum allowed value") // amount exceeds maximum allowed value

	ErrMissingDate = errors.New("missing date field")                                  // missing date field
	ErrInvalidDate = errors.New("invalid date format, expected RFC3339 or YYYY-MM-DD") // invalid date format, expected RFC3339 or YYYY-MM-DD
	ErrDateTooOld  = errors.New("date is too far in the past")                         // date is too far in the past
	ErrDateTooFar  = errors.New("date is too far in the future")                       // date is too far in the future

	ErrMissingCategory  = errors.New("missing category field")     // missing category field
	ErrCategoryTooShort = errors.New("category name is too short") // category name is too short
	ErrCategoryTooLong  = errors.New("category name is too long")  // category name is too long

	ErrDescriptionTooLong = errors.New("description exceeds maximum length")                   // description exceeds maximum length
	ErrInvalidSortOrder   = errors.New("invalid sort order, use 'asc' or 'desc'")              // invalid sort order, use 'asc' or 'desc'
	ErrInvalidSortBy      = errors.New("invalid sort_by field")                                // invalid sort_by field
	ErrInvalidGroupBy     = errors.New("invalid group_by value, allowed: day, week, category") // invalid group_by value, allowed: day, week, category

	ErrItemNotFound = errors.New("item not found")        // item not found
	ErrInternal     = errors.New("internal server error") // internal server error
)
