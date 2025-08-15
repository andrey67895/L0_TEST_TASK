package routes

import (
	"embed"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
)

// Встраиваем статические файлы (HTML, CSS, JS) в приложение
//
//go:embed assets/*
var assets embed.FS

// RegisterRoutes регистрирует все маршруты приложения
func RegisterRoutes(e *echo.Echo) {
	static := e.Group("/assets")
	{
		static.GET("/script.js", serveStaticFiles)
		static.GET("/styles.css", serveStaticFiles)
	}
}

func serveStaticFiles(c echo.Context) error {
	path := c.Request().URL.Path

	if len(path) > 0 {
		_, i := utf8.DecodeRuneInString(path)
		path = path[i:]
	}
	c.Logger().Error(path)
	// Чтение файла из встраиваемых ресурсов
	data, err := assets.ReadFile(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Файл не найден")
	}

	// Устанавливаем тип контента в зависимости от расширения файла
	if strings.HasSuffix(path, ".html") {
		c.Response().Header().Set(echo.HeaderContentType, "text/html")
	} else if strings.HasSuffix(path, ".css") {
		c.Response().Header().Set(echo.HeaderContentType, "text/css")
	} else if strings.HasSuffix(path, ".js") {
		c.Response().Header().Set(echo.HeaderContentType, "application/javascript")
	}

	// Отправляем содержимое файла клиенту
	return c.String(http.StatusOK, string(data))
}
