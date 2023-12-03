package handler

import (
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type BalanceWithdrawHandler struct {
}

func (h *BalanceWithdrawHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Info("Balance withdraw handler")
}

func NewBalanceWithdrawHandler() *BalanceWithdrawHandler {
	return &BalanceWithdrawHandler{}
}
