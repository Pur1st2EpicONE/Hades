package v1

import (
	"Hades/internal/errs"
	"strconv"

	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) DeleteItem(c *ginext.Context) {

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, errs.ErrInvalidID)
	}

	if err := h.service.DeleteItem(c.Request.Context(), itemID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, "deleted")

}
