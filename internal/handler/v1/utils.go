package v1

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func parseQuery(c *ginext.Context) (models.Options, error) {

	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = parseTime(fromStr)
		if err != nil {
			return models.Options{}, err
		}
	}

	if toStr != "" {
		to, err = parseTime(toStr)
		if err != nil {
			return models.Options{}, err
		}
	}

	return models.Options{
		Category: c.Query("category"),
		Type:     c.Query("type"),
		Sort:     strings.ToUpper(c.Query("sort")),
		SortBy:   strings.ToLower(c.Query("sort_by")),
		From:     from,
		To:       to,
		GroupBy:  strings.ToLower(c.Query("group_by")),
	}, nil

}

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
		errors.Is(err, errs.ErrInvalidSortOrder),
		errors.Is(err, errs.ErrInvalidSortBy),
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
