package v1

import (
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
