package server

import (
	"github.com/labstack/echo/v4"

	"github.com/andrey67895/L0_TEST_TASK/internal/config"
	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
	"github.com/andrey67895/L0_TEST_TASK/internal/routes"
	"github.com/andrey67895/L0_TEST_TASK/internal/transport/http/api"
)

type Server struct {
	cfg    *config.Config
	logger *logger.Logger
	echo   *echo.Echo
}

func New(cfg *config.Config, logger *logger.Logger, handlers *APIHandlers, middlewares []echo.MiddlewareFunc) *Server {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		echo:   echo.New(),
	}
	s.echo.Use(middlewares...)
	api.RegisterHandlers(s.echo, handlers)
	routes.RegisterRoutes(s.echo)
	return s
}

// Start запускает HTTP сервер
func (s *Server) Start() error {
	if s.cfg.App.SSLEnable {
		return s.echo.StartTLS(s.cfg.App.AppPort(), s.cfg.App.SSLConfig.CertPath, s.cfg.App.SSLConfig.KeyPath)
	}
	return s.echo.Start(s.cfg.App.AppPort())
}
