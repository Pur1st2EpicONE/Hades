// Package app wires application components, provides lifecycle management,
// runs database migrations, and exposes the entry point for booting and running the service.
package app

import (
	"Hades/internal/config"
	"Hades/internal/handler"
	"Hades/internal/logger"
	"Hades/internal/repository"
	"Hades/internal/server"
	"Hades/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
	"github.com/wb-go/wbf/dbpg"
)

// App represents the application's composition root.
// It holds long-lived resources (logger, DB, server) and
// the context/cancel function used for graceful shutdown.
type App struct {
	logger  logger.Logger      // logger is the structured logger used across application layers.
	logFile *os.File           // logFile is the file handle where logs are written.
	server  server.Server      // server is the HTTP server instance.
	ctx     context.Context    // ctx is the root context used to coordinate shutdown across components.
	cancel  context.CancelFunc // cancel cancels the root context when a shutdown signal is received.
	storage repository.Storage // storage is the data storage abstraction backed by the database.
}

// Boot loads configuration, initializes logger, connects to the database,
// runs pending migrations, wires all components and returns a fully constructed *App ready to run.
func Boot() *App {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger, logFile := logger.NewLogger(config.Logger)

	db, err := bootstrapDB(logger, config.Storage)
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err, "layer", "app")
	}

	return wireApp(db, logger, logFile, config)

}

// bootstrapDB establishes a database connection using repository.ConnectDB,
// logs successful connection, sets up goose migration dialect and runs all pending migrations.
func bootstrapDB(logger logger.Logger, config config.Storage) (*dbpg.DB, error) {

	db, err := repository.ConnectDB(config)
	if err != nil {
		return nil, err
	}

	logger.LogInfo("app — connected to database", "layer", "app")

	if err := goose.SetDialect(config.Dialect); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db.Master, config.MigrationsDir); err != nil {
		return nil, fmt.Errorf("failed to apply goose migrations: %w", err)
	}

	logger.Debug("app — migrations applied", "layer", "app")

	return db, nil

}

// wireApp constructs application components (storage, service, handler,
// server), creates a cancellable context and returns the assembled *App.
func wireApp(db *dbpg.DB, logger logger.Logger, logFile *os.File, config config.Config) *App {

	ctx, cancel := newContext(logger)
	storage := repository.NewStorage(logger, config.Storage, db)
	service := service.NewService(logger, storage)
	server := server.NewServer(logger, config.Server, handler.NewHandler(service), cancel)

	return &App{
		logger:  logger,
		logFile: logFile,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
		storage: storage,
	}

}

// newContext creates a context that is cancelled when the process
// receives SIGINT or SIGTERM. It also logs receipt of the signal
// and initiates graceful shutdown by calling the cancel function.
func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		sigString := sig.String()
		if sig == syscall.SIGTERM {
			sigString = "terminate" // sig.String() returns the SIGTERM string in past tense for some reason
		}
		logger.LogInfo("app — received signal "+sigString+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()

	return ctx, cancel

}

// Run starts the server in a background goroutine and blocks
// until the application's context is cancelled. After cancellation it invokes Stop.
func (a *App) Run() {

	go a.server.Run()

	<-a.ctx.Done()

	a.Stop()

}

// Stop performs an orderly shutdown of application components: it shuts down
// the server, waits for background work to finish, closes storage
// and the log file if it is not os.Stdout.
func (a *App) Stop() {

	a.server.Shutdown()
	a.storage.Close()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
