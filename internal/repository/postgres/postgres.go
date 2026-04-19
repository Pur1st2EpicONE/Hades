// Package postgres implements the Storage interface using PostgreSQL as the backend.
package postgres

import (
	"Hades/internal/config"
	"Hades/internal/logger"

	"github.com/wb-go/wbf/dbpg"
)

// Storage is the PostgreSQL implementation of the repository.Storage interface.
type Storage struct {
	db     *dbpg.DB       // db is the database connection pool.
	logger logger.Logger  // logger is used for structured logging.
	config config.Storage // config holds database and retry settings.
}

// NewStorage creates a new PostgreSQL storage instance.
func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{db: db, logger: logger, config: config}
}

// Close gracefully closes the database master connection.
// It logs success or failure of the close operation.
func (s *Storage) Close() {
	if err := s.db.Master.Close(); err != nil {
		s.logger.LogError("postgres — failed to close properly", err, "layer", "repository.postgres")
	} else {
		s.logger.LogInfo("postgres — database closed", "layer", "repository.postgres")
	}
}
