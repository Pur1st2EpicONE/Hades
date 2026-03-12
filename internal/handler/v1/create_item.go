package v1

import (
	"Hades/internal/errs"
	"Hades/internal/models"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) CreateItem(c *ginext.Context) {

	var request CreateItemDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	date, err := parseTime(request.Date)
	if err != nil {
		respondError(c, err)
		return
	}

	createdItem, err := h.service.CreateItem(c.Request.Context(), models.Item{
		Type:        request.Type,
		Amount:      request.Amount,
		Date:        date,
		Category:    request.Category,
		Description: request.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, createdItem)

}
