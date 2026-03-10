package impl

import (
	"Hades/internal/errs"
	"Hades/internal/models"
	"context"
	"database/sql"
	"errors"
)

func (s *Service) UpdateItem(ctx context.Context, itemID int, updatedItem models.Item) (models.Item, error) {

	if err := validateItem(updatedItem); err != nil {
		return models.Item{}, err
	}

	result, err := s.storage.UpdateItem(ctx, itemID, updatedItem)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, errs.ErrItemNotFound
		}
		s.logger.LogError("service — failed to update item in storage", err, "itemID", itemID, "layer", "service.impl")
		return models.Item{}, err
	}

	return result, nil

}
