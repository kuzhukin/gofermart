package handler

import (
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type BalanceHandler struct {
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Info("Balance handler")
}

func NewBalanceHandler() *BalanceHandler {
	return &BalanceHandler{}
}
