package v1

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
)

const dateLayout = "2006-01-02"

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
		Category:     c.Query("category"),
		Type:         c.Query("type"),
		Sort:         strings.ToUpper(c.Query("sort")),
		SortBy:       strings.ToLower(c.Query("sort_by")),
		From:         from,
		To:           to,
		GroupBy:      strings.ToLower(c.Query("group_by")),
		ExportFormat: strings.ToLower(c.Query("export")),
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

func respondCreated(c *ginext.Context, response any) {
	c.JSON(http.StatusCreated, ginext.H{"result": response})
}

func respondOK(c *ginext.Context, response any) {
	c.JSON(http.StatusOK, ginext.H{"result": response})
}

func fmtRespond(c *ginext.Context, data any, csvFilename string) {
	if c.Query("export") == "csv" {
		writeCSV(c, data, csvFilename)
		return
	}
	respondOK(c, data)
}

func writeCSV(c *ginext.Context, data any, filename string) {

	cd := fmt.Sprintf(`attachment; filename="%s_%s.csv"`, filename, time.Now().Format(dateLayout))

	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", cd)

	writer := csv.NewWriter(c.Writer)
	var writeErr error

	switch values := data.(type) {
	case []models.Item:
		writeErr = writeItems(writer, values)
	case models.Analytics:
		writeErr = writeAnalytics(writer, values)
	default:
		respondError(c, errs.ErrUnsupportedType)
		return
	}

	writer.Flush()
	if writeErr == nil {
		writeErr = writer.Error()
	}

	if writeErr != nil {
		respondError(c, errs.ErrFailedCSV)
	}

}

func writeItems(writer *csv.Writer, items []models.Item) error {

	if err := writer.Write([]string{"ID", "Type", "Amount", "Date", "Category", "Description"}); err != nil {
		return err
	}

	for _, item := range items {
		row := []string{
			fmt.Sprintf("%d", item.ID),
			item.Type,
			item.Amount.String(),
			item.Date.Format(time.RFC3339),
			item.Category,
			item.Description,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil

}

func writeAnalytics(writer *csv.Writer, a models.Analytics) error {

	if err := writer.Write([]string{"Metric", "Value"}); err != nil {
		return err
	}

	rows := [][]string{
		{"Count", fmt.Sprintf("%d", a.Count)},
		{"Total Income", a.TotalIncome.String()},
		{"Total Expense", a.TotalExpense.String()},
		{"Balance", a.Balance.String()},
		{"Average", a.AvgAmount.String()},
		{"Median", a.Median.String()},
		{"90th Percentile", a.Percentile90.String()},
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil

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
		return http.StatusInternalServerError, err.Error()
	}

}
