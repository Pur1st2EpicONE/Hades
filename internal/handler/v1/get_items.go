package v1

import (
	"github.com/wb-go/wbf/ginext"
)

// GetItems handles GET /api/v1/items.
// It parses query parameters, calls service.GetItems and returns paginated items.
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

	fmtRespond(c, items, "items")

}
