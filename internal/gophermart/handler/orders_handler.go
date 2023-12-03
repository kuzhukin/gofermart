package handler

import (
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type OrdersHandler struct {
}

func (h *OrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Info("Orders handler")
}

func NewOrdersHandler() *OrdersHandler {
	return &OrdersHandler{}
}
