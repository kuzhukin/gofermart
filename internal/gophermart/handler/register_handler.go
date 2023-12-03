package handler

import (
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type RegisterHandler struct {
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Info("Register handler")
}

func NewRegisterHandler() *RegisterHandler {
	return &RegisterHandler{}
}
