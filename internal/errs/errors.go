package errs

import "errors"

var (
	ErrInvalidJSON        = errors.New("invalid JSON format")                         // invalid JSON format
	ErrInvalidID          = errors.New("invalid itemID")                              // invalid itemID
	ErrInternal           = errors.New("internal server error")                       // internal server error
	ErrMissingType        = errors.New("missing type")                                // missing type
	ErrInvalidType        = errors.New("invalid type")                                // invalid type
	ErrInvalidSortOrder   = errors.New("invalid sort order")                          // invalid sort order
	ErrZeroAmount         = errors.New("zero amount")                                 // zero amount
	ErrNegativeAmount     = errors.New("negative amount")                             // negative amount
	ErrAmountTooLarge     = errors.New("amount too large")                            // amount too large
	ErrMissingDate        = errors.New("missing date")                                // missing date
	ErrInvalidDate        = errors.New("invalid event date format, expected RFC3339") // invalid event date format, expected RFC3339
	ErrDateTooOld         = errors.New("date too old")                                // date too old
	ErrDateTooFar         = errors.New("date too far")                                // date too far
	ErrMissingCategory    = errors.New("missing category")                            // missing category
	ErrCategoryTooShort   = errors.New("category too short")                          // category too short
	ErrCategoryTooLong    = errors.New("category too long")                           // category too long
	ErrDescriptionTooLong = errors.New("description too long")                        // description too long
	ErrItemNotFound       = errors.New("item not found")
)
