package order

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
	"github.com/andrey67895/L0_TEST_TASK/internal/service"
)

type Handler struct {
	log          *logger.Logger
	orderService *service.OrderService
}

func NewHandler(log *logger.Logger, orderService *service.OrderService) *Handler {
	return &Handler{
		log:          log,
		orderService: orderService,
	}
}

// ApiV1GetOrderByOrderUid (GET /order/{order_uid})
func (h *Handler) ApiV1GetOrderByOrderUid(c echo.Context, orderUid string) error {
	ctx := c.Request().Context()
	order, err := h.orderService.GetOrderByUID(ctx, orderUid)
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Заказ с таким ID не найден, обратитесь в поддержку",
		})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Произошла ошибка при получении заказа обратитесь в поддержку",
		})
	}
	return c.JSON(http.StatusOK, order)
}
