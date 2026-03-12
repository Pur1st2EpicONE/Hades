package v1

import (
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

	fmtRespond(c, items, "items")

}
