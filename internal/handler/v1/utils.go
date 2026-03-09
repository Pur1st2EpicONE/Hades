package v1

import (
	"Hades/internal/errs"
	"errors"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func parseTime(timeStr string) (time.Time, error) {

	if timeStr == "" {
		return time.Time{}, errs.ErrMissingDate
	}

	validTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, errs.ErrInvalidDate
	}

	return validTime.UTC(), nil

}

func respondOK(c *ginext.Context, response any) {
	c.JSON(http.StatusOK, ginext.H{"result": response})
}

func respondCreated(c *ginext.Context, response any) {
	c.JSON(http.StatusCreated, ginext.H{"result": response})
}

func respondError(c *ginext.Context, err error) {
	if err != nil {
		status, msg := mapErrorToStatus(err)
		c.AbortWithStatusJSON(status, ginext.H{"error": msg})
	}
}

func mapErrorToStatus(err error) (int, string) {

	switch {

	case errors.Is(err, errs.ErrInvalidJSON),
		errors.Is(err, errs.ErrInvalidID),
		errors.Is(err, errs.ErrMissingType),
		errors.Is(err, errs.ErrInvalidType),
		errors.Is(err, errs.ErrZeroAmount),
		errors.Is(err, errs.ErrNegativeAmount),
		errors.Is(err, errs.ErrAmountTooLarge),
		errors.Is(err, errs.ErrMissingDate),
		errors.Is(err, errs.ErrInvalidDate),
		errors.Is(err, errs.ErrDateTooOld),
		errors.Is(err, errs.ErrDateTooFar),
		errors.Is(err, errs.ErrMissingCategory),
		errors.Is(err, errs.ErrCategoryTooShort),
		errors.Is(err, errs.ErrCategoryTooLong),
		errors.Is(err, errs.ErrDescriptionTooLong):
		return http.StatusBadRequest, err.Error()

	case errors.Is(err, errs.ErrItemNotFound):
		return http.StatusNotFound, err.Error()

	default:
		return http.StatusInternalServerError, errs.ErrInternal.Error()
	}

}
