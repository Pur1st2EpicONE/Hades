package impl

import (
	"Hades/internal/errs"
	"context"
	"errors"
)

func (s *Service) DeleteItem(ctx context.Context, itemID int) error {

	if err := s.storage.DeleteItem(ctx, itemID); err != nil {
		if !errors.Is(err, errs.ErrItemNotFound) {
			s.logger.LogError("service — failed to delete item from storage", err, "itemID", itemID, "layer", "service.impl")
		}
		return err
	}

	return nil

}
