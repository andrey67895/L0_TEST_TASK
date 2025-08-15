package server

import (
	"github.com/andrey67895/L0_TEST_TASK/internal/transport/http/html"
	"github.com/andrey67895/L0_TEST_TASK/internal/transport/http/order"
)

type OrderHandler = order.Handler
type HtmlHandler = html.Handler

type APIHandlers struct {
	OrderHandler
	HtmlHandler
}

func NewAPIHandlers(orderHandler *order.Handler, htmlHandler *html.Handler) *APIHandlers {
	return &APIHandlers{
		OrderHandler: OrderHandler(*orderHandler),
		HtmlHandler:  HtmlHandler(*htmlHandler),
	}
}
