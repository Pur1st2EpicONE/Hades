// Package impl provides the concrete implementation of the Service interface.
package impl

import (
	"Hades/internal/logger"
	"Hades/internal/repository"
)

// Service is the concrete implementation of the business logic.
type Service struct {
	logger  logger.Logger      // logger is used for structured logging.
	storage repository.Storage // storage is the data persistence layer.
}

// NewService creates a new Service instance with the given dependencies.
func NewService(logger logger.Logger, storage repository.Storage) *Service {
	return &Service{logger: logger, storage: storage}
}
