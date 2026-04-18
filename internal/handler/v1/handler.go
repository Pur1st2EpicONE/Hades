// Package v1 implements the API version 1 HTTP handlers.
package v1

import (
	"Hades/internal/service"
)

// Handler holds the service dependency for all v1 endpoints.
type Handler struct {
	service service.Service // service is the business logic layer.
}

// NewHandler creates a new v1 Handler with the given service.
func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
