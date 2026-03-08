package server

import (
	"Hades/internal/config"
	"Hades/internal/logger"
	"Hades/internal/server/httpserver"
	"context"
	"net/http"
)

type Server interface {
	Run()
	Shutdown()
}

func NewServer(logger logger.Logger, config config.Server, handler http.Handler, cancel context.CancelFunc) Server {
	return httpserver.NewServer(logger, config, handler, cancel)
}
