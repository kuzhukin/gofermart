package handler

import (
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type WithdrawalsHandler struct {
}

func (h *WithdrawalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Info("Balance handler")
}

func NewWithdrawalsHandler() *WithdrawalsHandler {
	return &WithdrawalsHandler{}
}
