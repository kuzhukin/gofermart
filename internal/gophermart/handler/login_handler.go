package handler

import (
	"gophermart/internal/gophermart/zlog"
	"net/http"
)

type LoginHandler struct {
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Logger.Info("login handler")
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}
