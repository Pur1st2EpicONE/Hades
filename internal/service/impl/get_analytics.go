package impl

import (
	"Hades/internal/models"
	"context"
)

func (s *Service) GetAnalytics(ctx context.Context, options models.Options) (models.Analytics, error) {

	if err := validateOptions(&options); err != nil {
		return models.Analytics{}, err
	}

	analytics, err := s.storage.GetAnalytics(ctx, options)
	if err != nil {
		s.logger.LogError("service — failed to get analytics from storage", err, "layer", "service.impl")
		return models.Analytics{}, err
	}

	return analytics, nil

}
