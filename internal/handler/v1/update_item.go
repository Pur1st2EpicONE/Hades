package v1

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) UpdateItem(c *ginext.Context) {

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, errs.ErrInvalidID)
		return
	}

	var request UpdateItemDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondError(c, errs.ErrInvalidJSON)
		return
	}

	date, err := parseTime(request.Date)
	if err != nil {
		respondError(c, err)
		return
	}

	updatedItem, err := h.service.UpdateItem(c.Request.Context(), itemID, models.Item{
		Type: request.Type, Amount: request.Amount,
		Date: date, Category: request.Category,
		Description: request.Description,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, updatedItem)

}
