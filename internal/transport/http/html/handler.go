package html

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
)

//go:embed static/*
var static embed.FS

type Handler struct {
	log *logger.Logger
}

func NewHandler(log *logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// GetApiMainIndexHtml получение главной страницы
func (h *Handler) GetApiMainIndexHtml(c echo.Context) error {
	data, err := static.ReadFile("static/index.html")
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Файл не найден")
	}
	return c.HTMLBlob(http.StatusOK, data)
}
