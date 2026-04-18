package v1

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

// DeleteItem handles DELETE /api/v1/items/:id.
// It parses the ID from URL, calls service.DeleteItem and returns a status message.
func (h *Handler) DeleteItem(c *ginext.Context) {

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, errs.ErrInvalidID)
		return
	}

	if err := h.service.DeleteItem(c.Request.Context(), itemID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, models.StatusDeleted)

}
