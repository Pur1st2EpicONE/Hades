// Package service defines the business logic interface for the Hades service.
package service

import (
	"Hades/internal/logger"
	"Hades/internal/models"
	"Hades/internal/repository"
	"Hades/internal/service/impl"
	"context"
)

// Service defines the business operations for managing financial items and analytics.
type Service interface {
	GetItems(ctx context.Context, options models.Options) ([]models.Item, error)              // GetItems returns a list of items matching the given filter options.
	CreateItem(ctx context.Context, item models.Item) (models.Item, error)                    // CreateItem creates a new item and returns it with an assigned ID.
	UpdateItem(ctx context.Context, itemID int, updatedItem models.Item) (models.Item, error) // UpdateItem updates an existing item by ID and returns the updated version.
	DeleteItem(ctx context.Context, itemID int) error                                         // DeleteItem removes an item by ID.
	GetAnalytics(ctx context.Context, options models.Options) (any, error)                    // GetAnalytics returns aggregated statistics based on options.
}

// NewService creates a new Service instance using the default implementation.
func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
