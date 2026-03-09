package impl

import (
	"Hades/internal/logger"
	"Hades/internal/repository"
)

type Service struct {
	logger  logger.Logger
	storage repository.Storage
}

func NewService(logger logger.Logger, storage repository.Storage) *Service {
	return &Service{logger: logger, storage: storage}
}
