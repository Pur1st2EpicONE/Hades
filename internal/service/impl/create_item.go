package impl

import (
	"Hades/internal/models"
	"context"
	"time"
)

func (s *Service) CreateItem(ctx context.Context, item models.Item) (models.Item, error) {

	if err := validateItem(item); err != nil {
		return models.Item{}, err
	}

	initialize(&item)

	result, err := s.storage.CreateItem(ctx, item)
	if err != nil {
		s.logger.LogError("service — failed to save item to storage", err, "layer", "service.impl")
		return models.Item{}, err
	}

	return result, nil

}

func initialize(item *models.Item) {
	item.CreatedAt = time.Now().UTC()
}
