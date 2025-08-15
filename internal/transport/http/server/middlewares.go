package server

import (
	"compress/gzip"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/andrey67895/L0_TEST_TASK/internal/config"
)

func CreateMiddlewares(cfg *config.Config) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     cfg.CORS.AllowOrigins,
			AllowMethods:     cfg.CORS.AllowMethods,
			AllowHeaders:     cfg.CORS.AllowHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           int(cfg.CORS.MaxAge.Seconds()),
		}),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: gzip.BestSpeed,
		}),
	}
}
