package v1

import (
	"github.com/wb-go/wbf/ginext"
)

// GetAnalytics handles GET /api/v1/analytics.
// It parses query parameters, calls service.GetAnalytics and returns aggregated data.
func (h *Handler) GetAnalytics(c *ginext.Context) {

	options, err := parseQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}

	analytics, err := h.service.GetAnalytics(c.Request.Context(), options)
	if err != nil {
		respondError(c, err)
		return
	}

	fmtRespond(c, analytics, "analytics")

}
