// Package repository defines the data persistence layer interface and provides
// database connection utilities for the Hades service.
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

// Storage defines the data persistence operations for items and analytics.
type Storage interface {
	GetItems(ctx context.Context, options models.Options) ([]models.Item, error)              // GetItems returns a list of items matching the given filter options.
	CreateItem(ctx context.Context, item models.Item) (int, error)                            // CreateItem persists a new item and returns its generated ID.
	UpdateItem(ctx context.Context, itemID int, updatedItem models.Item) (models.Item, error) // UpdateItem updates an existing item by ID and returns the updated version.
	DeleteItem(ctx context.Context, itemID int) error                                         // DeleteItem removes an item by ID.
	GetAnalytics(ctx context.Context, options models.Options) (any, error)                    // GetAnalytics returns aggregated statistics based on filter options.
	Close()                                                                                   // Close releases the database connection pool.
}

// NewStorage creates a new Storage instance backed by PostgreSQL.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) Storage {
	return postgres.NewStorage(logger, config, db)
}

// ConnectDB establishes a database connection using the provided configuration.
// It returns a *dbpg.DB instance or an error if connection or ping fails.
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
