package impl

import (
	"Hades/internal/models"
	"context"
)

func (s *Service) GetItems(ctx context.Context, options models.Options) ([]models.Item, error) {

	if err := validateOptions(options); err != nil {
		return nil, err
	}

	items, err := s.storage.GetItems(ctx, options)
	if err != nil {
		s.logger.LogError("service — failed to get items from storage", err, "layer", "service.impl")
		return nil, err
	}

	return items, nil

}
