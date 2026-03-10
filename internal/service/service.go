package service

import (
	"Hades/internal/logger"
	"Hades/internal/models"
	"Hades/internal/repository"
	"Hades/internal/service/impl"
	"context"
)

type Service interface {
	GetItems(ctx context.Context, options models.Options) ([]models.Item, error)
	CreateItem(ctx context.Context, item models.Item) (int, error)
	UpdateItem(ctx context.Context, itemID int, updatedItem models.Item) (models.Item, error)
	DeleteItem(ctx context.Context, itemID int) error
}

func NewService(logger logger.Logger, storage repository.Storage) Service {
	return impl.NewService(logger, storage)
}
