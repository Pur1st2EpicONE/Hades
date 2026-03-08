package service

import (
	"Hades/internal/config"
	"Hades/internal/logger"
	"Hades/internal/models"
	"Hades/internal/repository"
	"Hades/internal/service/impl"
	"context"
)

type AuthService interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	CreateToken(userID int64) (string, error)
	GetUserId(ctx context.Context, user models.User) (int64, error)
	ParseToken(tokenString string) (int64, error)
}

type CoreService interface {
	CreateEvent(ctx context.Context, event *models.Event) (string, error)
	CreateBooking(ctx context.Context, userID int64, eventID string) (int64, error)
	ConfirmBooking(ctx context.Context, userID int64, eventID string) error
	CancelBooking(ctx context.Context, bookingID int64) error
	GetInfo(ctx context.Context, eventID string) (*models.Event, error)
	GetAllEvents(ctx context.Context) []models.Event
}

type Service struct {
	AuthService
	CoreService
}

func NewService(logger logger.Logger, config config.Service, storage *repository.Storage) *Service {
	return &Service{
		AuthService: impl.NewAuthService(logger, config, storage.AuthStorage),
		CoreService: impl.NewCoreService(logger, storage.CoreStorage),
	}
}
