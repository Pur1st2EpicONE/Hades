package repository

import (
	"Hades/internal/config"
	"Hades/internal/logger"
	"Hades/internal/models"
	"Hades/internal/repository/postgres"
	"context"
	"fmt"

	"github.com/wb-go/wbf/dbpg"
)

type Storage interface {
	GetItems(ctx context.Context, options models.Options) ([]models.Item, error)
	CreateItem(ctx context.Context, item models.Item) (int, error)
	UpdateItem(ctx context.Context, itemID int, updatedItem models.Item) (models.Item, error)
	DeleteItem(ctx context.Context, itemID int) error
	Close()
}

func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

func ConnectDB(config config.Storage) (*dbpg.DB, error) {

	options := &dbpg.Options{
		MaxOpenConns:    config.MaxOpenConns,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	}

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode), nil, options)
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %w", err)
	}

	if err := db.Master.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil

}
