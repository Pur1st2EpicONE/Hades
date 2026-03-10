package v1

import (
	"Hades/internal/models"
	"strings"
	"time"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) GetItems(c *ginext.Context) {

	options, err := parseQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}

	items, err := h.service.GetItems(c.Request.Context(), options)
	if err != nil {
		respondError(c, err)
		return
	}

	response := make([]ItemResponseDTO, len(items))

	for i, item := range items {
		response[i] = ItemResponseDTO{
			ID:          item.ID,
			Type:        item.Type,
			Amount:      item.Amount,
			Date:        item.Date.Format(time.RFC3339),
			Category:    item.Category,
			Description: item.Description,
		}
	}

	respondOK(c, response)

}

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
		From:     from,
		To:       to,
	}, nil

}
